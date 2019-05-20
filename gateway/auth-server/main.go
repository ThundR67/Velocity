package main

import (
	"net/http"

	"github.com/SonicRoshan/Velocity/gateway/auth-server/handlers"
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
)

//Loding Logger
var log = logger.GetLogger("logs/auth_server.log")

func main() {
	urlHandler := handlers.Handlers{}
	urlHandler.Init()
	http.HandleFunc("/signin", urlHandler.SignInHandler)
	http.HandleFunc("/signin-fresh", urlHandler.SignInFreshHandler)
	log.Critical(http.ListenAndServe(config.AuthServerConfigIPAddress, nil))
}
