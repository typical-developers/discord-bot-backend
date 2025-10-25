package pages

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	Format = message.NewPrinter(message.MatchLanguage("en"))
)

func Uppercase(s string) string {
	c := cases.Title(language.English)
	return c.String(s)
}
