package controller

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jsonp"
	"github.com/hoangbeatsb3/ttcourse/config"
	"github.com/hoangbeatsb3/ttcourse/service"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func Serve(cfg *config.Config) {
	r := chi.NewRouter()
	configureRouter(cfg, r)

	addr := cfg.Server.GetFullAddr()
	srv := &http.Server{
		Addr:         addr,
		Handler:      cors.Default().Handler(r),
		ReadTimeout:  cfg.Server.RTimeout,
		WriteTimeout: cfg.Server.WTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill)

	go func() {
		logrus.Info("---Main: Serving at http://%s", addr)
	}()

	// Graceful shutdown
	<-stop
	logrus.Info("---Main: Shutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	srv.Shutdown(ctx)
	logrus.Info("---Main:Goodbye!")
}

func configureRouter(cfg *config.Config, r *chi.Mux) {
	r.Use(middleware.DefaultLogger, middleware.DefaultCompress, middleware.Recoverer, middleware.Logger, jsonp.Handler)

	// Routing
	r.Route("/courses", func(r chi.Router) {
		r.Get("/", service.FindAllCourses)
		r.Get("/{name}", service.FindCoursesByName)
		r.Get("/alias/{alias}", service.FindCourseByAlias)
		r.Get("/highest-vote", service.FindHighestVote)
		r.Post("/", service.CreateCourse)
		r.Post("/vote", service.VoteCourse)
	})
}
