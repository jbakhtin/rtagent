package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

const gzipType string = "gzip"

var gzipWriterPool sync.Pool

func acquireGzipWriter(w io.Writer, lvl int) (zw *gzip.Writer, err error) {
	v := gzipWriterPool.Get()
	if v == nil {
		return gzip.NewWriterLevel(w, lvl)
	}

	zw = v.(*gzip.Writer)
	zw.Reset(w)
	return
}

func releaseGzipWriter(zw *gzip.Writer) {
	zw.Close()
	gzipWriterPool.Put(zw)
}

type gzipWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (wr gzipWriter) Write(b []byte) (int, error) {
	return wr.Writer.Write(b)
}

func GZIPCompressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), gzipType) {
			gz, err := acquireGzipWriter(w, gzip.BestSpeed)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", gzipType)
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

			releaseGzipWriter(gz)
			return
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), gzipType) {
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

		next.ServeHTTP(w, r)
	})
}
