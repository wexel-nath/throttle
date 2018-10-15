package main

import (
	"fmt"
	"net/http"
	"time"
	"sync"

	"github.com/wexel-nath/throttle"
)

const (
	numWorkers = 5
	numJobs    = 1000

	url = "http://localhost:23456/default"
)

type Job struct {
	ID int
}

var (
	jobs = make(chan Job, numJobs)
	wg   = sync.WaitGroup{}
)

func main() {
	createJobs()

	// start workers
	for i := 1; i <= numWorkers; i++ {
		time.Sleep(50 * time.Millisecond) // stagger workers
		wg.Add(1)
		go worker(i)
	}

	wg.Wait()
}

func createJobs() {
	for i := 1; i <= numJobs; i++ {
		jobs<- Job{
			ID: i,
		}
	}
}

func worker(workerID int) {
	throttler := throttle.NewThrottler(throttle.Config{}) // use the defaults

	for job := range jobs {
		time.Sleep(throttler.Duration())
		fmt.Printf("Worker: %d processing job: %d\n", workerID, job.ID)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		client := http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		code := resp.StatusCode
		if code != http.StatusOK {
			fmt.Printf("code: %d. increasing throttle\n", code)
			throttler.Increase()
		} else {
			fmt.Printf("code: %d. decreasing throttle\n", code)
			throttler.Decrease()
		}
	}
	wg.Done()
}
