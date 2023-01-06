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
	LinkLimit string `json:"linklimit"`
	MainLink  string `json:"mainlink"`
}

// !!!! Only change this if you want to change the Port. !!!!!
var Port = 7171

// These will be called when user clicks on the submit button on forms.
//
// Change these to your willings.
var SignupCallbackURL = fmt.Sprintf("http://localhost:%d/signup/callback", Port)
var SigninCallbackURL = fmt.Sprintf("http://localhost:%d/signin/callback", Port)

// These will be used pretty frequently.
var SignupURL = fmt.Sprintf("http://localhost:%d/signup", Port)
var SigninURL = fmt.Sprintf("http://localhost:%d/signin", Port)
var HomeURL = fmt.Sprintf("http://localhost:%d/home", Port)
var ScrapingURL = fmt.Sprintf("http://localhost:%d/home/scraping/callback", Port)
