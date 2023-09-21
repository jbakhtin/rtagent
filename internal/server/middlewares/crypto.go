package middlewares

import (
	"bytes"
	"crypto/rsa"
	"github.com/jbakhtin/rtagent/pkg/crypto"
	"io"
	"net/http"
)

func Decrypt(privateKey *rsa.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, request *http.Request) {
			if privateKey != nil {
				chyper, err := io.ReadAll(request.Body)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}
				request.Body.Close()
				data, err := crypto.GetDecryptedMessage(privateKey, chyper, request.Header.Get("Encrypted-Key"))
				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}

				request.Body = io.NopCloser(bytes.NewReader(data))
				request.ContentLength = int64(len(data))
			}
			next.ServeHTTP(writer, request)
		}
		return http.HandlerFunc(fn)
	}
}
