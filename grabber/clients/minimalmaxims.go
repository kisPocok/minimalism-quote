package clients

import (
	"minimal/grabber"
)

type Source interface {
	GrabQuote() (string, error)
}

func Minimalmaxism() Source {
	return grabber.NewSource(
		"https://minimalmaxims.com/",
		"span",
		"quotable-quote",
		grabber.SkipTheFirstTag,
	)
}
