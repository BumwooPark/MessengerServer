package main

import (
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
)

type (
	ClientBuilder interface {
		Producer(ipAddr string) ClientBuilder
		Consumer(ipAddr string, Topic string, Channel string) ClientBuilder
		Build() NsqClientAction
	}

	NsqClientBuilder struct {
		producer *nsq.Producer
		consumer *nsq.Consumer
	}

	NsqClientAction interface {
		Producer(topic string, message []byte)
		Consumer(ipAddr string,handler func(message *nsq.Message) error)
		Close()
	}

	NsqClient struct {
		producer *nsq.Producer
		consumer *nsq.Consumer
	}
)

func NewBuilder() ClientBuilder {
	return &NsqClientBuilder{nil, nil}
}

func (client *NsqClientBuilder) Producer(ipAddr string) ClientBuilder {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(ipAddr, config)
	if err != nil {
		log.Panic(err)
	}
	client.producer = p
	return client
}

func (client *NsqClientBuilder) Consumer(ipAddr string, Topic string, Channel string) ClientBuilder {
	decodeConfig := nsq.NewConfig()
	c, err := nsq.NewConsumer(Topic,Channel, decodeConfig)
	if err != nil {
		log.Panic(err)
	}
	client.consumer = c
	return client
}

func (client *NsqClientBuilder) Build() NsqClientAction {
	return &NsqClient{client.producer,client.consumer}
}

func (client *NsqClient) Producer(topic string, message []byte) {
	err := client.producer.Publish(topic,message)
	if err != nil {
		panic(err)
	}
}

func (client *NsqClient) Consumer(ipAddr string,handler func(message *nsq.Message) error) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	client.consumer.AddHandler(nsq.HandlerFunc(handler))

	err := client.consumer.ConnectToNSQD(ipAddr)
	if err != nil {
		log.Panic("Could not connect")
	}
	wg.Wait()
}


func (client *NsqClient) Close(){
	client.consumer.Stop()
	client.producer.Stop()
}