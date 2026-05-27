package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nitin-dixit/RepJot/internal/api"
	"github.com/nitin-dixit/RepJot/internal/middleware"
	"github.com/nitin-dixit/RepJot/internal/store"
	"github.com/nitin-dixit/RepJot/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     middleware.UserMiddleware
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	log := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	// our store will go here
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// our handler will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, log)
	userHandler := api.NewUserHanlder(userStore, log)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, log)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:         log,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		Middleware:     middlewareHandler,

		DB: pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
