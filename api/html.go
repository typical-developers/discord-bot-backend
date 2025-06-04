package api

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
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
//	@Param	file	path	string	true "Asset file location."
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
//	@Param		avatar_url	query		string	true	"The avatar of the member."
//
//	@Success	200			{object}	models.APIResponse[MemberProfile]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func MemberProfileCard(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	memberId := c.Params("member_id")

	displayName := c.Query("display_name")
	username := c.Query("username")
	avatarUrl := c.Query("avatar_url")

	if _, err := url.Parse(avatarUrl); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "invalid avatar url.",
			},
		})
	}

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

		logger.Log.Error("Failed to get guild settings", "error", err)
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

		logger.Log.Error("Failed to get member profile.", "guild_id", guildId, "member_id", memberId, "error", err)

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

	html := html_page.ProfileCard(html_page.ProfileCardProps{
		DisplayName: displayName,
		Username:    username,
		AvatarURL:   avatarUrl,
		ChatActivity: html_page.ActivityInfo{
			Rank:               int(profile.ChatRank),
			TotalPoints:        int(profile.ActivityPoints),
			RoleCurrentPoints:  roles.Next.Progress,
			RoleRequiredPoints: roles.Next.RequiredPoints,
		},
	})
	c.Set("content-type", fiber.MIMETextHTML)

	return html.Render(c.Context())
}
