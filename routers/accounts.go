package routers

import (
	"net/http"
	"context"
	"fmt"
	"github.com/leighmacdonald/steamid/v3/steamid"
	"github.com/gin-gonic/gin"
	"insider.db/m/models"
	"log"
	"strings"
	"strconv"
	"bytes"
	"encoding/json"
)

func (app *App) SubmitInsider(ctx *gin.Context) {
	type SubmitInsiderData struct {
		PartOfClan string `json:"part_of_clan"`
		Server string `json:"server"`
		// Just nice to know -- Not required because not everyone
		// will have the IP address of a server on the top of their head.
		ServerIP *string `json:"server_ip"`
		InsiderSteamLink string `json:"insider_steam_link"`
		ClanMemberSteamLinks *[]string `json:"other_clan_members"`
		// For non authenticated users. Just another nice piece of information to have when
		// dealing with reports.
		UserOwnerOfClan *bool `json:"clan_owner"`
		// We will get the insider's steam account and view all previous nick names.
		// if the discord happens to be on-record we will use that to hopefully find any other
		// discord names.
		InsiderCurrentName string `json:"insider_current_name"`
	}
	var insider_report SubmitInsiderData
	if err := ctx.BindJSON(&insider_report); err != nil {
		return
	}
	// if app.Config.RequireLoginToSubmit && *app.Config.RequireLoginToSubmit {
	// 	log.Println("login?")
	// }
	report_id := app.SnowflakeGeneratorNode.Generate().Int64()
	split_link := strings.Split(insider_report.InsiderSteamLink, "/")
	last_element := split_link[len(split_link)-1]
	var sid64 steamid.SID64
	if _, err := strconv.Atoi(last_element); err == nil {
		_sid64 := steamid.New(last_element)
		if !_sid64.Valid() {
			ctx.JSON(http.StatusBadRequest, struct{ Message string }{Message: "Provided steam link is invalid."})
			return
		}
		sid64 = _sid64
	} else {
		_sid64, err := steamid.ResolveVanity(context.Background(), insider_report.InsiderSteamLink)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, struct{ Message string }{Message: "Failed to get insider steam profile."})
			return
		}
		sid64 = _sid64
	}
	sid32 := steamid.SID64ToSID32(sid64)
	sid3 := steamid.SID64ToSID3(sid64)
	sid := steamid.SID64ToSID(sid64)

	fmt.Printf("Steam64: %d\n", sid64.Int64())
	fmt.Printf("Steam32: %d\n", steamid.SID64ToSID32(sid64))
	fmt.Printf("Steam3: %s\n", steamid.SID64ToSID3(sid64))
	fmt.Printf("Steam: %s\n", steamid.SID64ToSID(sid64))

	//func (database *Database) CreateReport(id uint64, reporter *uint64, steam_link string, insider_ids SID, server string) *Reports {
	report := app.Database.CreateReport(report_id, nil, insider_report.InsiderSteamLink, models.SID{
		Sid64: sid64.Int64(),
		Sid32: fmt.Sprintf("%d",sid32),
		Sid3: fmt.Sprintf("%s",sid3),
		Sid: fmt.Sprintf("%s",sid),
	},insider_report.Server)
	if report == nil {
		ctx.JSON(http.StatusInternalServerError, struct{ Message string }{Message: "Failed to create report"})
		return
	}
	if app.Config.AlertChannelsWebhook != nil {
		log.Println("Sending webhook?")
		type EmbedField struct {
			Name string `json:"name"`
			Value string `json:"value"`
			Inline bool `json:"inline"`
		}
		type EmbedsFooter struct {
			Text string `json:"text"`
		}
		type Embed struct {
			Title string `json:"title"`
			Description string `json:"description"`
			Color int `json:"color"`
			Fields []EmbedField `json:"fields"`
			Footer EmbedsFooter `json:"footer"`
		}
		type WebhookPayload struct{ 
			Content string `json:"content"` 
			Embeds []Embed `json:"embeds"`
		}
		reqBody, err := json.Marshal(WebhookPayload{
			Content: "A user has been reported. This user has not been confirmed as an insider and will be investigated at the next chance we get. Do not harass anyone posted in this channel, guilty or not.",
			Embeds: []Embed{
			{
				Title: "New Report",
				Description: fmt.Sprintf("User **%s** (current name) has been reported as an insider. Pending investigation, this user may or may not be a verified insider.", insider_report.InsiderCurrentName),
				Color: 10616832,
				Fields: []EmbedField{
					EmbedField{
					Name: "SID64",
					Value: strconv.Itoa(int(sid64.Int64())),
					Inline: true,
				  },
				  EmbedField{
					Name: "SID32",
					Value: fmt.Sprintf("%d",sid32),
					Inline: true,
				  },
				  EmbedField{
					Name: "SID3",
					Value: fmt.Sprintf("%s",sid3),
					Inline: true,
				  },
				  EmbedField{
					Name: "SID",
					Value: fmt.Sprintf("%s",sid),
					Inline: true,
				  },
				  EmbedField{
					Name: "Profile Link",
					Value: insider_report.InsiderSteamLink,
					Inline: true,
				  },
				  EmbedField{
					Name: "Server Reported On",
					Value: insider_report.Server,
					Inline: true,
				  },
				  EmbedField{
					Name: ":no_entry: Harassment of any reported user is not permitted",
					Inline: false,
					Value: "thanks :)",
				  },
				  EmbedField{
					Name: "Whats next?",
					Inline: false,
					Value: "If this user is found to be an insider their account will be watched and any changes will be logged. They will be made pubically available. If you find the user on a different account after they have been confirmed as an insider, please report the new account.",
				  },
				},
				Footer: EmbedsFooter{
				  Text: fmt.Sprintf("Report ID: %d", report.Id),
				},
			  },
		},
		  })
		if err != nil {
			ctx.JSON(http.StatusCreated, struct{ Message string }{Message: "Failed to create webhook object. Report still created."})
			return
		}
		log.Println(*app.Config.AlertChannelsWebhook)
		resp, send_webhook_err := http.Post(
			*app.Config.AlertChannelsWebhook,
			"application/json",
			bytes.NewBuffer(reqBody),
		)
		if send_webhook_err != nil {
			ctx.JSON(http.StatusCreated, struct{ Message string }{Message: "Failed to send webhook to alert channel. Report still created."})
			return
		}
		print(resp)
	}
	ctx.JSON(http.StatusOK, report)
}