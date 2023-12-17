package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const accuralEndpoint = "/api/orders/"

type AccuralClient struct {
	cl  http.Client
	url string
}

func New(addr string) *AccuralClient {
	return &AccuralClient{
		cl:  *http.DefaultClient,
		url: addr + accuralEndpoint,
	}
}

func (c *AccuralClient) UpdateOrderStatus(ctx context.Context, orderID string) (*AccuralResponse, error) {
	uri := c.makeURI(orderID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, nil
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("do request, err=%w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read from body, err=%w", err)
	}

	accural := &AccuralResponse{}
	if err := json.Unmarshal(data, accural); err != nil {
		return nil, fmt.Errorf("accural unmarshal err=%w", err)
	}

	return accural, nil
}

var tryingIntervals []time.Duration = []time.Duration{
	time.Millisecond * 1000,
	time.Millisecond * 3000,
	time.Millisecond * 5000,
}

func (c *AccuralClient) doRequest(req *http.Request) (*http.Response, error) {
	var joinedError error
	maxTryingsNum := len(tryingIntervals)

	for trying := 0; trying <= maxTryingsNum; trying++ {
		if resp, err := c.cl.Do(req); err != nil {
			if trying < maxTryingsNum {
				joinedError = errors.Join(joinedError, err)
				time.Sleep(tryingIntervals[trying])
			}
		} else {
			return resp, nil
		}
	}

	return nil, fmt.Errorf("request trying limit exceeded, errs=%w", joinedError)
}

func (c *AccuralClient) makeURI(orderId string) string {
	return c.url + orderId
}
