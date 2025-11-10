package layouts

import (
	"fmt"

	pages "github.com/typical-developers/discord-bot-backend/internal/pages"
	. "github.com/typical-developers/discord-bot-backend/internal/pages/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type AvatarProps struct {
	URL string
}

func Avatar(props AvatarProps) Node {
	return Img(
		Class("avatar"),
		Src(props.URL),
	)
}

type TagProps struct {
	Accent string
	Icon   Node
	Text   string
}

func Tag(props TagProps) Node {
	return Div(
		Class("tag"),
		Style(fmt.Sprintf("--accent: %s;", props.Accent)),
		props.Icon,
		Typography(TypographyProps{
			Size:   FontSizeNormal,
			Weight: FontWeightBold,
		}, Text(props.Text)),
	)
}

type ActivityRole struct {
	Accent string
	Text   string
}

type UserInfoProps struct {
	DisplayName string
	Username    string

	TopChatActivityRole *ActivityRole
}

func UserInfo(props UserInfoProps) Node {
	var topActivityRole Node
	if props.TopChatActivityRole != nil && props.TopChatActivityRole.Text != "" && props.TopChatActivityRole.Accent != "" {
		topActivityRole = Tag(TagProps{
			Accent: props.TopChatActivityRole.Accent,
			Icon:   ChatBubbleIcon(IconProps{Width: "18px", Height: "18px"}),
			Text:   props.TopChatActivityRole.Text,
		})
	}

	return Div(
		Class("user-info"),
		Div(
			Class("names"),
			Typography(TypographyProps{
				Size:   FontSizeLarge,
				Weight: FontWeightBlack,
			}, Text(props.DisplayName)),
			Typography(TypographyProps{
				Size:   FontSizeNormal,
				Weight: FontWeightSemiBold,
			}, Text(fmt.Sprintf("@%s", props.Username))),
		),
		Div(
			Class("tags"),

			If(props.TopChatActivityRole != nil, topActivityRole),
		))
}

type ProgressBarProps struct {
	CurrentProgress  int
	RequiredProgress int
}

func ProgressBar(props ProgressBarProps) Node {
	progress := int(float64(props.CurrentProgress) / float64(props.RequiredProgress) * 100)

	pos1 := 100 - progress
	pos2 := 200 - progress

	progressText := pages.Format.Sprintf("%d / %d", props.CurrentProgress, props.RequiredProgress)
	if props.CurrentProgress >= props.RequiredProgress {
		progressText = "MAX"

		pos1 = 0
		pos2 = 100
	}

	return Div(
		Class("progress-bar"),
		Div(
			Class("bar"),
			Style(fmt.Sprintf("--gradient-1-pos: %d%%; --gradient-2-pos: %d%%; transform: translateX(-%d%%);", pos1, pos2, 100-progress)),
		),
		Div(
			Class("progress"),
			Typography(TypographyProps{
				Size:   FontSizeNormal,
				Weight: FontWeightBold,
			}, Text(progressText)),
		),
	)
}

type RankingInfo struct {
	AllTime int
	Weekly  int
	Monthly int
}

type ProgressGroupHeaderProps struct {
	ActivityType string
	Icon         Node
	Ranking      RankingInfo
	TotalPoints  int
}

func ProgressGroupHeader(props ProgressGroupHeaderProps) Node {
	return Div(
		Class("header"),
		props.Icon,
		Div(
			Class("details with-icon"),
			Typography(TypographyProps{
				Size:   FontSizeNormal,
				Weight: FontWeightBlack,
			}, Text(pages.Format.Sprintf("%s Activity", props.ActivityType))),
			Typography(TypographyProps{
				Size:   FontSizeNormal,
				Weight: FontWeightRegular,
			}, Text(pages.Format.Sprintf("%d Points", props.TotalPoints))),
		),
		Div(
			Class("details ranking"),
			Typography(TypographyProps{
				Size:   FontSizeMedium,
				Weight: FontWeightBlack,
			}, RankingText(props.Ranking.AllTime)),

			Div(
				Class("leaderboards"),

				If(props.Ranking.Weekly != 0, Typography(TypographyProps{
					Size:   FontSizeSmall,
					Weight: FontWeightBold,
				}, RankingText(props.Ranking.Weekly, "Weekly"))),

				If(props.Ranking.Weekly != 0 && props.Ranking.Monthly != 0, Div(Class("divider"))),

				If(props.Ranking.Monthly != 0, Group{
					Typography(TypographyProps{
						Size:   FontSizeSmall,
						Weight: FontWeightBold,
					}, RankingText(props.Ranking.Monthly, "Monthly")),
				}),
			),
		),
	)
}

