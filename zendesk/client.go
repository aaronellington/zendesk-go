package zendesk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type client struct {
	httpClient           *http.Client
	zendeskAuth          authentication
	chatCredentials      ChatCredentials
	chatToken            *chatToken
	chatMutex            *sync.Mutex
	subdomain            string
	requestPreProcessors []RequestPreProcessor
}

func (c *client) do(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "https"
	r.Header.Set("Accept", "application/json")
	r.Header.Set("User-Agent", "aaronellington/zendesk-go")

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

func (c *client) json(ctx context.Context, method string, path string, body io.Reader, target any, requestPreProcessor RequestPreProcessor) error {
	request, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	if err := requestPreProcessor.ProcessRequest(request); err != nil {
		return err
	}

	response, err := c.do(request)
	if err != nil {
		return err
	}

	if target != nil {
		defer response.Body.Close()

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		os.WriteFile("stuff.json", bodyBytes, 0o755)

		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) ZendeskRequest(ctx context.Context, method string, path string, body io.Reader, target any) error {
	return c.json(ctx, method, path, body, target, RequestPreProcessorFunc(func(r *http.Request) error {
		r.URL.Host = fmt.Sprintf("%s.zendesk.com", c.subdomain)
		c.zendeskAuth.AddZendeskAuthentication(r)

		return nil
	}))
}

func (c *client) ChatRequest(ctx context.Context, method string, path string, body io.Reader, target any) error {
	if err := c.getAccessToken(ctx); err != nil {
		return err
	}

	if err := c.json(ctx, method, path, body, target, RequestPreProcessorFunc(func(r *http.Request) error {
		if c.chatToken == nil {
			return errors.New("no token")
		}

		if c.chatToken.AccessToken == "" {
			return errors.New("blank token")
		}

		r.URL.Host = "www.zopim.com"
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.chatToken.AccessToken))

		return nil
	})); err != nil {
		return err
	}

	return nil
}

type chatToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (c *client) getAccessToken(ctx context.Context) error {
	if c.chatToken != nil {
		return nil
	}

	c.chatMutex.Lock()
	defer c.chatMutex.Unlock()

	if c.chatToken != nil {
		return nil
	}

	target := chatToken{}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.chatCredentials.ClientID)
	data.Set("client_secret", c.chatCredentials.ClientSecret)

	if err := c.json(
		ctx,
		http.MethodPost,
		"/oauth2/token",
		strings.NewReader(data.Encode()),
		&target,
		RequestPreProcessorFunc(func(r *http.Request) error {
			r.URL.Host = "www.zopim.com"
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			return nil
		}),
	); err != nil {
		return err
	}

	c.chatToken = &target

	return nil
}
