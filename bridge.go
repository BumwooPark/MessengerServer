package main

import (
	"github.com/nsqio/go-nsq"
	"fmt"
)

func bridgeSocketQueue(NSQ NsqClientAction, wsClient *Client){
	defer NSQ.Close()

	go NSQ.Consumer(wsClient.nsqAddr, func(message *nsq.Message) error {
		wsClient.send <- message.Body
		return nil
	})

	for{
		select {
		case message := <- wsClient.recv:
			NSQ.Producer(wsClient.topic,message)
		case message := <- wsClient.send:
			fmt.Println("ws.send")
			fmt.Println(message)
		}
	}
}

