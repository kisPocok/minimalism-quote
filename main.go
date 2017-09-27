package main

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/sys/unix"
)

func main() {
	c, err := grabQuote()
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

func grabQuote() (q quote, err error) {
	response, err := http.Get("https://minimalmaxims.com/")
	if err != nil {
		return
	}
	defer response.Body.Close()

	token := html.NewTokenizer(response.Body)
	for {
		actual := token.Next()
		switch {
		case actual == html.StartTagToken:
			t := token.Token()
			if t.Data == "span" {
				if class(t.Attr) == "quotable-quote" {
					token.Next() // skip paragraph
					token.Next() // html.TextToken, the quote node
					t := token.Token()
					q.content = t.String()
				}
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
	return html.Attribute{}, errors.New("Missing key!")
}

func class(a []html.Attribute) string {
	attr, err := findAttr(a, "class")
	if err != nil {
		return ""
	}
	return attr.Val
}

/*
var class = func() func(a []html.Attribute) string {
	return func(a []html.Attribute) string {
		attr, err := findAttr(a, "class")
		if err != nil {
			return ""
		}
		return attr.Val
	}
}()
*/
