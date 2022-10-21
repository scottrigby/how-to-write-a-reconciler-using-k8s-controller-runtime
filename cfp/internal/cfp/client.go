package cfp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	SpeakerPath  = "/api/speakers"
	ProposalPath = "/api/proposals"
)

type Client struct {
	Client   *http.Client
	endpoint string
}

func NewClient(endpoint string, client *http.Client) (*Client, error) {
	_, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		Client:   client,
		endpoint: endpoint,
	}, nil
}

func (c *Client) Create(ctx context.Context, path string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.endpoint, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, &Error{Reason: ErrCreateRequest, Err: err}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, &Error{Reason: ErrMakeRequest, Err: err}
	}

	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, computeError(fmt.Errorf("error reading response: %w", err), path, http.MethodGet)
	}

	// return the response body which will contain important error info
	if resp.StatusCode != http.StatusOK {
		return nil, computeError(fmt.Errorf("create error: %s: %s", payload, resp.Status), path, http.MethodPost)
	}

	return payload, nil
}

func (c *Client) Update(ctx context.Context, path, id string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s%s/%s", c.endpoint, path, id), bytes.NewBuffer(body))

	if err != nil {
		return nil, &Error{Reason: ErrCreateRequest, Err: err}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, &Error{Reason: ErrMakeRequest, Err: err}
	}

	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, computeError(fmt.Errorf("error reading response: %w", err), path, http.MethodGet)
	}

	// return the response body which will contain important error info
	if resp.StatusCode != http.StatusOK {
		return nil, computeError(fmt.Errorf("update error: %s: %s", payload, resp.Status), path, http.MethodPut)
	}

	return payload, nil
}

func (c *Client) Get(ctx context.Context, path, id string) ([]byte, error) {
	// Get by id
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s/%s", c.endpoint, path, id), nil)
	if err != nil {
		return nil, &Error{Reason: ErrCreateRequest, Err: err}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, &Error{Reason: ErrMakeRequest, Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, computeError(fmt.Errorf("get error: %s", resp.Status), path, http.MethodGet)
	}

	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, computeError(fmt.Errorf("error reading response: %w", err), path, http.MethodGet)
	}

	// return the response body which will contain important error info
	if resp.StatusCode != http.StatusOK {
		return nil, computeError(fmt.Errorf("get error: %s: %s", payload, resp.Status), path, http.MethodGet)
	}

	return payload, nil
}

func (c *Client) Delete(ctx context.Context, path, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s%s/%s", c.endpoint, path, id), nil)
	if err != nil {
		return &Error{Reason: ErrCreateRequest, Err: err}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return &Error{Reason: ErrMakeRequest, Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return computeError(fmt.Errorf("error deleting speaker: %s", resp.Status), path, http.MethodDelete)
	}

	return nil
}

func computeError(err error, path string, method string) error {
	if err == nil {
		return nil
	}

	switch {
	case method == http.MethodGet && path == SpeakerPath:
		return &Error{Reason: ErrFetchSpeaker, Err: err}
	case method == http.MethodGet && path == ProposalPath:
		return &Error{Reason: ErrFetchProposal, Err: err}
	case method == http.MethodPost && path == SpeakerPath:
		return &Error{Reason: ErrCreateSpeaker, Err: err}
	case method == http.MethodPost && path == ProposalPath:
		return &Error{Reason: ErrCreateProposal, Err: err}
	case method == http.MethodPut && path == SpeakerPath:
		return &Error{Reason: ErrUpdateSpeaker, Err: err}
	case method == http.MethodPut && path == ProposalPath:
		return &Error{Reason: ErrUpdateProposal, Err: err}
	case method == http.MethodDelete && path == SpeakerPath:
		return &Error{Reason: ErrDeleteSpeaker, Err: err}
	case method == http.MethodDelete && path == ProposalPath:
		return &Error{Reason: ErrDeleteProposal, Err: err}
	default:
		return &Error{Reason: ErrUnknown, Err: err}
	}
}
