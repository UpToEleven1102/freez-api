package controllers

import (
	c "git.nextgencode.io/huyen.vu/freez-app-rest/config"
	auth "git.nextgencode.io/huyen.vu/freez-app-rest/identity"
	"net"
	"net/http"
)

/*AuthHandler ... - HandleFunc for auth route*/
func AuthHandler(w http.ResponseWriter, req *http.Request, route string, userType string) {
	if req.Method != "POST" {
		http.NotFound(w, req)
		return
	}
	switch route {
	case c.Email:
		if userType == c.Verify {
			auth.VerifyEmailPin(w, req)
		} else if len(userType) == 0 {
			auth.SendRandomPinEmail(w, req)
		} else {
			http.NotFound(w, req)
		}
	case c.PhoneNumber:
		if userType == c.Verify {
			auth.VerifySMSPin(w, req)
		} else if len(userType) == 0 {
			auth.SendRandomPinSMS(w, req)
		} else {
			http.NotFound(w, req)
		}
	case c.SignUp:
		if userType == c.Merchant {
			ip , _, _ := net.SplitHostPort(req.RemoteAddr)

			auth.SignUpMerchant(w, req, ip)
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
	case c.Exist:
		if userType == c.Email {
			auth.EmailExists(w, req)
		} else if userType == c.PhoneNumber {
			auth.PhoneNumberExists(w, req)
		} else {
			http.NotFound(w, req)
		}
	case c.Facebook:
		auth.AuthenticateFacebook(w, req, userType)
	default:
		http.NotFound(w, req)
	}
}
