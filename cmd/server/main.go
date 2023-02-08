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

	ctx := context.Background()
	repository := inmemory.NewMetricRepository(&ctx)

	metricService := services.NewMetricService(&ctx, repository)

	metricService.Create("test", "test", "test")
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
