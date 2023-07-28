package models

import (
	"time"
)

type ClanMembers struct {
	Id        	int64		`json:"id"`
	Clan		int64 		`json:"clan"`
	Member		int64 		`json:"member"`
	AddedOn	time.Time 	`json:"added_on"`
	Owner	*bool `json:"owner"`
}
