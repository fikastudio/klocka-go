package klocka

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
)

const defaultEndpoint = "https://klocka.dev/api"

type Client struct {
	apiKey    string
	apiSecret string
	endpoint  string
	hcl       *http.Client
}
type ClientOpts struct {
	Endpoint string
}

func NewClient(apiKey, apiSecret string) (*Client, error) {
	return NewClientWithOpts(apiKey, apiSecret, &ClientOpts{
		Endpoint: defaultEndpoint,
	})
}

func NewClientWithOpts(apiKey, apiSecret string, opts *ClientOpts) (*Client, error) {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		endpoint:  opts.Endpoint,
		hcl: &http.Client{
			Transport: gzhttp.Transport(http.DefaultTransport),
		},
	}, nil
}

func (cl *Client) CreateTask(ctx context.Context, spec TaskInput) (*Task, error) {
	b, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cl.endpoint+"/v1/task", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cl.apiKey, cl.apiSecret)
	req.Header.Set("content-type", "application/json")

	resp, err := cl.hcl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New("invalid status code returned")
		body, _ := io.ReadAll(resp.Body)

		return nil, &APIError{
			err:    err,
			body:   body,
			status: resp.StatusCode,
		}
	}

	var t Task
	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (cl *Client) UpdateTask(ctx context.Context, id string, spec TaskInput) (*Task, error) {
	b, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, cl.endpoint+"/v1/task/"+id, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cl.apiKey, cl.apiSecret)
	req.Header.Set("content-type", "application/json")

	resp, err := cl.hcl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.New("invalid status code returned")
		body, _ := io.ReadAll(resp.Body)

		return nil, &APIError{
			err:    err,
			body:   body,
			status: resp.StatusCode,
		}
	}

	var t Task
	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (cl *Client) DeleteTask(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, cl.endpoint+"/v1/task/"+id, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(cl.apiKey, cl.apiSecret)

	resp, err := cl.hcl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		err = errors.New("invalid status code returned")
		body, _ := io.ReadAll(resp.Body)

		return &APIError{
			err:    err,
			body:   body,
			status: resp.StatusCode,
		}
	}

	return nil
}

func (cl *Client) ListTasks(ctx context.Context, id string, opts *ListOpts) (*PaginatedResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cl.endpoint+"/v1/task?"+opts.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cl.apiKey, cl.apiSecret)
	req.Header.Set("content-type", "application/json")

	resp, err := cl.hcl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = errors.New("invalid status code returned")
		body, _ := io.ReadAll(resp.Body)

		return nil, &APIError{
			err:    err,
			body:   body,
			status: resp.StatusCode,
		}
	}

	var pr PaginatedResponse
	if err = json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (cl *Client) GetTask(ctx context.Context, id string) (*Task, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cl.endpoint+"/v1/task/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cl.apiKey, cl.apiSecret)
	req.Header.Set("content-type", "application/json")

	resp, err := cl.hcl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = errors.New("invalid status code returned")
		body, _ := io.ReadAll(resp.Body)

		return nil, &APIError{
			err:    err,
			body:   body,
			status: resp.StatusCode,
		}
	}

	var t Task
	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}

	return &t, nil
}
