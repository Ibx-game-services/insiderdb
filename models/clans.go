package models

import (
	"time"
)

type Clans struct {
	Id        	int64		`json:"id"`
	Name		string		`json:"steam"`
	Marker 		*string		`json:"marker"`
	Owner		*int64		`json:"owner"`
	AddedBy		int64		`json:"added_by`
	CreatedOn	time.Time	`json:"created_on"`
}
