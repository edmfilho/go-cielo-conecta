package go_cielo_conecta

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type requestLog struct {
	Request  string `json:"request"`
	Response string `json:"response,omitempty"`

	Status string `json:"status"`
	Code   int    `json:"-"`
}

func (req requestLog) LogValue() slog.Value {
	if req.Response == "" {
		return slog.GroupValue(
			slog.String("request", req.Request),
			slog.String("status", req.Status),
		)
	}

	return slog.GroupValue(
		slog.String("request", req.Request),
		slog.String("response", req.Response),
		slog.String("status", req.Status),
	)
}

func (c *Client) logger(r *http.Request, resp *http.Response) {
	if c.log == nil {
		return
	}

	logger := requestLog{
		Request: fmt.Sprintf("%s %s", r.Method, r.URL.String()),
		Status:  resp.Status,
		Code:    resp.StatusCode,
	}

	if logger.Code < 200 || logger.Code > 299 {
		c.log.Error("Error executing the request", "result", logger)
		return
	}

	c.log.Info("Request performed successfully.", "result", logger)
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

func (c *Client) writeLog(message string) {
	if c.log == nil {
		return
	}

	c.log.Info(message)
}
