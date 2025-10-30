package mqttclient

import (
	"context"
	"log"
	"os"
	"os/signal"
	"stormaaja/go-ha/data-store/store"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttMsgChan = make(chan mqtt.Message)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	mqttMsgChan <- msg
}

func processMsg(ctx context.Context, input <-chan mqtt.Message, topic string, memoryStore *store.MemoryStore) chan mqtt.Message {
	out := make(chan mqtt.Message)
	go func() {
		defer close(out)
		for {
			select {
			case msg, ok := <-input:
				if !ok {
					return
				}
				log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
				if msg.Topic() == topic {
					payload := string(msg.Payload())
					splitted := strings.Split(payload, ":")
					if len(splitted) != 2 {
						log.Printf("Invalid temperature payload: %s\n", payload)
					} else {
						clientId := strings.TrimSpace(splitted[0])
						value, err := strconv.ParseFloat(strings.TrimSpace(splitted[1]), 64)
						if err != nil {
							log.Printf("Parsing temperature from MQTT payload failed: %v\n", err)
						} else {
							measurement := store.Measurement{
								DeviceId:        clientId,
								MeasurementType: "temperature",
								Field:           "temperature",
								Value:           value,
								UpdatedAt:       time.Now(),
							}
							memoryStore.SetMeasurement("temperature", clientId, measurement)
						}
					}
				}
				out <- msg
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected to MQTT Broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v", err)
}

func Subscribe(clientId string, broker string, topic string, memoryStore *store.MemoryStore) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		finalChan := processMsg(ctx, mqttMsgChan, topic, memoryStore)
		for range finalChan {
			// just consuming these for now
		}
	}()

	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Printf("Subscribed to topic: %s\n", topic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	cancel()

	log.Println("Unsubscribing and disconnecting...")
	client.Unsubscribe(topic)
	client.Disconnect(250)

	wg.Wait()
	log.Println("Goroutine terminated, exiting...")
}
