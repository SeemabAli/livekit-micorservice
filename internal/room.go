package internal

import (
	"context"

	livekit "github.com/livekit/protocol/livekit"
)

func ListParticipants(roomName string) ([]*livekit.ParticipantInfo, error) {
	resp, err := roomClient.ListParticipants(context.Background(), &livekit.ListParticipantsRequest{
		Room: roomName,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetParticipants(), nil
}
