package go_cielo_conecta

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type ClientInterface interface {
	NewRequest(method, path string, body any) (*http.Request, error)
	NewRequestWithContext(ctx context.Context, method, path string, body any) (*http.Request, error)
	Send(req *http.Request, body any) error

	CreatePayment(orderId string, amount float64, productId uint) SaleInterface
	SharedLibrary(terminalID string, subMerchantId ...string) (map[string]any, error)
	GetPaymentBy(param GetParam, query string, transactionDate ...time.Time) (*Sale, error)

	Close()
	SetLogger(slog *slog.Logger)
}

// NewClient creates a new instance of the Cielo Conecta API client. It requires a Merchant struct containing
// the merchant's credentials and an Environment struct with the necessary API URLs. The function validates the input
// parameters, initializes the client, retrieves an access token, and starts a background goroutine to refresh the token
// as needed. If any required fields are missing or if there is an error during token retrieval, it returns an error.
// Otherwise, it returns a pointer to the initialized Client instance.
func NewClient(m Merchant, env Environment) (ClientInterface, error) {
	if m.ID == "" || m.Secret == "" || env.APIUrl == "" || env.OAuthURL == "" || env.APIQueryUrl == "" || env.ParamsURL == "" {
		return nil, errors.New("merchantId, merchantSecret and environment fields are required")
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := Client{
		Mutex:    sync.Mutex{},
		Client:   &http.Client{},
		merchant: m,
		env:      env,
		token:    nil,
		cancel:   cancel,
	}

	c.DefaultLogger()

	err := c.getToken()
	if err != nil {
		return nil, err
	}

	c.writeLog(fmt.Sprintf("Successfully got access_token, expires in %s", (c.token.ExpiresIn * time.Second).String()))

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.refreshToken(ctx)
	}()

	return &c, nil
}

// NewRequest creates a new HTTP request with the specified method, path, and body.
// If the body is not nil, it encodes it as JSON and includes it in the request.
//
// The function returns the created HTTP request or an error if there was an issue encoding the body.
func (c *Client) NewRequest(method, path string, body any) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	return http.NewRequest(method, path, &buf)
}

// NewRequestWithContext creates a new HTTP request with the specified context, method, path, and body.
// If the body is not nil, it encodes it as JSON and includes it in the request.
//
// The function returns the created HTTP request or an error if there was an issue encoding the body.
func (c *Client) NewRequestWithContext(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(b)
	}

	return http.NewRequestWithContext(ctx, method, path, buf)
}

// Send sends an HTTP request and decodes the response into the provided variable.
// It sets the necessary headers for authentication and content type, and logs the request and response.
//
// If the response status code indicates an error (not in the 200-299 range), it attempts to decode the error response
// and returns it. If there is an issue decoding the response, it returns an error with the status code and decoding error.
// If the request is successful, it decodes the response body into the provided variable.
func (c *Client) Send(req *http.Request, v any) error {
	if v == nil {
		return nil
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("User-Agent", "go-cielo-conecta-client/1.0")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))

	resp, err := c.Client.Do(req)
	c.logger(req, resp)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var errResp = []ErrorResponse{{Response: resp}}

		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return fmt.Errorf("unable to decode JSON response: code=%d error=%w", resp.StatusCode, err)
		}

		return errResp[0]
	}

	return json.NewDecoder(resp.Body).Decode(v)
}
