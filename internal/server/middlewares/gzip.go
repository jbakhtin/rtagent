package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

const GZIPType string = "gzip"

func isSupportsGZIP(encodings []string) bool {
	for _, encode := range encodings {
		if strings.Contains(encode, GZIPType) {
			return true
		}
	}
	return false
}

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (wr gzipWriter) Write(b []byte) (int, error) {
	return wr.Writer.Write(b)
}

func GZIPCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isSupportsGZIP(r.Header.Values("Accept-Encoding")) {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		//w.Header().Set("Content-Encoding", GZIPType)
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func GZIPDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isSupportsGZIP(r.Header.Values("Content-Encoding")) {
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

		//w.Header().Set("Content-Encoding", GZIPType)
		next.ServeHTTP(w, r)
	})
}