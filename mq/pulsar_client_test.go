package mq

//todo pulsar 依赖 需要单独编译
//import (
//	"context"
//	"fmt"
//	"github.com/apache/pulsar/pulsar-client-go/pulsar"
//	"log"
//	"testing"
//)
//
//func TestPulsarProducer(t *testing.T) {
//	println("start producer pulsar")
//	// Instantiate a Pulsar client
//	client, err := pulsar.NewClient(pulsar.ClientOptions{
//		URL: "pulsar://localhost:6650",
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Use the client to instantiate a producer
//	producer, err := client.CreateProducer(pulsar.ProducerOptions{
//		Topic: "my-topic",
//		Name:  "my-name",
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	ctx := context.Background()
//
//	for i := 0; i < 10; i++ {
//		// Create a message
//		msg := pulsar.ProducerMessage{
//			Payload: []byte(fmt.Sprintf("message-%d", i)),
//		}
//		// Attempt to send the message
//		if err := producer.Send(ctx, msg); err != nil {
//			log.Fatal(err)
//		}
//		// Create a different message to send asynchronously
//		asyncMsg := pulsar.ProducerMessage{
//			Payload: []byte(fmt.Sprintf("async-message-%d", i)),
//		}
//
//		producer.SendAsync(ctx, asyncMsg, func(message pulsar.ProducerMessage, e error) {
//			if err != nil {
//				log.Fatal(err)
//			}
//			fmt.Printf("the %s successfully published", string(msg.Payload))
//		})
//
//	}
//}
//
//func TestPulsarConsumer(t *testing.T) {
//	// Instantiate a Pulsar client
//	client, err := pulsar.NewClient(pulsar.ClientOptions{
//		URL: "pulsar://localhost:6650",
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Use the client object to instantiate a consumer
//	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
//		Topic:            "my-golang-topic",
//		SubscriptionName: "sub-1",
//		Type:             pulsar.Exclusive,
//	})
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer consumer.Close()
//
//	ctx := context.Background()
//
//	// Listen indefinitely on the topic
//	for {
//		msg, err := consumer.Receive(ctx)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		// Do something with the message
//		processMessage(msg)
//
//		if err == nil {
//			// Message processed successfully
//			consumer.Ack(msg)
//		} else {
//			// Failed to process messages
//			consumer.Nack(msg)
//		}
//	}
//}
//
//func processMessage(msg pulsar.Message) {
//	log.Printf("msg topic:%s,Id:%d,time:%d,value:%s",
//		msg.Topic(), msg.ID(), msg.EventTime(), msg.Topic())
//}
