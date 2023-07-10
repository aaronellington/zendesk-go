package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type client struct {
	httpClient           *http.Client
	auth                 authentication
	subdomain            string
	requestPreProcessors []RequestPreProcessor
}

func (c *client) Do(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "https"
	r.URL.Host = fmt.Sprintf("%s.zendesk.com", c.subdomain)
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("User-Agent", "aaronellington/zendesk-go")

	if c.auth != nil {
		c.auth.AddZendeskAuthentication(r)
	}

	for _, requestPreProcessor := range c.requestPreProcessors {
		if err := requestPreProcessor.ProcessRequest(r); err != nil {
			return nil, err
		}
	}

	response, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= http.StatusBadRequest {
		defer response.Body.Close()

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		responseErr := &Error{
			StatusCode: response.StatusCode,
			Body:       bodyBytes,
		}

		if err := json.Unmarshal(bodyBytes, responseErr); err != nil {
			return nil, err
		}

		return nil, responseErr
	}

	return response, nil
}

func (c *client) jsonRequest(ctx context.Context, method string, path string, body io.Reader, target any) error {
	request, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return err
	}

	response, err := c.Do(request)
	if err != nil {
		return err
	}

	if target != nil {
		defer response.Body.Close()

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return err
		}
	}

	return nil
}
