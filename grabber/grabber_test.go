package grabber

import (
	"io"
	"net/http"
	"testing"

	"net/http/httptest"

	"golang.org/x/net/html"
)

const expected = "Pörkölt, jó kutya."

func TestGrabberShouldFindMinimalismLikeText(t *testing.T) {
	srv := fakeHTTPServer()
	defer srv.Close()

	q, err := NewSource(srv.URL, "span", "quotable-quote", fn).GrabQuote()
	if err != nil {
		t.Fatal(err)
	}

	if q.content != expected {
		t.Errorf("Grabbed message does not match, expected %s, got: %s.", expected, q.content)
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

// TODO same like clients.searchMinimalmaxismQuote()
func fn(token html.Tokenizer) string {
	token.Next() // skip paragraph
	token.Next() // html.TextToken, the quote node
	return token.Token().String()
}
