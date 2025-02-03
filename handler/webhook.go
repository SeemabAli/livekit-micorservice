package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/LinuxSploit/Livekit-Mircroservice/internal"
	"github.com/gin-gonic/gin"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/webhook"
)

type room struct {
	Email  string `json:"email"`
	LiveID int    `json:"LiveID"`
}

type respEndLive struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

func WebhookHandler(c *gin.Context) {
	authProvider := auth.NewSimpleKeyProvider(
		os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"),
	)
	// Event is a livekit.WebhookEvent{} object
	event, err := webhook.ReceiveWebhookEvent(c.Request, authProvider)
	if err != nil {
		// Could not validate, handle error
		return
	}
	// Consume WebhookEvent

	switch event.Event {
	case "participant_left":

		if event.Participant.Permission.CanPublish {
			internal.DeleteRoom(event.Room.Name)
			id, err := strconv.Atoi(event.Room.Name)
			if err != nil {
				log.Println(err)
				return
			}

			aroom := room{
				Email:  event.Participant.Identity,
				LiveID: id,
			}
			bData, err := json.Marshal(aroom)
			if err != nil {
				log.Println(err)
				return
			}

			resp, err := http.Post("https://api.acebeauty.club/api/AceBeauty/endLive", "application/json", bytes.NewBuffer(bData))
			if err != nil {
				log.Println(err)
				return
			}

			bresp, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}

			endResp := respEndLive{}
			err = json.Unmarshal(bresp, &endResp)
			if err != nil {
				log.Println(err)
				return
			}

			if endResp.Valid {
				fmt.Println("room successfully deleted!")
			}
		} else {
			id, err := strconv.Atoi(event.Room.Name)
			if err != nil {
				log.Println(err)
				return
			}

			aroom := room{
				Email:  event.Participant.Identity,
				LiveID: id,
			}
			bData, err := json.Marshal(aroom)
			if err != nil {
				log.Println(err)
				return
			}

			resp, err := http.Post("https://api.acebeauty.club/api/AceBeauty/leftLive", "application/json", bytes.NewBuffer(bData))
			if err != nil {
				log.Println(err)
				return
			}

			bresp, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}

			endResp := respEndLive{}
			err = json.Unmarshal(bresp, &endResp)
			if err != nil {
				log.Println(err)
				return
			}

			if endResp.Valid {
				fmt.Println(aroom.Email+" has left Live Stram: ", aroom.LiveID)
			}
		}

		// http.Post("https://api.acebeauty.club/api/AceBeauty/leftLive", "application/json")
	}

}
