package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nitin-dixit/goProject/internal/api"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	log := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	// our store will go here

	// our handler will go here
	workoutHandler := api.NewWorkoutHandler()
	app := &Application{
		Logger:         log,
		WorkoutHandler: workoutHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Status is available")
}
