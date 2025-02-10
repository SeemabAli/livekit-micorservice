package internal

import (
	"context"
	"os"
	"time"
	"log"
	"github.com/livekit/protocol/auth"
	livekit "github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go"
)

var host = "http  s://livekit.tssclinicallabs.com"
var roomClient *lksdk.RoomServiceClient

func Init() {
	roomClient = lksdk.NewRoomServiceClient(host, os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"))
}

func GetHostToken(room, host, identity string) string {
	at := auth.NewAccessToken(os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"))

	canPublish := true
	canSubscripbe := true

	grant := &auth.VideoGrant{
		RoomAdmin:    true,
		RoomJoin:     true,
		Room:         room,
		CanPublish:   &canPublish,
		CanSubscribe: &canSubscripbe,
	}

	at.AddGrant(grant).SetName(host).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	token, _ := at.ToJWT()

	return token
}

func GetJoinToken(room, host, identity string) (string, error) {
	at := auth.NewAccessToken(os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"))

	canPublish := false
	canSubscripbe := true
	canPublishData := true

	grant := &auth.VideoGrant{
		RoomJoin:       true,
		Room:           room,
		CanPublish:     &canPublish,
		CanSubscribe:   &canSubscripbe,
		CanPublishData: &canPublishData,
	}

	at.AddGrant(grant).SetName(host).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		return token, err
	}

	return token, nil
}

func CreateRoom(RoomName string) (*livekit.Room, error) {
	var roomClient = lksdk.NewRoomServiceClient(host, os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"))
	room, err := roomClient.CreateRoom(context.Background(), &livekit.CreateRoomRequest{
		Name:            RoomName,
		EmptyTimeout:    30, // 30 sec
		MaxParticipants: 100,
	})
	if err != nil {
		return room, err
	}
	log.Println("Room Created Successfully" , room.Name)
	return room, err
}

func ListAllRooms() ([]*livekit.Room, error) {
	rooms, err := roomClient.ListRooms(context.Background(), &livekit.ListRoomsRequest{})
	if err != nil {
		return nil, err
	}
	log.Println("Rooms:" , rooms.GetRooms())
	return rooms.GetRooms(), nil
}

func DeleteRoom(roomName string) error {
	_, err := roomClient.DeleteRoom(context.Background(), &livekit.DeleteRoomRequest{
		Room: roomName,
	})
	if err != nil {
		return err
	}
	log.Println("Room Deleted Successfully", roomName)
	return nil
}

func KickOutParticipant(RoomName string, identity string) (*livekit.RemoveParticipantResponse, error) {
	var roomClient = lksdk.NewRoomServiceClient(host, os.Getenv("LIVEKIT_API_KEY"), os.Getenv("LIVEKIT_API_SECRET"))
	res, err := roomClient.RemoveParticipant(context.Background(), &livekit.RoomParticipantIdentity{
		Room:     RoomName,
		Identity: identity,
	})
	if err != nil {
		return nil, err
	}

	return res, err
}
