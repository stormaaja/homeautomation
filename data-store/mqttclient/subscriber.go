package mqttclient

import (
	"context"
	"encoding/json"
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

const (
	clientID         = "data-store"
	measurementTopic = "measurements"
	temperatureTopic = "temperatures"
)

var mqttMsgChan = make(chan mqtt.Message)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	mqttMsgChan <- msg
}

type MqttMessage struct {
	MeasurementType string
	MeasurmentKey   string
	Measurement     store.Measurement
}

func processMsg(ctx context.Context, input <-chan mqtt.Message, memoryStore *store.MemoryStore) chan mqtt.Message {
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
				if msg.Topic() == measurementTopic {
					mqttMessage := MqttMessage{}
					err := json.Unmarshal(msg.Payload(), &mqttMessage)
					if err != nil {
						log.Printf("Error parsing MQTT message: %v\n", err)
					} else {
						memoryStore.SetMeasurement(mqttMessage.MeasurementType, mqttMessage.MeasurmentKey, mqttMessage.Measurement)
					}
				} else if msg.Topic() == temperatureTopic {
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

func Subscribe(broker string, memoryStore *store.MemoryStore) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
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
		finalChan := processMsg(ctx, mqttMsgChan, memoryStore)
		for range finalChan {
			// just consuming these for now
		}
	}()

	token := client.Subscribe(temperatureTopic, 1, nil)
	token.Wait()
	log.Printf("Subscribed to topic: %s\n", temperatureTopic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	cancel()

	log.Println("Unsubscribing and disconnecting...")
	client.Unsubscribe(temperatureTopic)
	client.Disconnect(250)

	wg.Wait()
	log.Println("Goroutine terminated, exiting...")
}
