package go_cielo_conecta

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type ClientInterface interface {
	Close()
	LogWithWriter(io.Writer)

	NewRequest(method, path string, body any) (*http.Request, error)
	NewRequestWithContext(ctx context.Context, method, path string, body any) (*http.Request, error)
	Send(req *http.Request, body any) error

	Authorization(s *Sale) (*Sale, error)
}

// NewClient ini
func NewClient(m Merchant, env Environment) (ClientInterface, error) {
	if m.ID == "" || m.Secret == "" || env.APIUrl == "" || env.OAuthURL == "" || env.APIQueryUrl == "" || env.ParamsURL == "" {
		return nil, errors.New("merchantId, merchantSecret and environment fields are required")
	}

	ctx, cancel := context.WithCancel(context.Background())

	var c = Client{
		Mutex:    sync.Mutex{},
		Client:   &http.Client{},
		merchant: m,
		env:      env,
		token:    nil,
		cancel:   cancel,
	}

	err := c.getToken()
	if err != nil {
		return nil, err
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.refreshToken(ctx)
	}()

	c.logStdOut() // Default logger is stdout, can be changed to LogWithWriter.

	return &c, nil
}

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

func (c *Client) logger(r *http.Request, resp *http.Response) {
	if c.Log == nil {
		return
	}

	var (
		requestDump  string
		responseDump string
	)

	if r != nil {
		requestDump = fmt.Sprintf("%s -> %s", r.Method, r.URL.String())
	}

	if resp != nil {
		// copy response body to avoid consuming it
		bodyCopy := bytes.NewBuffer(nil)
		_, err := io.Copy(bodyCopy, resp.Body)
		if err != nil {
			responseDump = fmt.Sprintf("status=%s, error_copying_body=%v", resp.Status, err)
		} else {
			responseDump = fmt.Sprintf("status=%s, response=%s", resp.Status, bodyCopy.String())

			// reset original response body for further processing
			resp.Body = io.NopCloser(bodyCopy)
		}
	}

	_, _ = c.Log.Write([]byte(fmt.Sprintf("[CieloConecta] Request: %s \n[CieloConecta] Response: %s \n", requestDump, responseDump)))
}

func (c *Client) logStdOut() {
	c.Log = os.Stdout
}
func (c *Client) LogWithWriter(w io.Writer) {
	c.Log = w
}
