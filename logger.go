package go_cielo_conecta

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) logger(r *http.Request, resp *http.Response) {
	if c.log == nil {
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

	_, _ = c.log.Write([]byte(fmt.Sprintf("[CieloConecta] Request: %s \n[CieloConecta] Response: %s \n", requestDump, responseDump)))
}

func (c *Client) SetLogger(w io.Writer) {
	c.log = w
}

func (c *Client) writeLog(message string) {
	if c.log == nil {
		return
	}

	if strings.HasSuffix(message, "\n") {
		_, _ = c.log.Write([]byte(fmt.Sprintf("[CieloConecta] %s", message)))
		return
	}

	if c.log != nil {
		_, _ = c.log.Write([]byte(fmt.Sprintf("[CieloConecta] %s\n", message)))
	}
}
