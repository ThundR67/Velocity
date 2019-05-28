package router

import (
	"github.com/SonicRoshan/Velocity/gateway/auth-server/handler"
	"github.com/gorilla/mux"
)

//GetRouter is used to create all routes
func GetRouter() *mux.Router {
	handler := handler.Handler{}
	handler.Init()
	router := mux.NewRouter()

	//Routes
	router.HandleFunc("/sign-in", handler.SignInHandler).Methods("GET")
	router.HandleFunc("/sign-in-fresh", handler.SignInFreshHandler).Methods("GET")
	router.HandleFunc("/refresh", handler.RefreshHandler).Methods("GET")
	router.HandleFunc("/sign-up", handler.SignUpHandler).Methods("POST")

	return router
}
