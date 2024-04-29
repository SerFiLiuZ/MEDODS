package models

import "time"

type AuthUser struct {
	AtSigned     string    `json:"ats"`
	AtExpiration time.Time `json:"ate"`
	RtSigned     string    `json:"rts"`
	RtExpiration time.Time `json:"rte"`
}
