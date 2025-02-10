package handler

import (
	"bytes"
	"encoding/json"
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
	log.Println("Received Webhook", c.Request);
	authProvider := auth.NewSimpleKeyProvider(
		os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"),
	)
	// Event is a livekit.WebhookEvent{} object
	event, err := webhook.ReceiveWebhookEvent(c.Request, authProvider)
	if err != nil {
		log.Println("Webhook validation failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid webhook signature"})
		return
	}

	if(event.Event == "participant_left") {
		log.Printf("Participant %s left room %s", event.Participant.Identity, event.Room.Name)

		// convert room id into integer
		id, err := strconv.Atoi(event.Room.Name)
		if err != nil {
			log.Println("Invalid room name", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room name"})
			return
		}

		aroom := room{
			Email:  event.Participant.Identity,
			LiveID: id,
		}


		// prepare Json payload
		bData, err := json.Marshal(aroom)
		if err != nil {
			log.Println("Failed to marshal room data", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal room data"})
			return
		}

		// determine the API endpoint and delete the room if the participant is a host
		var apiUrl string
		if(event.Participant.Permission.CanPublish){
			apiUrl = "https://api.mindlinkstechnology.com/api/AceBeauty/endLive"
			internal.DeleteRoom(event.Room.Name)
		} else {
			apiUrl = "https://api.mindlinkstechnology.com/api/AceBeauty/leftLive"
		}

		// send HttpRequest to the API
		resp , err := http.Post(apiUrl, "application/json", bytes.NewBuffer(bData))
		if(err != nil){
			log.Println("Failed to send request to API", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to API"})
			return
		}
		//  Close the resp body to prevent leaks
		resp.Body.Close()

		// read the response body
		bresp, err := io.ReadAll(resp.Body)
		if(err != nil){
			log.Println("Failed to read response body", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
			return
		}

		endResp := respEndLive{}
        err = json.Unmarshal(bresp, &endResp)
        if err != nil {
            log.Println("Failed to parse API response:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Response decoding failed"})
            return
        }

        if endResp.Valid {
            log.Println("Room update successful:", endResp.Message)
        }

        c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
