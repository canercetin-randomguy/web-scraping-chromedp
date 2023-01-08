package pages

import "fmt"

type SignUpFormBinding struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
type SignInFormBinding struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type ScrapingFormBinding struct {
	Username  string `json:"username" binding:"required"`
	LinkLimit string `json:"linkLimit" binding:"required"`
	MainLink  string `json:"mainLink" binding:"required"`
	MaxDepth  string `json:"maxDepth" binding:"required"`
	AuthKey   string `json:"authkey"`
}
type SecretKeyFormBinding struct {
	Username string `json:"username" binding:"required"`
	AuthKey  string `json:"authkey" binding:"required"`
}

// !!!! Only change this if you want to change the Port. !!!!!
var Port = 7171
var LoopbackPort = 7172

// These will be called when user clicks on the submit button on forms.
//
// Change these to your willings.
var SignupCallbackURL = fmt.Sprintf("http://localhost:%d/private/signup/callback", LoopbackPort)
var SigninCallbackURL = fmt.Sprintf("http://localhost:%d/private/signin/callback", LoopbackPort)
var ScrapingURL = fmt.Sprintf("http://localhost:%d/private/home/scraping/callback", LoopbackPort)
var SecretKeyCallbackURL = fmt.Sprintf("http://localhost:%d/private/secretkey/callback", LoopbackPort)

// These will be used pretty frequently.
var SignupURL = fmt.Sprintf("http://localhost:%d/v1/signup", Port)
var SigninURL = fmt.Sprintf("http://localhost:%d/v1/signin", Port)
var HomeURL = fmt.Sprintf("http://localhost:%d/v1/home", Port)

var SignupPath = "/v1/signup"
var SigninPath = "/v1/signin"
var HomePath = "/v1/home"
var ScrapingPath = "/v1/home/scraping/callback"
var FileStoragePath = "/v1/storage"
var SecretKeyPath = "/v1/secretkey"
