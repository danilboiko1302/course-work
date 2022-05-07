package controller

import (
	"course-work/app/nats"
	"course-work/app/service"
	"course-work/app/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetRoute() *gin.Engine {
	route := gin.Default()

	route.LoadHTMLGlob("app/pages/html/*.html")

	route.Static("/js", "app/pages/js")
	route.Static("/css", "app/pages/css")

	route.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	route.GET("/:roomName", func(c *gin.Context) {
		c.HTML(http.StatusOK, "room.html", nil)
	})

	room := route.Group("/room")

	{
		room.GET("/:roomName", func(c *gin.Context) {
			roomName := c.Param("roomName")
			fmt.Println("New user in room", roomName)
			ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				log.Println("error get connection")
				log.Fatal(err)
			}
			defer ws.Close()

			user := &types.User{}
			messagesChan := make(chan string, 1)
			unsub, err := nats.Connection.Sub(roomName, func(data []byte) {
				messagesChan <- string(data)
			})

			go func() {
				for {
					message := <-messagesChan
					if user.LoggedIn {
						ws.WriteMessage(websocket.TextMessage, []byte(message))
					}
				}
			}()

			if err != nil {
				log.Println("error when sub to nats")
				return
			}

			go func() {
				<-c.Request.Context().Done()
				err := unsub()

				if user.LoggedIn {
					service.RemoveFromRoom(roomName, user)

					err = service.Pub(roomName, &types.MessageFront{
						Action: types.UserLeft,
						Data:   user.Name,
					})

					if err != nil {
						log.Println("error pub msg nats: " + err.Error())
					}

				}

				if err != nil {
					log.Println("error when unsub from nats")
				}
			}()

			for {
				var msg types.Message
				err = ws.ReadJSON(&msg)

				if err != nil {
					log.Println("error read json")
					return
				}

				pubMsg, myMsg, err := service.HandleMessage(roomName, msg, user)

				if err != nil {
					log.Println(err.Error())

					bytes, err := json.Marshal(gin.H{"error": err.Error()})
					if err != nil {
						log.Println("error json.Marshal: " + err.Error())
						continue
					}
					messagesChan <- string(bytes)
					continue
				}

				if pubMsg != nil {
					err = service.Pub(roomName, pubMsg)
					if err != nil {
						log.Println("error pub msg nats: " + err.Error())
					}
				}

				if myMsg != nil {
					bytes, err := json.Marshal(myMsg)
					if err != nil {
						log.Println("error json.Marshal: " + err.Error())
						continue
					}
					messagesChan <- string(bytes)
				}

			}

		})
	}

	return route
}
