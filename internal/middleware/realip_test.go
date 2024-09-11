package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// serve - заглушка для тестов.
func serve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := []byte("Hello world!")
		w.Write(b)
	}
}

func Test_headerIP(t *testing.T) {
	tests := []struct {
		name   string
		header string
		value  string
		want   string
	}{
		{
			name:   "X-Real-IP",
			header: "X-Real-IP",
			value:  "88.88.88.88:12345",
			want:   "88.88.88.88:12345",
		},
		{
			name:   "X-Forwarded-For",
			header: "X-Forwarded-For",
			value:  "88.88.88.88:12345,99.99.99.99:12345,192.168.0.1:12345",
			want:   "88.88.88.88:12345",
		},
		{
			name:   "Empty",
			header: "True-Client-IP",
			value:  "88.88.88.88:12345",
			want:   "",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(tt.header, tt.value)

			if got := headerIP(req); got != tt.want {
				t.Errorf("headerIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRealIP(t *testing.T) {
	tests := []struct {
		name   string
		header string
		value  string
		want   string
	}{
		{
			name:   "X-Real-IP",
			header: "X-Real-IP",
			value:  "192.168.0.1:1234",
			want:   "192.168.0.1:1234",
		},
		{
			name:   "X-Forwarded-For",
			header: "X-Forwarded-For",
			value:  "88.88.88.88:12345,99.99.99.99:12345,192.168.0.1:12345",
			want:   "88.88.88.88:12345",
		},
		{
			name:   "Empty",
			header: "True-Client-IP",
			value:  "88.88.88.88:12345",
			want:   "192.0.2.1:1234",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			mux.Handle("GET /", RealIP(serve()))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(tt.header, tt.value)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			got := req.Header.Get(xRealIP)

			if got != tt.want {
				t.Errorf("RealIP() error = %s, want %s", got, tt.want)
			}
		})
	}

	// mux := http.NewServeMux()
	// mux.Handle("GET /", RealIP(serve()))
	// srv := httptest.NewServer(mux)
	// defer srv.Close()

	// req := httptest.NewRequest(http.MethodGet, "/", nil)
	// rr := httptest.NewRecorder()

	// mux.ServeHTTP(rr, req)

	// if req.Header.Get(xRealIP) == "" {
	// 	t.Errorf("RealIP() error = X-Real-IP is not set")
	// }
}
