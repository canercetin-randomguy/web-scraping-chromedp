package backend

import (
	"canercetin/pkg/sqlpkg"
	"github.com/gin-gonic/gin"
	"log"
)

func HomeHandler(c *gin.Context) {
	// Search for username in cookies.
	user, _ := c.Cookie("username")
	dbConnection := sqlpkg.SqlConn{}
	err := dbConnection.GetSQLConn("clients")
	if err != nil {
		log.Println(err)
	}
	// then search for auth token in DB.
	auth, err := dbConnection.RetrieveAuthenticationToken(user)
	if err != nil {
		log.Println(err)
	}
	// if auth token is not found, redirect to login page.
	if auth == "" {
		// then delete the cookie.
		c.SetCookie("authtoken", "", -1, "/", "localhost", false, true)
		c.SetCookie("username", "", -1, "/", "localhost", false, true)
		c.Redirect(302, "/signin")
		return
	}
	// if auth token is found, compare it with the cookie.
	cookie, _ := c.Cookie("authtoken")
	if cookie != auth { // if they don't match, redirect to login page.
		c.SetCookie("authtoken", "", -1, "/", "localhost", false, true)
		c.SetCookie("username", "", -1, "/", "localhost", false, true)
		c.Redirect(302, "/signin")
		// then delete the cookie.
		return
	}
	// query to see if user is on free plan.
	planName, err := dbConnection.RetrieveUserPackageDetails(user)
	if err != nil {
		log.Println(err)
		return
	}
	limitAmount, err := dbConnection.RetrieveUserLinkLimit(user)
	if err != nil {
		log.Println(err)
		return
	}
	err = dbConnection.DB.Close()
	if err != nil {
		log.Println(err)
		return
	}
	c.HTML(
		200,
		"home.html",
		gin.H{
			// When users submit the form, it will be sent to /scrape as a post request. Dont forget to hide it.
			"CallbackURL": ScrapingURL,
			"Username":    user,
			"Plan":        planName,
			"Limit":       limitAmount,
		},
	)
}
