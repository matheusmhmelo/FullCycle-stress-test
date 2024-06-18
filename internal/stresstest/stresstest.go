package stresstest

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"sync"
	"time"
)

type testsConfig struct {
	mu sync.Mutex

	url         string
	requests    int
	concurrency int
	running     bool
}

type testsResult struct {
	requests  int
	successes int
	failures  map[string]int
}

func ExecuteTests(cmd *cobra.Command, _ []string) {
	config, err := initializeTests(cmd)
	if err != nil {
		newError(err)
		return
	}
	results := make(chan testsResult, config.concurrency)
	start := make(chan struct{})
	limitChan := make(chan struct{}, config.requests)
	ctx, cancel := context.WithCancel(context.Background())

	go config.limit(ctx, cancel, limitChan)

	fmt.Println("Preparing concurrency...")
	wg := &sync.WaitGroup{}
	for i := 0; i < config.concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			config.run(ctx, start, limitChan, results, i)
		}(i)
	}

	fmt.Println("Running tests...")

	startTime := time.Now()
	close(start)
	wg.Wait()
	endTime := time.Since(startTime)

	finalResult := testsResult{
		failures: map[string]int{},
	}
	for i := 0; i < config.concurrency; i++ {
		res := <-results
		finalResult.successes += res.successes
		finalResult.requests += res.requests

		for code, count := range res.failures {
			_, ok := finalResult.failures[code]
			if !ok {
				finalResult.failures[code] = count
			} else {
				finalResult.failures[code] = finalResult.failures[code] + count
			}
		}
	}

	fmt.Println("Tests finished in ", endTime.String())
	fmt.Println("Total Request: ", finalResult.requests)
	fmt.Println("Successes (HTTP 200): ", finalResult.successes)
	fmt.Println("Failures (By Status code): ", finalResult.failures)
}

func (cfg *testsConfig) run(ctx context.Context, start <-chan struct{}, limit chan<- struct{}, res chan<- testsResult, i int) {
	results := testsResult{
		failures: map[string]int{},
	}
	<-start

	for {
		select {
		case <-ctx.Done():
			res <- results
			return
		default:
			client := http.DefaultClient
			req, _ := http.NewRequestWithContext(ctx, "GET", cfg.url, nil)
			response, err := client.Do(req)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					limit <- struct{}{}
					continue
				}
				fmt.Println("Error while making request: ", err)
				results.requests++
				results.failures["unknown error"]++
				limit <- struct{}{}
				continue
			}

			if response.StatusCode == http.StatusOK {
				results.successes++
			} else {
				_, ok := results.failures[response.Status]
				if ok {
					results.failures[response.Status]++
				} else {
					results.failures[response.Status] = 1
				}
			}
			results.requests++
			limit <- struct{}{}
		}
	}
}

func (cfg *testsConfig) limit(ctx context.Context, cancel context.CancelFunc, limitChan <-chan struct{}) {
	var total int
	for {
		select {
		case <-ctx.Done():
			return
		case <-limitChan:
			total++
			if total >= cfg.requests {
				cancel()
				return
			}
		}
	}
}

func initializeTests(cmd *cobra.Command) (*testsConfig, error) {
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return nil, err
	}

	requests, err := cmd.Flags().GetInt("requests")
	if err != nil {
		return nil, err
	}

	concurrency, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return nil, err
	}

	return &testsConfig{
		url:         url,
		requests:    requests,
		concurrency: concurrency,
	}, nil
}

func newError(err error) {
	fmt.Println("Error during stress tests:", err)
}
