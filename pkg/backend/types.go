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

// These will be called when user clicks on the submit button on forms.
//
// Change these to your willings.
var SignupCallbackURL = "http://localhost:6969/signup/callback"
var SigninCallbackURL = "http://localhost:6969/signin/callback"
