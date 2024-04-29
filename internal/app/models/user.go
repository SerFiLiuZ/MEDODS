package models

type User struct {
	GUID string   `json:"guid"`
	Rts  []string `json:"rts"`
}
