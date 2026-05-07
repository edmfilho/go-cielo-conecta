package go_cielo_conecta

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type logData struct {
	URL    string `json:"url"`
	Method string `json:"method"`
	Status string `json:"status"`
	Body   string `json:"body,omitempty"`
}

func (req logData) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("url", req.URL),
		slog.String("method", req.Method),
		slog.String("status", req.Status),
		slog.String("body", req.Body),
	)
}

func (c *Client) logger(r *http.Request, resp *http.Response) {
	if c.log == nil {
		return
	}

	logger := logData{
		URL:    r.URL.String(),
		Method: r.Method,
		Status: resp.Status,
	}

	if !strings.Contains(logger.URL, "parametersdownloadsandbox") {
		bodyBytes, err := readBody(resp)
		if err != nil {
			return
		}

		logger.Body = string(bodyBytes)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		c.log.Error("error executing the request", "request", logger)
		return
	}

	c.log.Info("request was successful", "request", logger)
}

func readBody(resp *http.Response) ([]byte, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	copiedBody := bodyBytes

	// Restore the original body for further processing
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return copiedBody, nil
}

func (c *Client) SetLogger(slog *slog.Logger) {
	c.log = slog.With("source", "cielo-conecta-client")
}

func (c *Client) DefaultLogger() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(time.Now().Format(time.RFC3339))
			}
			return a
		},
	}))

	c.log = l.With("source", "cielo-conecta-client")
}

func (c *Client) LogInfo(msg string, args ...any) {
	if c.log == nil {
		return
	}

	c.log.Info(msg, args...)
}

func (c *Client) LogError(msg string, args ...any) {
	if c.log == nil {
		return
	}

	c.log.Error(msg, args...)
}
