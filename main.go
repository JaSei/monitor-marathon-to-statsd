package main

import (
	"fmt"
	marathon "github.com/gambol99/go-marathon"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/alexcesaro/statsd.v2"
	"log"
	"strings"
)

const VERSION = "0.1.0"

var (
	marathonUrl   = kingpin.Flag("marathon", "Marathon URL").Required().Envar("MARATHON").String()
	statsdAddress = kingpin.Flag("statsd", "statsd/statsite address").Required().String()
)

func main() {
	kingpin.Version(VERSION)
	kingpin.Parse()

	stat, err := statsd.New(statsd.Address(*statsdAddress), statsd.Prefix("mmcloud.marathon.status_update"))

	if err != nil {
		log.Fatal(err)
	}
	// Configure client
	config := marathon.NewDefaultConfig()
	config.URL = *marathonUrl
	config.EventsTransport = marathon.EventsTransportSSE

	client, err := marathon.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create a client for marathon, error: %s", err)
	}

	// Register for events
	updates, err := client.AddEventsListener(marathon.EventIDStatusUpdate)
	if err != nil {
		log.Fatalf("Failed to register for events: %s", err)
	}

	for {
		select {
		case event := <-updates:
			statusUpdate := event.Event.(*marathon.EventStatusUpdate)
			appId := strings.Replace(statusUpdate.AppID, "/", "_", -1)
			appId = strings.Replace(appId, ",", "_", -1)

			log.Printf("Received event: %s (id: %s)", statusUpdate.TaskStatus, appId)
			stat.Increment(fmt.Sprintf("%s.%s", statusUpdate.TaskStatus, appId))
		}
	}
}
