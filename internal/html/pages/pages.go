package html_page

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	format = message.NewPrinter(message.MatchLanguage("en"))
)

func Uppercase(s string) string {
	c := cases.Title(language.English)
	return c.String(s)
}
