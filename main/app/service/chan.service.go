package service

import (
	"time"
)

var chans map[string]chan string = make(map[string]chan string)

var limitForTimer time.Duration = 10 * time.Second

func getChan(room string) chan string {
	messageChan := chans[room]

	if messageChan != nil {
		return messageChan
	}

	newChan := make(chan string, 1)
	timer := time.NewTimer(limitForTimer)
	chans[room] = newChan

	go func() {
		open := true
		for open {
			select {
			case message := <-newChan:
				{
					saveMessageRedis(room, message)
				}
			case <-timer.C:
				{
					chans[room] = nil
					open = false
					break
				}
			}
		}
	}()

	return newChan
}
