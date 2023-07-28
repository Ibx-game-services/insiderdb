package models

import (
	"time"
)

type Accounts struct {
	Id        	int64		`json:"id"`
	Steam		uint64		`json:"steam"`
	Discord		*int64		`json:"discord"`
	SteamName	string		`json:"steam_name"`
	DiscordName	*string		`json:"discord_name"`
	CreatedOn	time.Time	`json:"created_on"`
	Flags		uint64		`json:"flags"`
	AcceptedBy	*int64		`json:"accepted_by"`
	AcceptedOn	*time.Time	`json:"accepted_on"`
}
