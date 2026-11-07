package common

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resp(code int, body string) *http.Response {
	return &http.Response{StatusCode: code, Body: io.NopCloser(strings.NewReader(body))}
}

func TestReadBody_OK(t *testing.T) {
	b, err := ReadBody(resp(200, "hello"), 10)
	if err != nil || string(b) != "hello" {
		t.Fatalf("b=%q err=%v", b, err)
	}
	if _, err := ReadBody(resp(204, ""), 10); err != nil {
		t.Fatal(err)
	}
}

func TestReadBody_NonSuccessStatus(t *testing.T) {
	for _, code := range []int{199, 301, 400, 404, 500, 502} {
		_, err := ReadBody(resp(code, `{"errcode":0}`), 100)
		var se *HTTPStatusError
		if !errors.As(err, &se) || se.StatusCode != code {
			t.Fatalf("code %d: err=%v", code, err)
		}
		if !strings.Contains(err.Error(), "HTTP") {
			t.Fatalf("message = %q", err)
		}
	}
}

func TestHTTPStatusError_Error(t *testing.T) {
	if got := (&HTTPStatusError{StatusCode: 500}).Error(); got != "HTTP 500" {
		t.Fatal(got)
	}
	if got := (&HTTPStatusError{StatusCode: 500, Body: "x"}).Error(); got != "HTTP 500: x" {
		t.Fatal(got)
	}
}

func TestReadBody_SnippetTruncated(t *testing.T) {
	_, err := ReadBody(resp(500, strings.Repeat("a", 5000)), 100)
	var se *HTTPStatusError
	if !errors.As(err, &se) || len(se.Body) != statusSnippetBytes {
		t.Fatalf("err=%v", err)
	}
}

func TestReadBody_TooLarge(t *testing.T) {
	if _, err := ReadBody(resp(200, "abcdef"), 5); err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("err = %v", err)
	}
	if _, err := ReadBody(resp(200, "abcde"), 5); err != nil {
		t.Fatalf("exact limit: %v", err)
	}
}

type badBody struct{}

func (badBody) Read([]byte) (int, error) { return 0, errors.New("read boom") }
func (badBody) Close() error             { return nil }

func TestReadBody_ReadErrorAndNil(t *testing.T) {
	if _, err := ReadBody(&http.Response{StatusCode: 200, Body: badBody{}}, 10); err == nil {
		t.Fatal("expected read error")
	}
	if _, err := ReadBody(nil, 10); err == nil {
		t.Fatal("expected nil response error")
	}
}

func TestFetchAccessToken_Non2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"access_token":"leak","expires_in":7200}`))
	}))
	defer server.Close()
	_, err := FetchAccessToken(context.Background(), AppCredential{Host: server.URL, AppID: "wx", Secret: "s"}, server.Client())
	var se *HTTPStatusError
	if !errors.As(err, &se) || se.StatusCode != 500 {
		t.Fatalf("err = %v", err)
	}
}

func TestFetchAccessToken_TooLarge(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"access_token":"`))
		_, _ = w.Write([]byte(strings.Repeat("a", int(MaxResponseBytes)+10)))
		_, _ = w.Write([]byte(`"}`))
	}))
	defer server.Close()
	_, err := FetchAccessToken(context.Background(), AppCredential{Host: server.URL, AppID: "wx", Secret: "s"}, server.Client())
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("err = %v", err)
	}
}
