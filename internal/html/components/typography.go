package html_components

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
