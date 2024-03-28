package zendesk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type client struct {
	userAgent            string
	httpClient           *http.Client
	authentication       Authentication
	subDomain            string
	requestPreProcessors []RequestPreProcessor
}

func (c *client) request(
	ctx context.Context,
	method string,
	path string,
	body any,
	target any,
) error {
	var bodyReader io.Reader = http.NoBody

	if body != nil && body != http.NoBody {
		payloadBytes, err := json.MarshalIndent(body, "", "\t")
		if err != nil {
			panic(err)
		}

		bodyReader = bytes.NewReader(payloadBytes)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		method,
		path,
		bodyReader,
	)
	if err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if target != nil {
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(target); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) do(request *http.Request) (*http.Response, error) {
	c.requestFixer(request)

	for _, requestPreProcessor := range c.requestPreProcessors {
		if err := requestPreProcessor.ProcessRequest(request); err != nil {
			return nil, err
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= http.StatusBadRequest {
		defer response.Body.Close()

		zendeskError := Error{
			Response: response,
		}

		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&zendeskError); err != nil {
			return nil, err
		}

		return nil, zendeskError
	}

	return response, nil
}

func (c *client) requestFixer(request *http.Request) {
	zendeskHost := fmt.Sprintf("%s.zendesk.com", c.subDomain)

	if request.URL.Scheme == "" {
		request.URL.Scheme = "https"
	}

	if request.Header.Get("User-Agent") == "" {
		request.Header.Set("User-Agent", c.userAgent)
	}

	if request.URL.Host == "" {
		request.URL.Host = zendeskHost
		request.Host = zendeskHost
	}

	if request.Header.Get("Accept") == "" {
		request.Header.Set("Accept", "application/json")
	}

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	if request.URL.Host == zendeskHost {
		c.authentication.AddZendeskAuthentication(request)
	}
}
