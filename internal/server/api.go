package server

import (
	"GoExamGatewayAPI/internal/logger"
	"GoExamGatewayAPI/internal/middleware"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Post - структура новостного поста.
type Post struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	PubTime time.Time `json:"pubTime"`
	Link    string    `json:"link"`
}

// Comment - структура комментария к посту.
type Comment struct {
	ID       string    `json:"id"`
	ParentID string    `json:"parentId"`
	PostID   string    `json:"postId"`
	PubTime  time.Time `json:"pubTime"`
	Content  string    `json:"content"`
}

type FullComment struct {
	Comment
	Childs []FullComment `json:"childs"`
}

// FullPost - структура ноавостного поста с деревом комментариев к нему.
type FullPost struct {
	Post
	Comments []FullComment `json:"comments"`
}

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequset = errors.New("bad request")
)

// News проксирует запрос на получение списка новостей с пагинацией
// в сервис новостей по переданному адресу.
func News(host string, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.News"

		log := slog.Default().With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("request to receive posts")

		resp, err := request(host, r, client)
		if err != nil {
			log.Error("failed to receive posts", logger.Err(err))
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		log.Debug("posts received successfully")

		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Error("failed to copy response body", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		log.Info("request served successfully")
	}
}

// NewsById проксирует запрос на получение одной новости по ее ID
// вместе с комментариями в сервис новостей и с сервис комментариев
// по переданным адресам.
func NewsById(hostNews, hostComments string, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.NewsById"

		log := slog.Default().With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("request to receive post by ID with comments")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		// Объявим функцию, которая вызывает request и декодирует
		// ответ в структуру.
		fn := func(host string, req *http.Request, data any) error {
			resp, err := request(host, req, client)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound {
				return ErrNotFound
			}
			if resp.StatusCode == http.StatusBadRequest {
				return ErrBadRequset
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("request status = %d", resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("readAll err = %w", err)
			}

			err = json.Unmarshal(body, data)
			if err != nil {
				return fmt.Errorf("unmarshal err = %w", err)
			}
			return nil
		}

		var post Post
		var comments []FullComment
		var errProxy error
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		rNews := r.Clone(r.Context())
		rComm := r.Clone(ctx)

		// Сделаем запросы в сервис новостей и сервис комментариев
		// в отдельных горутинах.
		wg.Add(2)
		go func() {
			defer wg.Done()
			err := fn(hostNews, rNews, &post)
			// Если при получении новости возникла ошибка, то сохраняем ее
			// для дальнейшей обработки, затем прерываем запрос комментариев.
			if err != nil {
				log.Error("failed to receive post", logger.Err(err))
				errProxy = err
				cancel()
				return
			}
			log.Debug("post received successfuly")
		}()

		go func() {
			defer wg.Done()
			uri := rComm.URL.Path
			uri = strings.ReplaceAll(uri, "news/id", "comments")
			rComm.URL.Path = uri
			err := fn(hostComments, rComm, &comments)
			// Если при получении комментариев возникла ошибка, то не обрабатываем
			// ее как в горутине получения новости, так как ошибка получения новости
			// критична, а получения комментариев - нет.
			if err != nil {
				log.Error("failed to receive comments", logger.Err(err))
				return
			}
			log.Debug("comments received successfuly")
		}()

		wg.Wait()

		if errProxy != nil {
			log.Error("failed to find post by ID", logger.Err(errProxy))
			if errors.Is(errProxy, ErrNotFound) {
				http.Error(w, "post not found", http.StatusNotFound)
				return
			}
			if errors.Is(errProxy, ErrBadRequset) {
				http.Error(w, "incorrect post id", http.StatusBadRequest)
				return
			}
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		fullPost := FullPost{
			Post:     post,
			Comments: comments,
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "\t")
		err := enc.Encode(fullPost)
		if err != nil {
			log.Error("failed to encode post and comments", logger.Err(err))
			http.Error(w, "failed to encode post and comments", http.StatusInternalServerError)
			return
		}
		log.Info("request served successfuly")
	}
}

// AddComment проверяет запрос в сервисе цензурирования комментариев.
// В случае успеха проксирует запрос на добавление нового комментария
// в сервис комментариев.
func AddComment(hostComments, hostCensor string, client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.AddComment"

		log := slog.Default().With(
			slog.String("op", operation),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("request to add new comment")

		// Создаем копии тела запроса.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request body", logger.Err(err))
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		rc1 := io.NopCloser(bytes.NewBuffer(body))
		rc2 := io.NopCloser(bytes.NewBuffer(body))

		// Клонируем запрос с немодифицированнымм телом и меняем
		// путь запроса.
		rCens := r.Clone(r.Context())
		rCens.Body = rc1
		rCens.URL.Path = ""

		log.Debug("checking new comment")

		respCensor, err := request(hostCensor, rCens, client)
		if err != nil {
			log.Error("failed to check comment", logger.Err(err))
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer respCensor.Body.Close()
		io.Copy(io.Discard, respCensor.Body)

		if respCensor.StatusCode != http.StatusOK {
			log.Error("comment contains inappropriate words", slog.Int("code", respCensor.StatusCode))
			http.Error(w, "comment contains inappropriate words", http.StatusBadRequest)
			return
		}
		log.Debug("comment checked successfully")

		// Клонируем запрос с немодифицированнымм телом.
		rComm := r.Clone(r.Context())
		rComm.Body = rc2

		respComm, err := request(hostComments, rComm, client)
		if err != nil {
			log.Error("failed to add new comment", logger.Err(err))
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		defer respComm.Body.Close()

		// Полностью копируем ответ сервиса комментариев в ResponseWriter.
		copyHeader(w.Header(), respComm.Header)
		w.WriteHeader(respComm.StatusCode)
		_, err = io.Copy(w, respComm.Body)
		if err != nil {
			log.Error("failed to copy response body", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		log.Info("request served successfully")
	}
}
