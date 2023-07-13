package router

import(
	"go-postgres/middleware"
    "github.com/gorilla/mux"
)

//Router is exported and used in main.go

func Router() *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	return router
}