package models

import (
	"time"
	"log"
	"net"
	"os"
)

type Reports struct {
	Id        	int64		`json:"id"`
	Reporter	*int64 	`json:"reporter"`
	CreatedOn	time.Time `json:"created_on"`
	TargetSteamLink string `json:"target_steam_link"`
	InsiderSid64 int64 `json:"insider_sid64"`
	InsiderSid32 string `json:"insider_sid32"`
	InsiderSid3 string `json:"insider_sid3"`
	InsiderSid string `json:"insider_sid"`
	ServerReportedOn string `json:"server_reported_on`
	ServerReportedOnIp net.IP `json:"server_reported_on_ip"`
	AcceptedBy *int64 `json:"-"`
	AcceptedOn *time.Time `json:"-"`
	Flags uint64 `json:"flags"`
}

func (database *Database) CreateReport(id int64, reporter *uint64, steam_link string, insider_ids SID, server string) *Reports {
	report := &Reports{}
	err := database.Database.
		QueryRow(database.Context, "INSERT INTO reports(id,reporter,target_steam_link,insider_sid64,insider_sid32,insider_sid3,insider_sid,server_reported_on,server_reported_on_ip) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING *", id, reporter, steam_link, insider_ids.Sid64, insider_ids.Sid32, insider_ids.Sid3, insider_ids.Sid, server, nil).
		Scan(&report.Id, &report.Reporter, &report.CreatedOn,&report.TargetSteamLink,&report.InsiderSid64,&report.InsiderSid32,&report.InsiderSid3,&report.InsiderSid,&report.ServerReportedOn,&report.ServerReportedOnIp,&report.AcceptedBy,&report.AcceptedOn,&report.Flags)

		
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	return report
}
