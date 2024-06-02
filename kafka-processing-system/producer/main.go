package main

import (
	"encoding/json"
	"fmt"
	"io"
	"kafka-processing-system/model"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var producer *kafka.Producer
var topic = "solicitudes"
var currentID int
var mu sync.Mutex

func main() {
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	http.HandleFunc("/send", handleSend)

	fmt.Println("Server is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var solicitud model.Solicitud
	err = json.Unmarshal(body, &solicitud)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
		return
	}

	// Generate and assign an incremental ID
	mu.Lock()
	currentID++
	solicitud.ID = fmt.Sprintf("%d", currentID)
	mu.Unlock()

	// Launch a goroutine to handle the Kafka message production asynchronously
	func(solicitud model.Solicitud) {
		solicitud.Estado = "recibido"
		solicitud.Time = time.Now()
		solicitud.Correo = "distritest681@gmail.com"

		value, err := json.Marshal(solicitud)
		if err != nil {
			log.Printf("Failed to marshal JSON: %s", err)
			return
		}

		// Print the JSON message
		log.Printf("Sending JSON: %s", string(value))

		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value,
		}, nil)
		if err != nil {
			log.Printf("Failed to produce message: %s", err)
			return
		}

		// Wait for the message to be delivered
		select {
		case e := <-producer.Events():
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		case <-time.After(10 * time.Second):
			log.Printf("Delivery timeout\n")
		}
	}(solicitud)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message sent to Kafka"))
}
