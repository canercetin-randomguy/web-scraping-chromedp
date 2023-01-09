package links

type LinkStruct struct {
	Link          string `json:"Link"`
	IsBroken      bool   `json:"IsBroken"`
	OutsideDomain bool   `json:"OutsideDomain"`
	StatusCode    int    `json:"StatusCode"`
}
