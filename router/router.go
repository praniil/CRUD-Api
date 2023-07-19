package router

import (
	"go-postgres/middleware"

	"github.com/gorilla/mux"
)

//Router is exported and used in main.go

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/newstudent", middleware.CreateStudent).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/getstudent/{id}", middleware.GetStudent).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/getallstudent", middleware.GetAllStudent).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/updatestudent", middleware.UpdateStudent).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deletestudent/{id}", middleware.DeleteStudent).Methods("POST", "OPTIONS")
return router
}
