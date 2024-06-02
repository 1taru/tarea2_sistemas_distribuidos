package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	latency   time.Duration
	startTime time.Time
	endTime   time.Time
}

type Solicitud struct {
	ID     string    `json:"id"`
	Estado string    `json:"estado"`
	Time   time.Time `json:"time"`
	Correo string    `json:"correo"`
}

func main() {
	baseURL := "http://localhost:8080/send"
	var wg sync.WaitGroup

	messageCounts := []int{10, 50, 100, 200, 500, 2000, 5000, 10000} // different message loads
	results := make(map[int][]Result)

	for _, count := range messageCounts {
		var res []Result
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				solicitud := Solicitud{
					ID:     fmt.Sprintf("solicitud-%d", i),
					Estado: "inicial",
					Time:   time.Now(),
					Correo: "distritest681@gmail.com",
				}

				data, err := json.Marshal(solicitud)
				if err != nil {
					fmt.Printf("Error marshaling JSON: %s\n", err)
					return
				}

				startTime := time.Now()
				resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(data))
				endTime := time.Now()

				if err != nil {
					fmt.Printf("HTTP request failed: %s\n", err)
					return
				}

				if resp.StatusCode != http.StatusOK {
					fmt.Printf("Failed to send message: %s\n", resp.Status)
				} else {
					body, _ := io.ReadAll(resp.Body)
					fmt.Printf("Response: %s\n", string(body))
				}

				resp.Body.Close()

				res = append(res, Result{
					latency:   endTime.Sub(startTime),
					startTime: startTime,
					endTime:   endTime,
				})
			}(i)
		}
		wg.Wait()
		results[count] = res
	}

	// Output results as CSV
	fmt.Println("MessageCount, AverageLatency(ms), Throughput(msg/sec)")
	for count, result := range results {
		var totalLatency time.Duration
		for _, res := range result {
			totalLatency += res.latency
		}
		averageLatency := totalLatency / time.Duration(len(result))
		throughput := float64(len(result)) / result[len(result)-1].endTime.Sub(result[0].startTime).Seconds()
		fmt.Printf("%d, %.2f, %.2f\n", count, averageLatency.Seconds()*1000, throughput)
	}
}
