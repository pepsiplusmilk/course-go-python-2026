package urls

import (
	"context"
	"log"
	"sync"
)

type ResponseCodeAggregator struct {
	client HttpClient
}

func NewURLAggregator(client HttpClient) *ResponseCodeAggregator {
	return &ResponseCodeAggregator{client: client}
}

func (u *ResponseCodeAggregator) process(ctx context.Context, url string, result *SimpleAggregatorResult) {
	type info struct {
		statusCode     int
		serverAnswered bool
		err            error
	}

	ch := make(chan info, 1)

	go func() {
		resp, err := u.client.Get(url)

		var inf info
		inf.err = err

		if resp != nil {
			inf.statusCode = resp.StatusCode
			inf.serverAnswered = true
			_ = resp.Body.Close()
		}

		ch <- inf
	}()

	select {
	case <-ctx.Done():
		return
	case inf := <-ch:
		if inf.err != nil {
			log.Printf("failed to get response from the server %q. got error: %v", url, inf.err)
			return
		}

		if inf.serverAnswered {
			result.count(inf.statusCode)
		}
	}
}

func (u *ResponseCodeAggregator) Aggregate(ctx context.Context, urls <-chan string) (AggregationResult, error) {
	result := NewSimpleAggregatorResult()

	go func() {
		defer result.finish()

		var group sync.WaitGroup
		defer group.Wait()

		for {

			if ctx.Err() != nil {
				return
			}

			select {
			case <-ctx.Done():
				return
			case url, ok := <-urls:
				if !ok {
					return
				}

				if ctx.Err() != nil {
					return
				}

				group.Add(1)
				go func(url string) {
					defer group.Done()
					u.process(ctx, url, result)
				}(url)
			}
		}
	}()

	return result, nil
}
