// Пакет для работы с обработчиками API.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddComment добавляет новый комментарий.
func AddComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// const operation = "server.api.AddComment"

		var req []Comment
		err := json.NewDecoder(r.Body).Decode(&req)
		fmt.Println(req, err)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req[0].Content == "" {
			http.Error(w, "Empty Comment", http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusCreated)
	}
}
