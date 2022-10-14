package cfp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type SpeakerClient struct {
	Client   *http.Client
	endpoint string
}

func NewSpeakerClient(endpoint string, client *http.Client) *SpeakerClient {
	if client == nil {
		client = http.DefaultClient
	}

	return &SpeakerClient{
		Client:   client,
		endpoint: endpoint,
	}
}

func (c *SpeakerClient) Create(ctx context.Context, path string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", c.endpoint, path), bytes.NewBuffer(body))
	if err != nil {
		return &Error{Reason: ErrCreateRequest, Err: err}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return &Error{Reason: ErrMakeRequest, Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return &Error{Reason: ErrCreateSpeaker, Err: fmt.Errorf("error creating speaker: %d", resp.StatusCode)}
	}

	return nil
}

func (c *SpeakerClient) Update(ctx context.Context, path, id string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s%s/%s", c.endpoint, path, id), bytes.NewBuffer(body))

	if err != nil {
		return &Error{Reason: ErrCreateRequest, Err: err}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return &Error{Reason: ErrMakeRequest, Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return &Error{Reason: ErrUpdateSpeaker, Err: fmt.Errorf("error updating speaker: %s", resp.Status)}
	}

	return nil
}

func (c *SpeakerClient) Get(ctx context.Context, path, id string) ([]byte, error) {
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
		return nil, &Error{Reason: ErrFetchSpeaker, Err: fmt.Errorf("error getting speaker: %s", resp.Status)}
	}

	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &Error{Reason: ErrFetchSpeaker, Err: fmt.Errorf("error reading response: %w", err)}
	}

	return payload, nil
}

func (c *SpeakerClient) Delete(ctx context.Context, path, id string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s%s/%s", c.endpoint, path, id), nil)
	if err != nil {
		return &Error{Reason: ErrCreateRequest, Err: err}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return &Error{Reason: ErrMakeRequest, Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return &Error{Reason: ErrDeleteSpeaker, Err: fmt.Errorf("error deleting speaker: %s", resp.Status)}
	}

	return nil
}
