package components

import (
	"fmt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type FontSize string

func (fs FontSize) String() string {
	return string(fs)
}

const (
	FontSizeSmall  FontSize = "small"
	FontSizeMedium FontSize = "medium"
	FontSizeNormal FontSize = "normal"
	FontSizeLarge  FontSize = "large"
	FontSizeXLarge FontSize = "extra-large"
)

type FontWeight string

func (fw FontWeight) String() string {
	return string(fw)
}

const (
	FontWeightRegular  FontWeight = "regular"
	FontWeightSemiBold FontWeight = "semi-bold"
	FontWeightBold     FontWeight = "bold"
	FontWeightBlack    FontWeight = "black"
)

type TypographyProps struct {
	Size   FontSize
	Weight FontWeight
}

func Typography(props TypographyProps, children ...Node) Node {
	fontSize := fmt.Sprintf("font-size-%s", props.Size.String())
	fontWeight := fmt.Sprintf("font-weight-%s", props.Weight.String())

	return Div(
		Classes{
			"typography": true,
			fontSize:     true,
			fontWeight:   true,
		},
		Div(children...),
	)
}

func RankingText(rank int, text ...string) Node {
	var colorOverride string
	switch rank {
	case 1:
		colorOverride = "linear-gradient(90deg, #82F5FF 0%, #14AAB8 100%)"
	case 2:
		colorOverride = "linear-gradient(90deg, #FFD54C 0%, #F8A304 100%)"
	case 3:
		colorOverride = "linear-gradient(90deg, #BFD7D9 0%, #859EAD 100%)"
	}

	rankIfno := fmt.Sprintf("#%d", rank)
	if len(text) > 0 {
		for i := range text {
			content := text[i]
			if content == "" {
				continue
			}

			rankIfno += fmt.Sprintf(" %s", text[i])
		}
	}

	return Div(Group{
		If(colorOverride != "", Group{
			Class("gradient-text"),
			Style(fmt.Sprintf(`background-image: %s`, colorOverride)),
		}),

		Text(rankIfno),
	})
}
