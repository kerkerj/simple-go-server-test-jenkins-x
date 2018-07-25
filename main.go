package gosrv

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"tmps/simple-go-server-test/models"
)

type Route interface {
}

type App struct {
	Server *http.Server
	Db     sql.DB
	Router *http.ServeMux
}

var app App
var logger = log.New(os.Stdout, "http: ", log.LstdFlags)

func main() {
	NewApp().Run()
}

func NewApp() *App {
	app := &App{
		Router: http.NewServeMux(),
		Server: &http.Server{
			Addr:         ":9876",
			Handler:      nil,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			ErrorLog:     logger,
		},
	}
	app.setRoutes()

	return app
}

func (a *App) setRoutes() {
	a.Router.HandleFunc("/", LogMiddleware(IndexRouter()))

	a.Server.Handler = a.Router
}

func (a *App) Run() {
	logger.Printf("Listening on 9876\n")
	if err := a.Server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Printf("Cannot start server on 9876\nError: %s", err.Error())
	}
}

func IndexRouter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.NewUser(1, "test", "bio")
		fmt.Fprintf(w, "Index - "+user.Name)
	}
}

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Printf("[%s] %q [error: %s]\n", r.Method, r.URL.String(), err)
			}
		}()

		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		logger.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
}
