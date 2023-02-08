package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jbakhtin/rtagent/internal/repositories/storages/inmemory"
	"github.com/jbakhtin/rtagent/internal/services"
	"net/http"
)

// TODO: создать репозиторий с подменой

var form = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/login" method="post">
            <label>Логин</label><input type="text" name="login">
            <label>Пароль<input type="password" name="password">
            <input type="submit" value="Login">
        </form>
    </body>
</html>`

func getMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, form)
}

func updateMetric(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query().Get("type")
	if q == "" {
		http.Error(w, "The type parameter is missing", http.StatusBadRequest)
		return
	}

	q = r.URL.Query().Get("key")
	if q == "" {
		http.Error(w, "The key parameter is missing", http.StatusBadRequest)
		return
	}

	q = r.URL.Query().Get("value")
	if q == "" {
		http.Error(w, "The value parameter is missing", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	repository := inmemory.NewMetricRepository(&ctx)

	metricService := services.NewMetricService(&ctx, repository)

	metricService.Create("test", "test", "test")

	w.Header().Set("content-type", "text/plain")
	// устанавливаем статус-код 200
	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println("Server started")

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Get("/", getMetrics)

		r.Post("/update/{type}/{key}/{value}", updateMetric)
	})

	http.ListenAndServe(":8080", r)
}
