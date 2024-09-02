package api

import (
	"encoding/json"
	"net/http"
)

func Detailed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// const operation = "server.api.Detailed"

		id := r.PathValue("id")
		_ = id

		var resp = NewsFullDetailed{
			News:     HardCodeNews[1],
			Comments: CommentsNews2,
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		enc := json.NewEncoder(w)
		err := enc.Encode(resp)
		if err != nil {
			http.Error(w, "failed to encode detailed news", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}
}
