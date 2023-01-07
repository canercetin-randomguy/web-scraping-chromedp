package backend

import "fmt"

type SignUpFormBinding struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type SignInFormBinding struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type ScrapingFormBinding struct {
	Username  string `json:"username"`
	LinkLimit string `json:"linkLimit"`
	MainLink  string `json:"mainLink"`
	MaxDepth  string `json:"maxDepth"`
}

// !!!! Only change this if you want to change the Port. !!!!!
var Port = 7171

// These will be called when user clicks on the submit button on forms.
//
// Change these to your willings.
var SignupCallbackURL = fmt.Sprintf("http://localhost:%d/v1/signup/callback", Port)
var SigninCallbackURL = fmt.Sprintf("http://localhost:%d/v1/signin/callback", Port)

// These will be used pretty frequently.
var SignupURL = fmt.Sprintf("http://localhost:%d/v1/signup", Port)
var SigninURL = fmt.Sprintf("http://localhost:%d/v1/signin", Port)
var HomeURL = fmt.Sprintf("http://localhost:%d/v1/home", Port)
var ScrapingURL = fmt.Sprintf("http://localhost:%d/v1/home/scraping/callback", Port)

var SignupPath = "/v1/signup"
var SigninPath = "/v1/signin"
var HomePath = "/v1/home"
var ScrapingPath = "/v1/home/scraping/callback"
var FileStoragePath = "/v1/storage"
