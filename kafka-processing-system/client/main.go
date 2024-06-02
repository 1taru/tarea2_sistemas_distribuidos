package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kafka-processing-system/model"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	url := "http://localhost:8080/send"
	datasetPath := filepath.Join("..", "dataset.json")

	// Leer el archivo JSON
	jsonFile, err := os.Open(datasetPath)
	if err != nil {
		log.Fatalf("Failed to open dataset file: %s", err)
	}
	defer jsonFile.Close()

	// Leer el contenido del archivo
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Failed to read dataset file: %s", err)
	}

	// Definir una variable para almacenar las solicitudes
	var solicitudes []model.Solicitud

	// Deserializar el JSON en la variable solicitudes
	err = json.Unmarshal(byteValue, &solicitudes)
	if err != nil {
		log.Fatalf("Failed to unmarshal dataset file: %s", err)
	}

	// Enviar cada solicitud al servidor con un ID incremental
	for i, solicitud := range solicitudes {

		// Establecer el ID incremental
		solicitud.ID = fmt.Sprintf("%d", i+1)

		jsonData, err := json.Marshal(solicitud)
		if err != nil {
			log.Printf("Failed to marshal message: %s", err)
			continue
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Failed to send request: %s", err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			log.Printf("Solicitud %s enviada exitosamente", solicitud.ID)
		} else {
			log.Printf("Failed to send request: received status code %d", resp.StatusCode)
		}
	}
}
