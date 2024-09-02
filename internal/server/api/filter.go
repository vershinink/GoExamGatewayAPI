// Пакет для работы с обработчиками API.
package api

import (
	"encoding/json"
	"net/http"
)

// Filter возвращает отфильтрованный список постов.
func Filter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// const operation = "server.api.Filter"

		var resp []NewsShortDetailed
		w.Header().Set("Access-Control-Allow-Origin", "*")
		resp = HardCodeNews
		enc := json.NewEncoder(w)
		err := enc.Encode(resp)
		if err != nil {
			http.Error(w, "failed to encode news", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}
}
