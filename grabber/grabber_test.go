package grabber

import (
	"io"
	"net/http"
	"testing"

	"net/http/httptest"
)

const expected = "Pörkölt, jó kutya."

func TestGrabberShouldFindMinimalismLikeText(t *testing.T) {
	srv := fakeHTTPServer()
	defer srv.Close()

	quote, err := NewSource(srv.URL, "span", "quotable-quote", SkipTheFirstTag).GrabQuote()
	if err != nil {
		t.Fatal(err)
	}

	if quote != expected {
		t.Errorf("Grabbed message does not match, expected %s, got: %s.", expected, quote)
	}
}

func fakeHTTPServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dummyHTMLResponder)
	return httptest.NewServer(mux)
}

func dummyHTMLResponder(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<html><body><span class="quotable-quote"><p>`+expected+`</p></span></body></html>`)
}
