package main

import (
 "fmt"
 "time"
 "net/http"
 lgr "github.com/go-pkgz/lgr"
 "github.com/pkg/errors"

 //"github.com/didip/tollbooth/v6"
 //"github.com/didip/tollbooth_chi"
 "github.com/go-chi/chi/v5"
 "github.com/go-chi/chi/v5/middleware"
)

type Server struct {
    Host           string
    Port           string
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
}

func (s Server) Run() error {
	if err := http.ListenAndServe(s.Host+":"+s.Port, s.routes()); err != http.ErrServerClosed {
		return errors.Wrap(err, "server failed")
	}
	return nil
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()

    router.Use(middleware.Logger)
	router.Use(middleware.Throttle(1000), middleware.Timeout(60*time.Second))
	//router.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(10, nil)))


	router.Route("/", func(r chi.Router) {
	    r.Get("/", s.homePage)
	    r.Get("/ws", s.serveWS)
	})

    lgr.Printf("[INFO] Activate rest server")
    lgr.Printf("[INFO] Host: %s", s.Host)
    lgr.Printf("[INFO] Port: %s", s.Port)

	return router
}

func (s Server) serveWS(w http.ResponseWriter, r *http.Request) {
    manager := NewManager();
    manager.serveWS(w, r)
}

func (s Server) homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home Page")
}
