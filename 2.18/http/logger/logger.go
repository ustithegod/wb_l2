package logger

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)

		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Msg("")
	})
}
