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

func TestGrabberShouldFailOn404(t *testing.T) {
	srv := fakeHTTPServer()
	defer srv.Close()
	_, err := NewSource(srv.URL+"/404", "span", "quotable-quote", SkipTheFirstTag).GrabQuote()
	if err == nil {
		t.Error("Non existent url should fail")
	}
}

func TestGrabberShouldFailOnNoServer(t *testing.T) {
	_, err := NewSource("http://localhost:8080", "span", "quotable-quote", SkipTheFirstTag).GrabQuote()
	if err == nil {
		t.Error("It should fail, really")
	}
}

func TestGrabberShouldFailIfNoHTMLFound(t *testing.T) {
	srv := fakeHTTPServer()
	defer srv.Close()
	_, err := NewSource(srv.URL+"/empty", "span", "quotable-quote", SkipTheFirstTag).GrabQuote()
	if err == nil || err.Error() != "cannot find desired parts" {
		t.Error("Non existent part, should fail")
	}
}

func fakeHTTPServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dummyHTMLResponder)
	mux.HandleFunc("/empty", emptyHTMLResponder)
	mux.Handle("/404", http.NotFoundHandler())
	return httptest.NewServer(mux)
}

func dummyHTMLResponder(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<html><body><span class="quotable-quote"><p>`+expected+`</p></span></body></html>`)
}

func emptyHTMLResponder(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<html><body></body></html>`)
}
