package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Latest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// const operation = "server.api.Latest"

		param := r.URL.Query().Get("page")
		if param == "" {
			param = "1"
		}
		page, err := strconv.Atoi(param)
		if err != nil {
			http.Error(w, "incorrect page number", http.StatusBadRequest)
			return
		}
		_ = page

		var resp []NewsShortDetailed
		w.Header().Set("Access-Control-Allow-Origin", "*")
		resp = HardCodeNews
		enc := json.NewEncoder(w)
		err = enc.Encode(resp)
		if err != nil {
			http.Error(w, "failed to encode news", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}
}
