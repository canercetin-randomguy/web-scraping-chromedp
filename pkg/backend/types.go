package backend

type SignUpFormBinding struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type SignInFormBinding struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var SignupCallbackURL = "http://localhost:6969/signup/callback"
var SigninCallbackURL = "http://localhost:6969/signin/callback"
