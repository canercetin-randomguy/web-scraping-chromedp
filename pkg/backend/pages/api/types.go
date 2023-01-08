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
type CrawlAPIResponse struct {
	CSVDownloadLink  string `json:"csvDownloadLink"`
	JSONDownloadLink string `json:"jsonDownloadLink"`
	Status           string `json:"status"`
}
type DeletePOSTBinding struct {
	Username string `json:"username" binding:"required"`
	AuthKey  string `json:"authkey" binding:"required"`
	FilePath string `json:"filepath" binding:"required"`
}
