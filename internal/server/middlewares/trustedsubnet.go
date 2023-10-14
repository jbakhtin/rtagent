package middlewares

import (
	"github.com/go-faster/errors"
	"net/http"
	"strings"
)

type Configer interface {
	GetTrustedSubnet() string
}

func TrustedSubnet(cfg Configer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, request *http.Request) {
			if cfg.GetTrustedSubnet() == "" {
				next.ServeHTTP(writer, request)
				return
			}

			if !strings.Contains(request.Header.Get("X-Real-IP"), cfg.GetTrustedSubnet()) {
				http.Error(writer, errors.New("Request rejected").Error(), http.StatusForbidden)
				return
			}

			next.ServeHTTP(writer, request)
		}

		return http.HandlerFunc(fn)
	}
}
