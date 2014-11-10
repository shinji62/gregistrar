package mbus

import (
	"encoding/json"
	"fmt"
	"github.com/apcera/nats"
	"github.com/cloudfoundry/yagnats"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/shinji62/gregistrar/config"
	"net/url"
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

func NewMessageBusConnection(c *config.Config) (yagnats.NATSConn, error) {

	NatsServer := make([]string, 0)

	for _, info := range c.Nats {
		uri := url.URL{
			Scheme: "nats",
			User:   url.UserPassword(info.User, info.Pass),
			Host:   fmt.Sprintf("%s:%d", info.Host, info.Port),
		}
		NatsServer = append(NatsServer, uri.String())
	}
	return yagnats.Connect(NatsServer)
}

func SendRegisterMessage(natsClient yagnats.NATSConn, msgList []RegisterMessage) error {

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

func RouterStartSubscribe(natsClient yagnats.NATSConn, msgList []RegisterMessage) error {
	//We need to respond to router.start
	_, err := natsClient.Subscribe("router.start", func(m *nats.Msg) {
		SendRegisterMessage(natsClient, msgList)

	})
	if err != nil {
		return err
	}

	return nil

}

func SendGreetingMessage(natsClient yagnats.NATSConn, replyUUID *uuid.UUID) error {
	//send first greeting
	err := natsClient.PublishRequest(greetingTopic, replyUUID.String(), []byte{})
	if err != nil {
		return err
	}
	return nil

}
