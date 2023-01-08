package api

type AuthPOSTBinding struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email"`
	SecretKey string `json:"secretkey"`
}
type CrawlPOSTBinding struct {
	Username string `json:"username" binding:"required"`
	AuthKey  string `json:"authkey" binding:"required"`
	MaxDepth string `json:"maxdepth" binding:"required"`
	MainLink string `json:"mainlink" binding:"required"`
}