type ProgressGroupProps struct {
	ActivityType   string
	Icon           Node
	Ranking        RankingInfo
	TotalPoints    int
	CurrentPoints  int
	RequiredPoints int
}

func ProgressGroup(props ProgressGroupProps) Node {
	return Div(
		Class("progress-group"),
		ProgressGroupHeader(ProgressGroupHeaderProps{
			Icon:         props.Icon,
			Ranking:      props.Ranking,
			TotalPoints:  props.TotalPoints,
			ActivityType: props.ActivityType,
		}),
		ProgressBar(ProgressBarProps{
			CurrentProgress:  props.CurrentPoints,
			RequiredProgress: props.RequiredPoints,
		}),
	)
}

type ActivityInfo struct {
	Ranking            RankingInfo
	TotalPoints        int
	RoleCurrentPoints  int
	RoleRequiredPoints int
	CurrentTitleInfo   *ActivityRole
}

type CardStyling struct {
	Gradient1HSL       string
	Gradient2HSL       string
	BackgroundColor    string
	BackgroundImageURL string
}

type ProfileCardProps struct {
	CardStyle          int32
	CardStyleOverrides CardStyling

	DisplayName  string
	Username     string
	AvatarURL    string
	ChatActivity ActivityInfo
}

func ProfileCard(props ProfileCardProps) Node {
	cardStyling := &CardStyling{
		Gradient1HSL:       "21, 97%, 69%",
		Gradient2HSL:       "270, 94%, 64%",
		BackgroundColor:    "#0E0911",
		BackgroundImageURL: "url(/static/images/card-style_0-background.png) no-repeat",
	}

	switch props.CardStyle {
	// This will set based on overrides.
	case 1:
		overrides := props.CardStyleOverrides

		if overrides.Gradient1HSL != "" {
			cardStyling.Gradient1HSL = overrides.Gradient1HSL
		}
		if overrides.Gradient2HSL != "" {
			cardStyling.Gradient2HSL = overrides.Gradient2HSL
		}
		if overrides.BackgroundImageURL != "" {
			cardStyling.BackgroundImageURL = fmt.Sprintf("url(%s) no-repeat center/cover", overrides.BackgroundImageURL)
		}
	case 2:
		cardStyling.Gradient1HSL = "263, 97%, 70%"
		cardStyling.Gradient2HSL = "234, 95%, 64%"
		cardStyling.BackgroundColor = "linear-gradient(180deg, #9F66FD 0.60%, #4D5EFA 25%);"
		cardStyling.BackgroundImageURL = "url(/static/images/card-style_2-background.png) no-repeat"
	}

	return HTML5(HTML5Props{
		Head: []Node{
			Link(Rel("stylesheet"), Href("/static/css/index.css")),
			Link(Rel("stylesheet"), Href("/static/css/profile-card.css")),
			Link(Rel("stylesheet"), Href("/static/css/typography.css")),
			Link(Rel("stylesheet"), Href("/static/css/icons.css")),

			Link(Rel("stylesheet"), Href("/static/css/fixel.css"), As("font")),
		},
		Body: []Node{
			Style(fmt.Sprintf(`
				--profile-card-background: %s;
				--profile-card-image: %s;
			`, cardStyling.BackgroundColor, cardStyling.BackgroundImageURL)),
			Div(ID("root"),
				Div(
					Class("content"),
					Div(
						Class("user"),
						Avatar(AvatarProps{
							URL: props.AvatarURL,
						}),
						UserInfo(UserInfoProps{
							DisplayName:         props.DisplayName,
							Username:            props.Username,
							TopChatActivityRole: props.ChatActivity.CurrentTitleInfo,
						}),
					),
					Div(
						Class("progress-info"),

						Style(fmt.Sprintf(`
							--gradient-1-hsl: %s;
							--gradient-2-hsl: %s;
						`, cardStyling.Gradient1HSL, cardStyling.Gradient2HSL)),
						ProgressGroup(ProgressGroupProps{
							ActivityType:   "Chat",
							Icon:           ChatBubbleIcon(IconProps{Width: "26px", Height: "26px"}),
							Ranking:        props.ChatActivity.Ranking,
							TotalPoints:    props.ChatActivity.TotalPoints,
							CurrentPoints:  props.ChatActivity.RoleCurrentPoints,
							RequiredPoints: props.ChatActivity.RoleRequiredPoints,
						}),
					),
				),
				Div(
					Classes{"background": true},
				),
			),
		},
	})
}
