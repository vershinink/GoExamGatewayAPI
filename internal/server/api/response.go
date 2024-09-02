// Пакет для работы с обработчиками API.
package api

import "time"

// NewsFullDetailed - структура подробного представления новости.
type NewsFullDetailed struct {
	News     NewsShortDetailed
	Comments []Comment
}

// NewsShortDetailed - структура короткого представления новости.
type NewsShortDetailed struct {
	ID      string `json:"ID"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
	PubTime int64  `json:"PubTime"`
	Link    string `json:"Link"`
}

// Comment - структура комментария.
type Comment struct {
	ID       string `json:"CommentID"`
	ParentID string `json:"ParentID"`
	NewsID   string `json:"NewsID"`
	Content  string `json:"CommentContent"`
}

// Тестовые значения.
var HardCodeNews = []NewsShortDetailed{
	{ID: "news01", Title: "Title 1", Content: "Content 1", PubTime: time.Now().Unix(), Link: "https://google.com"},
	{ID: "news02", Title: "Title 2", Content: "Content 2", PubTime: time.Now().Unix(), Link: "https://yandex.ru"},
	{ID: "news03", Title: "Title 3", Content: "Content 3", PubTime: time.Now().Unix(), Link: "https://bing.com"},
}

// Тестовые значения.
var CommentsNews1 = []Comment{
	{ID: "n01c01", ParentID: "", NewsID: "news01", Content: "Comment 1"},
}

// Тестовые значения.
var CommentsNews2 = []Comment{
	{ID: "n02c01", ParentID: "", NewsID: "news02", Content: "Comment 1"},
	{ID: "n02c02", ParentID: "", NewsID: "news02", Content: "Comment 2"},
	{ID: "n02c03", ParentID: "n02c02", NewsID: "news02", Content: "Comment 3"},
	{ID: "n02c04", ParentID: "n02c02", NewsID: "news02", Content: "Comment 4"},
}

// Тестовые значения.
var CommentsNews3 = []Comment{
	{ID: "n03c01", ParentID: "", NewsID: "news03", Content: "Comment 1"},
	{ID: "n03c02", ParentID: "n03c01", NewsID: "news03", Content: "Comment 2"},
}
