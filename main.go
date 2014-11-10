package main

import (
	"flag"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/shinji62/gregistrar/config"
	"github.com/shinji62/gregistrar/helpers"
	"github.com/shinji62/gregistrar/mbus"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "", "Configuration file ")
	flag.Parse()
}

func main() {
	c := config.DefaultConfig()
	replyUUID, err := uuid.NewV4()
	if err != nil {
		fmt.Print("Could not get an uuid")
	}

	if configFile != "" {
		c = config.InitConfigFromFile(configFile)
	}

	// setup number of procs
	if c.GoMaxProcs != 0 {
		runtime.GOMAXPROCS(c.GoMaxProcs)
	}

	natsClient, err := mbus.NewMessageBusConnection(c)
	if err != nil {
		fmt.Print("Error connecting to Nats")
		os.Exit(1)
	}
	defer natsClient.Close()

	//Get localIp
	ipList, err := helpers.LocalsIP()

	if err != nil {
		fmt.Print("Unable to get any local IP")
		os.Exit(1)
	}

	uris := make([]string, 0)
	for _, info := range c.Uris {
		uris = append(uris, info.Uri)
	}

	payloadList := make([]mbus.RegisterMessage, 0)

	for _, ip := range ipList {

		message := mbus.RegisterMessage{
			Host: ip,
			Port: c.Port,
			Uris: uris,
		}
		payloadList = append(payloadList, message)
	}

	signals := make(chan os.Signal, 1)
	errChan := make(chan error)

	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	err = mbus.SendGreetingMessage(natsClient, replyUUID)
	if err != nil {
		fmt.Printf("Could not send Greetings: %s", err)
	}

	err = mbus.SendRegisterMessage(natsClient, payloadList)
	if err != nil {
		fmt.Printf("Could not send First Register: %s", err)
	}

	err = mbus.RouterStartSubscribe(natsClient, payloadList)
	if err != nil {
		fmt.Printf("Could not subscribe: %s", err)
	}

	//TODO use router interval message
	//Sending registring message every Interval time
	ticker := time.NewTicker(time.Second * 2)

	go func() {
		for t := range ticker.C {
			fmt.Println("Ticker %s", t)
			errChan <- mbus.SendRegisterMessage(natsClient, payloadList)
		}
	}()
	for {
		select {
		case err := <-errChan:
			if err != nil {
				fmt.Printf("Exited by error %s", err)
				os.Exit(1)
			}
		case sig := <-signals:
			if sig == syscall.SIGUSR1 {
				//Dump config
			}
			fmt.Printf("Exited by signal %s", sig)
			os.Exit(0)
		}
	}
	defer close(signals)
	os.Exit(0)

}
