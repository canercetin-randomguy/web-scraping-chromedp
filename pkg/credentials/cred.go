package credentials

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetCredentials() LoginCredentials {
	viper.SetConfigName("credentials")
	viper.SetConfigType("env")
	// credentials.env will be dropped in same folder with main.go, or exe whatever.
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// If login is necessary, get login form and credentials.
	var tempCredentials LoginCredentials
	tempCredentials.Username = viper.Get("USERNAME").(string)
	tempCredentials.Password = viper.Get("PASSWORD").(string)
	tempCredentials.LoginLink = viper.Get("LOGIN_LINK").(string)
	tempCredentials.LoginUsernameField = viper.Get("LOGIN_USERNAME_FIELD").(string)
	tempCredentials.LoginPasswordField = viper.Get("LOGIN_PASSWORD_FIELD").(string)
	tempCredentials.LoginButtonField = viper.Get("LOGIN_BUTTON_FIELD").(string)
	return tempCredentials

}
