package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Initialize Redis client (should ideally be done once outside this function)
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis default address
	})

	// Channels for communication between goroutines
	g1Ticks := make(chan int)
	g2Ticks := make(chan int)
	g1Alerts := make(chan bool)

	// Start goroutines for G1 and G2
	go g1Generator(client, g1Ticks, g1Alerts)
	go g2Generator(g2Ticks, g1Alerts)

	// Main loop to simulate ticks
	g1Ticker := time.NewTicker(1 * time.Second)
	defer g1Ticker.Stop()

	g2Ticker := time.NewTicker(1 * time.Second)
	defer g2Ticker.Stop()

	var g2TickCount int

	for tick := 1; ; tick++ {
		select {
		case <-g1Ticker.C:
			fmt.Printf("Tick %d - G1\n", tick)
			g1Ticks <- tick // Send tick to G1

		case <-g2Ticker.C:
			g2TickCount++
			fmt.Printf("Tick %d - G2\n", tick)

			// Send tick to G2 and reset count every 7 ticks
			g2Ticks <- tick
			if g2TickCount == 7 {
				g2TickCount = 0
				g1Alerts <- true // Send alert to G1
			}
		}
	}
}

func g1Generator(client *redis.Client, g1Ticks <-chan int, g1Alerts <-chan bool) {
	var g1TickCount int

	for {
		select {
		case tick := <-g1Ticks:
			g1TickCount++
			log.Printf("G1 generated tick %d\n", tick)

			// Check if G1 should send data to Redis (every 3 ticks)
			if g1TickCount%3 == 0 {
				sendDataToRedis(client, "G1", tick)
			}

		case <-g1Alerts:
			// Handle alert from G2
			tick := <-g1Ticks
			log.Printf("Received alert from G2 to send data at tick %d\n", tick)
			sendDataToRedis(client, "G1", tick)
		}
	}
}

func g2Generator(g2Ticks <-chan int, g1Alerts chan<- bool) {
	for {
		select {
		case tick := <-g2Ticks:
			log.Printf("G2 generated tick %d\n", tick)
		}
	}
}

func sendDataToRedis(client *redis.Client, goroutine string, tick int) {
	// Prepare data to store in Redis
	data := fmt.Sprintf("Data sent to Redis by %s at tick %d, generated at %s", goroutine, tick, time.Now().UTC().Format(time.RFC3339))

	// Construct Redis key based on tick
	key := fmt.Sprintf("%s_%d_%d", goroutine, tick, tick/3)

	// Store data in Redis with the constructed key
	err := client.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		log.Printf("Failed to write data to Redis for key %s: %v\n", key, err)
		return
	}

	// Log success message
	log.Printf("Data written to Redis by %s with key %s: %s\n", goroutine, key, data)
}
