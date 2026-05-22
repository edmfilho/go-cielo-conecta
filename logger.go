package go_cielo_conecta

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const maxLogSize = 1024 * 100 // 102KB

type LogInfo struct {
	URL        string `json:"url"`
	Method     string `json:"method"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Body       []byte `json:"body,omitempty"`
}

func (l LogInfo) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("method", l.Method),
		slog.String("status", l.Status),
		slog.String("url", l.URL),
		slog.String("body", string(l.Body)),
	)
}

func (c *Client) logger(r *http.Request, resp *http.Response) {
	if c.log == nil {
		return
	}

	l := readBody(r, resp)

	if l.StatusCode < 200 || l.StatusCode > 299 {
		c.log.Error("error executing the request", "request", l)
		return
	}

	c.log.Info("request was successful", "request", l)
}

func readBody(r *http.Request, resp *http.Response) LogInfo {
	logInfo := LogInfo{
		URL:        r.URL.String(),
		Method:     r.Method,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       nil,
	}

	if r.Method == http.MethodGet && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return logInfo
	}

	content, _ := io.ReadAll(resp.Body)

	if int64(len(content)) > maxLogSize {
		resp.Body = io.NopCloser(bytes.NewBuffer(content))
		return logInfo
	}

	// Restore the original body for further processing
	resp.Body = io.NopCloser(bytes.NewBuffer(content))

	logInfo.Body = content
	return logInfo
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
