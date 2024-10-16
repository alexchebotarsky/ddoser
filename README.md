# DDOSER

## Usage of ddoser:

- -url string

  Target URL to make requests to.

- -rate int

  Amount of requests per second.

- -method string (Optional. Default: GET)

  HTTP method to use.

- -body string (Optional)

  Body to send with the request.

- -header string (Optional)

  Header to send with the request.

- -http-timeout duration (Optional. Default: 1s)

  HTTP client timeout.

## Example:

### GET

```
ddoser -url http://example.com -rate 10
```

### POST

```
ddoser -url http://example.com -rate 10 -method POST -body '{"key": "value"}' -header 'Content-Type: application/json'
```
