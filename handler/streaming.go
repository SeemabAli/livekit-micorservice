package handler

import (
	"fmt"
	"net/http"

	"github.com/LinuxSploit/Livekit-Mircroservice/internal"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type RoomData struct {
	RoomName string `binding:"required"`
	Name     string `binding:"required"`
	Identity string `binding:"required"`
}

func HostLiveHandler(c *gin.Context) {

	var data RoomData
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	room, err := internal.CreateRoam(data.RoomName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	hostToken := internal.GetHostToken(room.Name, data.Name, data.Identity)

	c.JSON(http.StatusOK, gin.H{"token": hostToken})

}

func JoinLiveHandler(c *gin.Context) {

	var data RoomData
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	//debug line
	fmt.Println("Join data: ", data)

	rooms, err := internal.ListAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	for _, r := range rooms {
		if r.Name == data.RoomName {
			joinToken, err := internal.GetJoinToken(data.RoomName, data.Name, data.Identity)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": joinToken})

			return
		}
	}

	c.JSON(http.StatusBadRequest, err)

}

type Room struct {
	Name string
}

func ListLiveHandler(c *gin.Context) {
	var LiveRoomsData []Room

	rooms, err := internal.ListAllRooms()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	for _, v := range rooms {

		LiveRoomsData = append(LiveRoomsData, Room{
			Name: v.GetName(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"LiveRooms": LiveRoomsData})

}
