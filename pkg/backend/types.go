package backend

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	// This can be either:
	//
	// Sign-up page
	//
	// Sign-in page
	//
	// Both are sending callbacks to the same endpoint, so we can use the same struct. And FormType is necessary.
	FormType string `json:"formType"`
}

var callbackURL = "http://localhost:6969/callback"
