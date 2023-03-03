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
	gh := func(writer http.ResponseWriter, request *http.Request) {
		if !isSupportsGZIP(request.Header.Values("Accept-Encoding")) {
			next.ServeHTTP(writer, request)
			return
		}

		gzWriter, err := gzip.NewWriterLevel(writer, gzip.BestSpeed)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gzWriter.Close()

		writer.Header().Set("Content-Encoding", GZIPType)
		next.ServeHTTP(gzipWriter{ResponseWriter: writer, Writer: gzWriter}, request)

	}
	return http.HandlerFunc(gh)
}

func GZIPDecompress(next http.Handler) http.Handler {
	gh := func(writer http.ResponseWriter, request *http.Request) {
		if !isSupportsGZIP(request.Header.Values("Content-Encoding")) {
			next.ServeHTTP(writer, request)
			return
		}

		gzReader, err := gzip.NewReader(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gzReader.Close()

		request.Body = gzReader //ReadCloser = Reader

		writer.Header().Set("Content-Encoding", GZIPType)
		next.ServeHTTP(writer, request)
	}
	return http.HandlerFunc(gh)
}