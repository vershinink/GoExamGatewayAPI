package server

import (
	"GoExamGatewayAPI/internal/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var comment string = `{"parentId": "","postId": "66e3d4e492358bac466b0861","content": "Test comment!"}`

func TestNews(t *testing.T) {
	t.Parallel()
	logger.Discard()

	client := &http.Client{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", News("google.com", client))
	srv := httptest.NewServer(mux)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("News() error = unexpected status code, got = %v, want = %v", rr.Code, http.StatusOK)
	}
}

func TestNewsById(t *testing.T) {
	t.Parallel()
	logger.Discard()

	tests := []struct {
		name     string
		newsId   string
		wantCode int
	}{
		{
			name:     "Status Code 200",
			newsId:   "66e3d4e492358bac466b0861",
			wantCode: http.StatusOK,
		},
		{
			name:     "Status Code 404",
			newsId:   "66e3d4e492358bac466b0999",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Status Code 200",
			newsId:   "asdfg",
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &http.Client{}
			mux := http.NewServeMux()
			mux.HandleFunc("GET /news/id/{id}", NewsById("192.168.0.150:10501", "192.168.0.150:10502", client))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			uri := "/news/id/" + tt.newsId
			req := httptest.NewRequest(http.MethodGet, uri, nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("NewsById() error = unexpected status code, got = %v, want = %v", rr.Code, tt.wantCode)
			}
		})
	}
}

func TestAddComment(t *testing.T) {
	t.Parallel()
	logger.Discard()

	client := &http.Client{}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /comments/new", AddComment("192.168.0.150:10502", "192.168.0.150:10503", client))
	srv := httptest.NewServer(mux)
	defer srv.Close()

	body := strings.NewReader(comment)
	req := httptest.NewRequest(http.MethodPost, "/comments/new", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("News() error = unexpected status code, got = %v, want = %v", rr.Code, http.StatusCreated)
	}
}
