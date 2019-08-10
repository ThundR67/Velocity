package handler

import (
	"net/http"

	"github.com/SonicRoshan/Velocity/global/config"
)

//EmailVerificationHandler is used to handle email verification route
func (handler Handler) EmailVerificationHandler(w http.ResponseWriter, r *http.Request) {
	verificationCode := r.URL.Query().Get(config.AuthServerConfigVerificationCodeField)
	if verificationCode == "" {
		handler.respond(w, nil, "Verification Code Not Provided", nil)
		return
	}

	email, err := handler.emailVerification.Verify(verificationCode)
	if err != nil {
		handler.respond(w, nil, "", err)
		return
	}

	msg, err := handler.users.Activate(email)
	if msg != "" || err != nil {
		handler.respond(w, nil, msg, err)
		return
	}

	handler.respond(w, map[string]string{
		"Status": "Your Account Has Been Verified",
	}, "", nil)

}
