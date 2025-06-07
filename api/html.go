package api

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/internal/discord"
	html_page "github.com/typical-developers/discord-bot-backend/internal/html/pages"
	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

var (
	version = time.Now().Unix()
	// css     = map[string]string{}
	// js      = map[string]string{}
	// assets  = map[string]string{}
	assets = map[string]string{}
)

func init() {
	initPaths()
}

func initPaths() {
	extWhitelist := []string{".js", ".css", ".png", ".jpg", ".jpeg", ".svg", ".ico", ".ttf", ".woff", ".woff2"}

	err := filepath.Walk("./internal/html", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileExt := filepath.Ext(path)
		if info.IsDir() || !slices.Contains(extWhitelist, fileExt) {
			return nil
		}

		fileName := filepath.Base(path)
		currentDir := filepath.Dir(path)
		relativeDir := strings.TrimPrefix(filepath.ToSlash(currentDir), "internal/html")

		fileLocation := fmt.Sprintf("%s/%s", relativeDir, fileName)
		fileLocation = strings.TrimPrefix(fileLocation, "/")

		assets[fileLocation] = path

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func Version(c *fiber.Ctx) error {
	return c.JSON(map[string]int64{
		"version": version,
	})
}

//	@Router	/{file} [get]
//	@Tags	HTML Generation
//
//	@Param	file	path	string	true	"Asset file location."
//
// nolint:staticcheck
func GetHTMLAsset(c *fiber.Ctx) error {
	location := c.Params("*")

	filePath := assets[location]
	if filePath == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	ext := filepath.Ext(filePath)

	mimeType := mime.TypeByExtension(ext)

	//
	cssMimeOverrides := []string{".ttf", ".woff", ".woff2"}
	if slices.Contains(cssMimeOverrides, ext) {
		mimeType = "text/css"
	} else if mimeType == "" {
		mimeType = http.DetectContentType(fileContents)
	}

	c.Set("content-type", mimeType)
	return c.Send(fileContents)
}

//	@Router		/guild/{guild_id}/member/{member_id}/profile/card [get]
//	@Summary	Get the HTML generation for a member's profile card.
//	@Tags		HTML Generation
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//	@Param		member_id	path		string	true	"The member ID."
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func MemberProfileCard(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	memberId := c.Params("member_id")

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "guild settings not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	profile, err := dbutil.GetMemberProfile(ctx, queries, guildId, memberId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "member not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get member profile.", "guild_id", guildId, "member_id", memberId, "error", err)

		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	roles := dbutil.MapMemberRoles(int(profile.ActivityPoints), settings.ChatActivityRoles)

	if roles.Next == nil {
		roles.Next = &models.ActivityRoleProgress{}
	}

	guild, err := discord.Client.Cache.Guild(guildId)
	if err != nil {
		logger.Log.WithSource.Error("Failed to get guild info.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	member, err := discord.Client.Cache.GuildMember(guildId, memberId)
	if err != nil {
		logger.Log.WithSource.Error("Failed to get member info.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	protocol := c.Protocol()
	if c.IsFromLocal() {
		protocol = "http"
	}

	avatarUrl := member.AvatarURL("100")
	var chatActivityRole *html_page.ActivityRole
	if roles.Current != nil {
		for _, role := range guild.Roles {
			if role.ID != roles.Current.RoleID {
				continue
			}

			chatActivityRole = &html_page.ActivityRole{
				Text:   role.Name,
				Accent: fmt.Sprintf("#%06X", role.Color),
			}
			break
		}
	}

	baseUrl := fmt.Sprintf("%s://%s", protocol, c.Hostname())
	html := html_page.ProfileCard(html_page.ProfileCardProps{
		APIUrl:      baseUrl,
		DisplayName: member.DisplayName(),
		Username:    member.User.Username,
		AvatarURL:   avatarUrl,
		ChatActivity: html_page.ActivityInfo{
			Rank:               int(profile.ChatRank),
			TotalPoints:        int(profile.ActivityPoints),
			RoleCurrentPoints:  roles.Next.Progress,
			RoleRequiredPoints: roles.Next.RequiredPoints,
			CurrentTitleInfo:   chatActivityRole,
		},
	})
	c.Set("content-type", fiber.MIMETextHTML)

	return html.Render(c.Context())
}

//	@Router		/guild/{guild_id}/activity-leaderboard/card [get]
//	@Summary	Get the HTML generation for a member's profile card.
//	@Tags		HTML Generation
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id		path		string					true	"The guild ID."
//	@Param		activity_type	query		models.ActivityType		true	"The activity type."
//	@Param		display			query		models.LeaderboardType	true	"The leaderboard display type."
//
//	@Failure	400				{object}	models.APIResponse[ErrorResponse]
//	@Failure	500				{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func ActivityLeaderboardCard(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	activityType := models.ActivityType(c.Query("activity_type"))
	displayType := models.LeaderboardType(c.Query("display"))

	if !activityType.Valid() {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "activity_type is not valid (chat).",
			},
		})
	}

	if !displayType.Valid() {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "display is not valid (all, monthly, weekly).",
			},
		})
	}

	protocol := c.Protocol()
	if c.IsFromLocal() {
		protocol = "http"
	}

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	guild, err := discord.Client.Cache.Guild(guildId)
	if err != nil {
		logger.Log.WithSource.Error("Failed to get guild info.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	leaderboardData := []html_page.LeaderboardDataField{}
	var leaderboardName string
	switch displayType {
	case models.LeaderboardTypeAllTime:
		leaderboardName = fmt.Sprintf("%s Activity - All Time", html_page.Uppercase(string(activityType)))
		leaderboard, err := queries.GetAllTimeChatActivityRankings(ctx, db.GetAllTimeChatActivityRankingsParams{
			GuildID:  guildId,
			OffsetBy: 0,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
					Success: false,
					Data: models.ErrorResponse{
						Message: "guild not found.",
					},
				})
			}

			logger.Log.WithSource.Error("Failed to get guild info.", "guild_id", guildId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "internal server error.",
				},
			})
		}

		for _, rank := range leaderboard {
			member, err := discord.Client.Cache.GuildMember(guildId, rank.MemberID)
			if err != nil {
				leaderboardData = append(leaderboardData, html_page.LeaderboardDataField{
					Rank:     int(rank.Rank),
					Username: rank.MemberID,
					Value:    int(rank.ActivityPoints),
				})
				continue
			}

			leaderboardData = append(leaderboardData, html_page.LeaderboardDataField{
				Rank:     int(rank.Rank),
				Username: fmt.Sprintf("@%s", member.User.Username),
				Value:    int(rank.ActivityPoints),
			})
		}
	case models.LeaderboardTypeMonthly:
		leaderboardName = fmt.Sprintf("%s Activity - This Month", html_page.Uppercase(string(activityType)))
		leaderboard, err := queries.GetMonthlyActivityLeaderboard(ctx, db.GetMonthlyActivityLeaderboardParams{
			GuildID:  guildId,
			OffsetBy: 0,
		})

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
					Success: false,
					Data: models.ErrorResponse{
						Message: "guild not found.",
					},
				})
			}

			logger.Log.WithSource.Error("Failed to get guild info.", "guild_id", guildId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "internal server error.",
				},
			})
		}

		for _, rank := range leaderboard {
			member, err := discord.Client.Cache.GuildMember(guildId, rank.MemberID)
			if err != nil {
				leaderboardData = append(leaderboardData, html_page.LeaderboardDataField{
					Rank:     int(rank.Rank),
					Username: rank.MemberID,
					Value:    int(rank.EarnedPoints),
				})
				continue
			}

			leaderboardData = append(leaderboardData, html_page.LeaderboardDataField{
				Rank:     int(rank.Rank),
				Username: fmt.Sprintf("@%s", member.User.Username),
				Value:    int(rank.EarnedPoints),
			})
		}
	case models.LeaderboardTypeWeekly:
		leaderboardName = fmt.Sprintf("%s Activity - This Week", html_page.Uppercase(string(activityType)))
		leaderboard, err := queries.GetWeeklyActivityLeaderboard(ctx, db.GetWeeklyActivityLeaderboardParams{
			GuildID:  guildId,
			OffsetBy: 0,
		})

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
					Success: false,
					Data: models.ErrorResponse{
						Message: "guild not found.",
					},
				})
			}

			logger.Log.WithSource.Error("Failed to get guild info.", "guild_id", guildId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "internal server error.",
				},
			})
		}

		for _, rank := range leaderboard {
			member, err := discord.Client.Cache.GuildMember(guildId, rank.MemberID)
			if err != nil {
				leaderboardData = append(leaderboardData, html_page.LeaderboardDataField{
					Rank:     int(rank.Rank),
					Username: rank.MemberID,
					Value:    int(rank.EarnedPoints),
				})
				continue
			}

			leaderboardData = append(leaderboardData, html_page.LeaderboardDataField{
				Rank:     int(rank.Rank),
				Username: fmt.Sprintf("@%s", member.User.Username),
				Value:    int(rank.EarnedPoints),
			})
		}
	}

	baseUrl := fmt.Sprintf("%s://%s", protocol, c.Hostname())
	html := html_page.SeverLeaderboard(html_page.SeverLeaderboardProps{
		APIUrl: baseUrl,
		ServerInfo: html_page.ServerInfo{
			Icon: guild.IconURL("100"),
			Name: guild.Name,
		},
		LeaderboardInfo: html_page.LeaderboardInfo{
			Name: leaderboardName,
			Data: leaderboardData,
		},
	})

	c.Set("content-type", fiber.MIMETextHTML)
	return html.Render(c.Context())
}
