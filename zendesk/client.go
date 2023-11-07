package zendesk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type client struct {
	httpClient           *http.Client
	zendeskAuth          authentication
	chatCredentials      ChatCredentials
	chatToken            *chatToken
	chatMutex            *sync.Mutex
	subDomain            string
	userAgent            string
	requestPreProcessors []RequestPreProcessor
}

func (c *client) doWithRetry(request *http.Request, target any) error {
	attempts := 0
	maxAttempts := 3
	retryAfter := int64(0)

	for attempts < maxAttempts {
		attempts++

		time.Sleep(time.Duration(retryAfter) * time.Second)

		response, err := c.httpClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if response.StatusCode >= http.StatusBadRequest {
			if response.StatusCode == http.StatusTooManyRequests {
				// Check for a "retry-after" header and then continue
				retryAfterString := response.Header.Get("retry-after")
				if retryAfterString != "" {
					retryAfter, err = strconv.ParseInt(retryAfterString, 10, 64)
					if err != nil {
						return err
					}
				}

				continue
			}

			responseErr := &Error{
				StatusCode: response.StatusCode,
				Body:       bodyBytes,
			}

			if err := json.Unmarshal(bodyBytes, responseErr); err != nil {
				return err
			}

			return responseErr
		}

		if target != nil {
			if err := json.Unmarshal(bodyBytes, target); err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("unable to complete request after retries (%d attempts)", maxAttempts)
}

func (c *client) do(request *http.Request, target any) error {
	if request.URL.Host == "" {
		request.URL.Host = fmt.Sprintf("%s.zendesk.com", c.subDomain)
	}

	request.URL.Scheme = "https"
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", c.userAgent)

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	for _, requestPreProcessor := range c.requestPreProcessors {
		if err := requestPreProcessor.ProcessRequest(request); err != nil {
			return err
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= http.StatusBadRequest {

		responseErr := &Error{
			StatusCode: response.StatusCode,
			Body:       bodyBytes,
		}

		if err := json.Unmarshal(bodyBytes, responseErr); err != nil {
			return err
		}

		return responseErr
	}

	if target != nil {
		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) ZendeskRequest(request *http.Request, target any, autoRetry bool) error {
	c.zendeskAuth.AddZendeskAuthentication(request)

	if autoRetry {
		return c.doWithRetry(request, target)
	}

	return c.do(request, target)
}

func (c *client) LiveChatRequest(request *http.Request, target any) error {
	attempts := 0
	maxAttempts := 2

	for attempts < maxAttempts {
		attempts++

		if err := c.getAccessToken(request.Context()); err != nil {
			return err
		}

		if c.chatToken == nil {
			return errors.New("no token")
		}

		if c.chatToken.AccessToken == "" {
			return errors.New("blank token")
		}

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.chatToken.AccessToken))

		if err := c.do(request, target); err != nil {
			if zdError, ok := err.(*Error); ok {
				if zdError.StatusCode == http.StatusUnauthorized {
					// Clear out the token
					c.chatToken = nil

					continue
				}
			}

			return err
		}

		break
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

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.chatCredentials.ClientID)
	data.Set("client_secret", c.chatCredentials.ClientSecret)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.zopim.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	target := chatToken{}

	if err := c.do(request, &target); err != nil {
		return err
	}

	c.chatToken = &target

	return nil
}
