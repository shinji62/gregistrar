package mbus

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry/yagnats"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/shinji62/gregistrar/config"
)

const (
	registerTopic    string = "router.register"
	greetingTopic    string = "router.greet"
	routerStartTopic string = "router.start"
)

type RegisterMessage struct {
	Host string   `json:"host"`
	Port uint16   `json:"port"`
	Uris []string `json:"uris"`
}

func NewMessageBusConnection(c *config.Config) (yagnats.NATSClient, error) {
	natsClient := yagnats.NewClient()
	NatsServer := []yagnats.ConnectionProvider{}
	for _, info := range c.Nats {
		NatsServer = append(NatsServer, &yagnats.ConnectionInfo{
			Addr:     fmt.Sprintf("%s:%d", info.Host, info.Port),
			Username: info.User,
			Password: info.Pass,
		})
	}
	err := natsClient.Connect(&yagnats.ConnectionCluster{
		Members: NatsServer,
	})

	return natsClient, err

}

func SendRegisterMessage(natsClient yagnats.NATSClient, msgList []RegisterMessage) error {

	for _, payload := range msgList {
		encodeMessage, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		errNats := natsClient.Publish(registerTopic, encodeMessage)
		fmt.Printf("Error %s", errNats)
		if errNats != nil {
			return errNats
		}

	}

	return nil

}

func RouterStartSubscribe(natsClient yagnats.NATSClient, msgList []RegisterMessage) error {
	//We need to respond to router.start
	_, err := natsClient.Subscribe("router.start", func(mgs *yagnats.Message) {
		SendRegisterMessage(natsClient, msgList)

	})
	return err

}

func SendGreetingMessage(natsClient yagnats.NATSClient, replyUUID *uuid.UUID) error {
	//send first greeting
	err := natsClient.Publish(greetingTopic, []byte{})
	return err

}
