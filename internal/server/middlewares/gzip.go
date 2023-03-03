package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

const gzipType string = "gzip"

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (wr gzipWriter) Write(b []byte) (int, error) {
	return wr.Writer.Write(b)
}

func GZIPCompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k := range r.Header {
			switch k {
			case "Accept-Encoding":
				if !strings.Contains(r.Header.Get(k), gzipType) {
					next.ServeHTTP(w, r)
					return
				}

				gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer gz.Close()

				w.Header().Set("Content-Encoding", gzipType)
				next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
				return
			case "Content-Encoding":
				if !strings.Contains(r.Header.Get(k), gzipType) {
					next.ServeHTTP(w, r)
					return
				}

				gzReader, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer gzReader.Close()

				r.Body = gzReader

				w.Header().Set("Content-Encoding", gzipType)
				next.ServeHTTP(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
