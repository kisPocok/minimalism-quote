package main

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/sys/unix"
)

func main() {
	source := NewSource(
		"https://minimalmaxims.com/",
		"span",
		"quotable-quote",
		func(token html.Tokenizer) string {
			token.Next() // skip paragraph
			token.Next() // html.TextToken, the quote node
			return token.Token().String()
		})

	c, err := grabQuote(source)
	if err != nil {
		unix.Exit(1)
	}
	fmt.Println(c)
	unix.Exit(0)
}

type quote struct {
	content string
}

func (q quote) String() string {
	return q.content
}

type source struct {
	url           string
	htmlTag       string
	htmlClassName string
	handle        func(html.Tokenizer) string
}

func NewSource(url, tag, class string, fn func(html.Tokenizer) string) *source {
	return &source{
		url:           url,
		htmlTag:       tag,
		htmlClassName: class,
		handle:        fn,
	}
}

func grabQuote(s *source) (q quote, err error) {
	response, err := http.Get(s.url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var getClass = getAttr("class")
	token := html.NewTokenizer(response.Body)
	for {
		actual := token.Next()
		switch {
		case actual == html.StartTagToken:
			t := token.Token()
			if t.Data == s.htmlTag && getClass(t.Attr) == s.htmlClassName {
				q.content = s.handle(*token)
				return
			}
		case actual == html.ErrorToken:
			// We are done
			return
		}
	}
	return q, nil
}

func findAttr(a []html.Attribute, key string) (html.Attribute, error) {
	for _, v := range a {
		if v.Key == key {
			return v, nil
		}
	}
	return html.Attribute{}, errors.New("missing attribute")
}

func getAttr(attrName string) func(a []html.Attribute) string {
	return func(a []html.Attribute) string {
		attr, err := findAttr(a, attrName)
		if err != nil {
			return ""
		}
		return attr.Val
	}
}
