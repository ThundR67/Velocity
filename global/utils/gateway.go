package utils

import (
	"encoding/json"
	"net/http"

	"github.com/SonicRoshan/Velocity/global/config"
	"go.uber.org/zap"
)

//GatewayRespond is used to respond to client with json data
func GatewayRespond(
	w http.ResponseWriter, data map[string]string, msg string, err error, log *zap.Logger) {

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Error("Error While Reponding", zap.Error(err))

		output := config.InternalServerError
		if config.DebugMode {
			output = err.Error()
		}

		data = map[string]string{
			config.AuthServerConfigErrField: output,
		}
	} else if msg != "" {
		data = map[string]string{
			config.AuthServerConfigErrField: msg,
		}
	}

	json.NewEncoder(w).Encode(data)
}
