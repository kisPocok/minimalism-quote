package clients

import (
	"minimal/grabber"

	"golang.org/x/net/html"
)

func Minimalmaxism() grabber.Source {
	return grabber.NewSource(
		"https://minimalmaxims.com/",
		"span",
		"quotable-quote",
		searchMinimalmaxismQuote,
	)
}

func searchMinimalmaxismQuote(token html.Tokenizer) string {
	token.Next() // skip paragraph
	token.Next() // html.TextToken, the quote node
	return token.Token().String()
}
