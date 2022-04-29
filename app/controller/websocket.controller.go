package controller

import (
	"bytes"
	"course-work/app/nats"
	"course-work/app/types"
	"encoding/gob"
	"fmt"
	"io"
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

	room := route.Group("/room")

	{
		room.GET("/:name", func(c *gin.Context) {
			name := c.Param("name")
			fmt.Println("New user in room", name)
			ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				log.Println("error get connection")
				log.Fatal(err)
			}
			defer ws.Close()

			// dataChan := make(chan types.Message)

			// go func() {
			// 	for {
			// 		// select {}
			// 		message := <-dataChan
			// 		fmt.Println(message)
			// 	}

			// }()

			unsub, err := nats.Connection.Sub(name, func(data []byte) {
				go func(r io.Reader) {
					var msg types.Message

					err := gob.NewDecoder(r).Decode(&msg)

					if err != nil {
						log.Println("error decoding message")
						return
					}

					err = ws.WriteJSON(msg)
					if err != nil {
						log.Println("error write json: " + err.Error())
					}

					// dataChan <- msg
				}(bytes.NewReader(data))
			})

			if err != nil {
				log.Println("error when sub to nats")
				return
			}

			go func() {
				<-c.Request.Context().Done()
				err := unsub()
				if err != nil {
					log.Println("error when unsub from nats")
				}
				// close(dataChan)
			}()

			for {
				var msg types.Message
				err = ws.ReadJSON(&msg)
				if err != nil {
					log.Println("error read json")
					// log.Fatal(err)
					return
				}

				var b bytes.Buffer

				err := gob.NewEncoder(&b).Encode(msg)

				if err != nil {
					log.Println("error convert msg to bytes")
					// log.Fatal(err)
					return
				}
				err = nats.Connection.Pub(name, b.Bytes())
				if err != nil {
					log.Println("error convert msg to bytes")
					return
				}
			}

			// var data struct {
			// 	A string `json:"a"`
			// 	B int    `json:"b"`
			// }
			// //Read data in ws
			// err = ws.ReadJSON(&data)
			// data.A = name
			// if err != nil {
			// 	log.Println("error read json")
			// 	log.Fatal(err)
			// }

			// //Write ws data, pong 10 times
			// var count = 0
			// for {
			// 	count++
			// 	if count > 10 {
			// 		break
			// 	}

			// 	err = ws.WriteJSON(struct {
			// 		A string `json:"a"`
			// 		B int    `json:"b"`
			// 		C int    `json:"c"`
			// 	}{
			// 		A: data.A,
			// 		B: data.B,
			// 		C: count,
			// 	})
			// 	if err != nil {
			// 		log.Println("error write json: " + err.Error())
			// 	}
			// 	time.Sleep(1 * time.Second)
			// }
		})
	}

	return route
}
