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

const (
	httpSecure      string = "https"
	websocketSecure string = "wss"
)

const (
	RequestHeaderRetryAttempts string = "X-Attempt-Count"
)

type client struct {
	httpClientForZendesk *http.Client
	httpClientForZopim   *http.Client
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

	var latestAttemptError error

	for attempts < maxAttempts {
		attempts++

		// Set the current attempt count in the request header, so that consumers of the library can identify when a request is being retried
		request.Header.Set(RequestHeaderRetryAttempts, fmt.Sprintf("%d", attempts))

		time.Sleep(time.Duration(retryAfter) * time.Second)

		latestAttemptError = c.do(request, target)
		if latestAttemptError == nil {
			return nil
		}

		// Check to see if the error is a network error
		networkErr, ok := latestAttemptError.(*url.Error)
		if ok {
			// If the error is determined temporary, retry
			if networkErr.Temporary() {
				continue
			}

			return networkErr
		}

		zendeskErr, ok := latestAttemptError.(*Error)
		if !ok {
			return latestAttemptError
		}

		if zendeskErr.Response.StatusCode != http.StatusTooManyRequests {
			return latestAttemptError
		}

		// Check for a "retry-after" header and then continue
		retryAfterString := zendeskErr.Response.Header.Get("Retry-After")
		if retryAfterString != "" {
			var err error

			retryAfter, err = strconv.ParseInt(retryAfterString, 10, 64)
			if err != nil {
				return err
			}
		}

		continue
	}

	return latestAttemptError
}

func (c *client) do(request *http.Request, target any) error {
	if request.URL.Host == "" {
		request.URL.Host = fmt.Sprintf("%s.zendesk.com", c.subDomain)
	}

	request.URL.Scheme = httpSecure
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

	clientToUse := c.httpClientForZendesk
	if strings.Contains(request.Host, "zopim.com") {
		clientToUse = c.httpClientForZopim
	}

	response, err := clientToUse.Do(request)
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
			Response: response,
		}

		// There are times where Zendesk will report an error, but not provide a json encoded body response - handle these by
		// providing the body bytes directly along with what the content type header is
		contentType := response.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			responseErr.Description = fmt.Sprintf("encountered error - response content is '%s', not JSON", contentType)
			responseErr.Message = string(bodyBytes)

			return responseErr
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

func (c *client) ZendeskRequest(request *http.Request, target any) error {
	c.zendeskAuth.AddZendeskAuthentication(request)

	if request.Method == http.MethodGet {
		return c.doWithRetry(request, target)
	}

	return c.do(request, target)
}

func (c *client) ZendeskGetRequest(ctx context.Context, url string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	c.zendeskAuth.AddZendeskAuthentication(request)

	return c.httpClientForZendesk.Do(request)
}

/*
NOTE: The RealTimeChat API uses a different URL to the LiveChat API/Zendesk API

https://developer.zendesk.com/api-reference/live-chat/real-time-chat-api/rest/
*/
func (c *client) RealTimeChatRequest(request *http.Request, target any) error {
	if request.URL.Host == "" {
		request.URL.Host = "rtm.zopim.com"
	}

	request.URL.Scheme = httpSecure
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", c.userAgent)

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

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
				if zdError.Response.StatusCode == http.StatusUnauthorized {
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

func (c *client) ChatRequest(request *http.Request, target any) error {
	if request.URL.Host == "" {
		request.URL.Host = "www.zopim.com"
	}

	request.URL.Scheme = httpSecure
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", c.userAgent)

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

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
				if zdError.Response.StatusCode == http.StatusUnauthorized {
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

/*
NOTE: You must configure the API Client in Zendesk Chat to be "client_type": "confidential". By default it is "public"
and this setting reverts to "public" if you interact with the client via the Zendesk Chat Web UI.

https://developer.zendesk.com/documentation/live-chat/getting-started/auth/#implementing-an-oauth-authorization-flow

https://developer.zendesk.com/documentation/live-chat/getting-started/auth/#confidential-grant-types
*/
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

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://www.zopim.com/oauth2/token",
		strings.NewReader(data.Encode()),
	)
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
