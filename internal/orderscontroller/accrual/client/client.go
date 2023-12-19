package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/internal/zlog"
	"io"
	"net/http"
	"time"
)

const accrualEndpoint = "/api/orders/"

type AccrualClient struct {
	cl  http.Client
	url string
}

func New(addr string) *AccrualClient {
	return &AccrualClient{
		cl:  *http.DefaultClient,
		url: addr + accrualEndpoint,
	}
}

func (c *AccrualClient) UpdateOrderStatus(ctx context.Context, orderID string) (*AccrualResponse, error) {
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

	switch resp.StatusCode {
	case succesHandlingStatusCode:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read from body, err=%w", err)
		}

		accrual := &AccrualResponse{}
		if err := json.Unmarshal(data, accrual); err != nil {
			return nil, fmt.Errorf("accrual unmarshal err=%w", err)
		}

		return accrual, nil
	case orderIsNotRegistredStatusCode:
		return nil, ErrOrderIsNotRegisted
	case requestLimitExcededStatusCode:
		return nil, ErrRequestLimitExceeded
	default:
		zlog.Logger.Errorf("Unknown status code from accrual system code=%d", resp.StatusCode)
		return nil, ErrUnknownStatusCode
	}
}

var tryingIntervals []time.Duration = []time.Duration{
	time.Millisecond * 100,
	time.Millisecond * 200,
	time.Millisecond * 300,
}

func (c *AccrualClient) doRequest(req *http.Request) (*http.Response, error) {
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

func (c *AccrualClient) makeURI(orderID string) string {
	return c.url + orderID
}
