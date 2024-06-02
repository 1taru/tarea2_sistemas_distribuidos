package main

import (
	"encoding/json"
	"fmt"
	"kafka-processing-system/model"
	"log"
	"net/http"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	lastSentEmail model.Solicitud
	mu            sync.RWMutex
	queue         chan model.Solicitud
)

func main() {
	// Inicializar la cola
	queue = make(chan model.Solicitud, 10)

	go consumeKafkaMessages()
	go processQueue()

	http.HandleFunc("/last", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lastSentEmail)
	})

	fmt.Println("Server is listening on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func consumeKafkaMessages() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "notificador",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{
		"solicitudes-recibido",
		"solicitudes-preparando",
		"solicitudes-entregando",
		"solicitudes-finalizado",
	}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %s", err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		var solicitud model.Solicitud
		err = json.Unmarshal(msg.Value, &solicitud)
		if err != nil {
			log.Printf("Failed to unmarshal message: %s", err)
			continue
		}

		queue <- solicitud
	}
}

func processQueue() {
	for solicitud := range queue {
		// Convertir la solicitud a JSON para imprimir
		jsonData, err := json.Marshal(solicitud)
		if err != nil {
			log.Printf("Failed to marshal solicitud to JSON: %s", err)
		} else {
			log.Printf("Processing solicitud: %s", string(jsonData))
		}

		// Actualizar el último correo enviado
		//mu.Lock()
		//lastSentEmail = solicitud
		//mu.Unlock()
		log.Printf(":v")

		// Actualizar el último correo enviado
		mu.Lock()
		lastSentEmail = solicitud
		mu.Unlock()
		// Añadir retraso para evitar la saturación del servidor de correo
		//time.Sleep(1 * time.Millisecond)
	}
}

/*
package main

import (
	"encoding/json"
	"fmt"
	"kafka-processing-system/model"
	"log"
	"net/http"
	"net/smtp"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	lastSentEmail model.Solicitud
	mu            sync.RWMutex
	queue         chan model.Solicitud
)

func sendEmail(to, subject, body string) error {
	from := "distritest681@gmail.com"
	password := "bxxp mcpx khty iqvk"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	return err
}

func main() {
	// Inicializar la cola
	queue = make(chan model.Solicitud, 10)

	go consumeKafkaMessages()
	go processQueue()

	http.HandleFunc("/last", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lastSentEmail)
	})

	fmt.Println("Server is listening on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func consumeKafkaMessages() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "notificador",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{
		"solicitudes-recibido",
		"solicitudes-preparando",
		"solicitudes-entregando",
		"solicitudes-finalizado",
	}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topics: %s", err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		var solicitud model.Solicitud
		err = json.Unmarshal(msg.Value, &solicitud)
		if err != nil {
			log.Printf("Failed to unmarshal message: %s", err)
			continue
		}

		queue <- solicitud
	}
}

func processQueue() {
	for solicitud := range queue {
		// Convertir la solicitud a JSON para imprimir
		jsonData, err := json.Marshal(solicitud)
		if err != nil {
			log.Printf("Failed to marshal solicitud to JSON: %s", err)
		} else {
			log.Printf("Processing solicitud: %s", string(jsonData))
		}

		subject := fmt.Sprintf("Actualización de Solicitud: %s", solicitud.Estado)
		body := fmt.Sprintf("Tu solicitud %s ha cambiado al estado: %s", solicitud.ID, solicitud.Estado)
		err = sendEmail(solicitud.Correo, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %s", err)
		} else {
			log.Printf("Email sent for solicitud %s with new state %s", solicitud.ID, solicitud.Estado)

			// Actualizar el último correo enviado
			//mu.Lock()
			//lastSentEmail = solicitud
			//mu.Unlock()
		}
		// Actualizar el último correo enviado
		mu.Lock()
		lastSentEmail = solicitud
		mu.Unlock()
		// Añadir retraso para evitar la saturación del servidor de correo
		//time.Sleep(1 * time.Millisecond)
	}
}

*/
