# throttle [![Build Status](https://travis-ci.com/wexel-nath/throttle.svg?branch=master)](https://travis-ci.com/wexel-nath/throttle)
A simple package to dynamically throttle outgoing http requests

# Installation
```
go get -v github.com/wexel-nath/throttle
```

# Usage
Throttling a single worker pulling from a queue:
```go
func worker(jobs chan struct{}) {
    throttler := throttle.NewThrottler(throttle.Config{}) // use the defaults

    for job := range jobs {
        time.Sleep(throttler.Duration())

        request, _ := http.NewRequest("GET", "http://example.com/path", nil)

        client := http.Client{}
        resp, _ := client.Do(request)

        // increase or decrease the throttler
        code := resp.StatusCode
        if code != http.StatusOK {
            throttler.Increase()
        } else {
            throttler.Decrease()
        }
    }
}
```

# Examples
Run the server:
```
go run examples/server/server.go
```

Run the client:
```
go run examples/client/client.go
```
