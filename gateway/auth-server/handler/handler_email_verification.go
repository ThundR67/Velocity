package handler

import (
	"net/http"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/utils"
)

//EmailVerificationHandler is used to handle email verification route
func (handler Handler) EmailVerificationHandler(w http.ResponseWriter, r *http.Request) {
	verificationCode := r.URL.Query().Get(config.AuthServerConfigVerificationCodeField)
	if verificationCode == "" {
		utils.GatewayRespond(w, nil, "Verification Code Not Provided", nil, log)

		return
	}

	email, err := handler.emailVerification.Verify(verificationCode)
	if err != nil {
		utils.GatewayRespond(w, nil, "", err, log)
		return
	}

	msg := handler.users.Activate(email)
	if msg != "" {
		utils.GatewayRespond(w, nil, msg, nil, log)
		return
	}

	utils.GatewayRespond(w, map[string]string{
		"Status": "Your Account Has Been Verified",
	}, "", nil, log)

}
