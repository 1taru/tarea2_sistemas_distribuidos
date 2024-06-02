package main

import (
	"encoding/json"
	"fmt"
	"kafka-processing-system/model"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var estados = []string{"recibido", "preparando", "entregando", "finalizado"}

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "procesador",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{"solicitudes"}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %s", err)
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close()

	var wg sync.WaitGroup

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		wg.Add(1)
		go func(msg *kafka.Message) {
			defer wg.Done()
			go handleSolicitud(msg, p)
		}(msg)
	}
	wg.Wait()
}

func handleSolicitud(msg *kafka.Message, p *kafka.Producer) {
	var solicitud model.Solicitud
	err := json.Unmarshal(msg.Value, &solicitud)
	if err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
		return
	}

	for _, estado := range estados {
		solicitud.Estado = estado
		solicitud.Time = time.Now()

		value, err := json.Marshal(solicitud)
		if err != nil {
			log.Printf("Failed to marshal message: %s", err)
			return
		}

		topic := fmt.Sprintf("solicitudes-%s", solicitud.Estado)
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value,
		}, nil)

		if err != nil {
			log.Printf("Failed to produce message: %s", err)
			return
		}

		// Ensure the message is delivered
		e := <-p.Events()
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}

		//time.Sleep(1 * time.Second) // Simulate processing time
	}
}
