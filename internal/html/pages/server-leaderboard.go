package html_page

import (
	"fmt"
	"time"

	. "github.com/typical-developers/discord-bot-backend/internal/html/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ServerIconProps struct {
	URL string
}

func ServerIcon(props ServerIconProps) Node {
	return Img(
		Class("server-icon"),
		Src(props.URL),
	)
}

type LeaderboardProps struct {
	ServerName      string
	LeaderboardName string
	ResetTime       *time.Time
}

func Leaderboard(props LeaderboardProps) Node {
	return Div(
		Class("leaderboard-info"),
		Typography(TypographyProps{
			Size:   FontSizeLarge,
			Weight: FontWeightBlack,
		}, Text(props.ServerName)),
		Typography(TypographyProps{
			Size:   FontSizeMedium,
			Weight: FontWeightSemiBold,
		}, Text(props.LeaderboardName)),
	)
}

type LeaderboardDataField struct {
	Rank     int
	Username string
	Value    int
}

type LeaderboardDataProps struct {
	Key  string
	Data []LeaderboardDataField
}

func LeaderboardData(props LeaderboardDataProps) Node {
	var rowItems []Node

	for _, data := range props.Data {
		value := format.Sprintf("%d", data.Value)

		d := Li(
			Class("row-item"),
			Typography(
				TypographyProps{
					Size:   FontSizeMedium,
					Weight: FontWeightBlack,
				},
				Class("rank"), Text(fmt.Sprintf("#%d", data.Rank)),
			),
			Typography(
				TypographyProps{
					Size:   FontSizeMedium,
					Weight: FontWeightSemiBold,
				},
				Class("member"), Text(data.Username),
			),
			Typography(
				TypographyProps{
					Size:   FontSizeMedium,
					Weight: FontWeightSemiBold,
				},
				Class("key"), Text(value),
			),
		)

		rowItems = append(rowItems, d)
	}

	return Div(
		Class("leaderboard"),
		Ol(
			Class("rows"),
			Li(
				Class("row-header"),
				Typography(
					TypographyProps{
						Size:   FontSizeMedium,
						Weight: FontWeightSemiBold,
					},
					Class("rank"), Text("Rank"),
				),
				Typography(
					TypographyProps{
						Size:   FontSizeMedium,
						Weight: FontWeightSemiBold,
					},
					Class("member"), Text("Member"),
				),
				Typography(
					TypographyProps{
						Size:   FontSizeMedium,
						Weight: FontWeightSemiBold,
					},
					Class("key"), Text(props.Key),
				),
			),
			Group(rowItems),
		),
	)
}

type ServerInfo struct {
	Icon string
	Name string
}

type LeaderboardInfo struct {
	Name string
	Data []LeaderboardDataField
}

type SeverLeaderboardProps struct {
	APIUrl          string
	ServerInfo      ServerInfo
	LeaderboardInfo LeaderboardInfo
}

func SeverLeaderboard(props SeverLeaderboardProps) Node {
	return HTML5(HTML5Props{
		Head: []Node{
			// If(os.Getenv("ENVIRONMENT") == "development",
			// 	Script(Src("/html/hot-reload.js")),
			// ),
			Base(Href(props.APIUrl)),

			Link(Rel("stylesheet"), Href("/html/index.css")),
			Link(Rel("stylesheet"), Href("/html/pages/server-leaderboard.css")),
			Link(Rel("stylesheet"), Href("/html/components/typography.css")),
			Link(Rel("stylesheet"), Href("/html/components/icons.css")),

			Link(Rel("stylesheet"), Href("/html/fonts/Fixel/fixel.css"), As("font")),
		},
		Body: []Node{
			Div(ID("root"),
				Div(
					Class("content"),
					Div(
						Class("server"),
						ServerIcon(ServerIconProps{
							URL: props.ServerInfo.Icon,
						}),
						Leaderboard(LeaderboardProps{
							ServerName:      props.ServerInfo.Name,
							LeaderboardName: props.LeaderboardInfo.Name,
						}),
					),
					LeaderboardData(LeaderboardDataProps{
						Key:  "Activity Points",
						Data: props.LeaderboardInfo.Data,
					}),
				),
				Div(
					Classes{"background": true},
				),
			),
		},
	})
}
