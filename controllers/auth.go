package controllers

import (
	c "git.nextgencode.io/huyen.vu/freeze-app-rest/config"
	auth "git.nextgencode.io/huyen.vu/freeze-app-rest/identity"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, req *http.Request, route string, userType string) {
	switch route {
	case c.SignUp:
		if userType == c.Merchant {
			auth.SignUpMerchant(w, req)
		} else if userType == c.User {
			auth.SignUpUser(w, req)
		} else {
			http.NotFound(w, req)
		}
	case c.SignIn:
		if userType == c.Merchant {
			auth.SignInMerchant(w, req)
		} else if userType == c.User {
			auth.SignInUser(w, req)
		} else {
			http.NotFound(w, req)
		}
	case c.UserInfo:
		if len(userType) > 0 {
			http.NotFound(w, req)
		}
		auth.GetUserInfo(w, req)
	case c.Verify:
		if len(userType) == 0 {
			auth.AccountExists(w, req)
		} else {
			http.NotFound(w, req)
		}
	default:
		http.NotFound(w, req)
	}
}
