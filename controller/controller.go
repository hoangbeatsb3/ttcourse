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
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/hoangbeatsb3/ttcourse/config"
	"github.com/hoangbeatsb3/ttcourse/service"
	"golang.org/x/net/context"
)

func Serve(cfg *config.Config) {
	r := chi.NewRouter()
	configureRouter(cfg, r)

	addr := cfg.Server.GetFullAddr()
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.RTimeout,
		WriteTimeout: cfg.Server.WTimeout,
	}
	srv.Handler = cors.Default().Handler(r)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt, os.Kill)

	go func() {
		log.Infof("---Main: Serving at http://%s", addr)
		log.Fatal(srv.ListenAndServe())
	}()

	// Graceful shutdown
	<-stop
	log.Info("---Main: Shutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	srv.Shutdown(ctx)
	log.Info("---Main:Goodbye!")
}

func configureRouter(cfg *config.Config, r *chi.Mux) {
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(jsonp.Handler)

	// Routing
	r.Route("/courses", func(r chi.Router) {
		r.Get("/", service.FindAllCourses)
		r.Get("/{name}", service.FindCoursesByName)
		r.Get("/alias/{alias}", service.FindCourseByAlias)
		r.Get("/highest-vote", service.FindHighestVote)
		r.Post("/new", service.CreateCourse)
		r.Post("/vote", service.VoteCourse)
	})
}
