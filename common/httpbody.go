package common

import (
	"fmt"
	"io"
	"net/http"
)

const (
	// MaxResponseBytes caps JSON API response bodies (4 MiB).
	MaxResponseBytes int64 = 4 << 20
	// MaxImageResponseBytes caps binary image responses such as QR codes (10 MiB).
	MaxImageResponseBytes int64 = 10 << 20

	statusSnippetBytes = 512
)

// HTTPStatusError reports a non-2xx HTTP response.
type HTTPStatusError struct {
	// StatusCode is the HTTP status code.
	StatusCode int
	// Body is a truncated copy of the response body.
	Body string
}

func (e *HTTPStatusError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

// ReadBody checks the HTTP status and then reads the body, limited to limit bytes.
// Non-2xx responses return *HTTPStatusError without being parsed. Bodies larger
// than limit return an error instead of being silently truncated.
func ReadBody(resp *http.Response, limit int64) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("common: nil response")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, statusSnippetBytes))
		return nil, &HTTPStatusError{StatusCode: resp.StatusCode, Body: string(snippet)}
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("common: response body exceeds %d bytes", limit)
	}
	return b, nil
}
