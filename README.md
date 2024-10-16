# DDOSER

## Usage of ddoser:

- -url string

  Target URL to make requests to.

- -rate int

  Amount of requests per second.

- -method string (Optional. Default: "GET")

  HTTP method to use.

- -body string (Optional. Default: "")

  Body to send with the request.

- -http-timeout duration (Optional. Default: 1s)

  HTTP client timeout.

## Example:

```
ddoser -url http://example.com -rate 10
```
