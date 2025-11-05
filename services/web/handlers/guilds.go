package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	u "github.com/typical-developers/discord-bot-backend/internal/usecase"
	"github.com/typical-developers/discord-bot-backend/pkg/httpx"
)

type GuildHandler struct {
	uc u.GuildsUsecase
}

func NewGuildHandler(r *chi.Mux, uc u.GuildsUsecase) {
	h := GuildHandler{uc: uc}

	r.Route("/v1/guild/{guildId}", func(r chi.Router) {
		r.Get("/settings", h.GetGuildSettings)
		r.Post("/settings", h.CreateGuildSettings)
		r.Patch("/settings/activity", h.UpdateGuildActivitySettings)
		r.Post("/settings/activity-roles", h.CreateActivityRole)

		r.Patch("/settings/message-embeds", h.UpdateGuildMessageEmbedSettings)

		r.Get("/activity-leaderboard-card", h.GenerateGuildActivityLeaderboardCard)

		r.Route("/voice-room-lobby/{originChannelId}", func(r chi.Router) {
			r.Post("/", h.CreateVoiceRoomLobby)
			r.Get("/", h.GetVoiceRoomLobby)
			r.Patch("/", h.UpdateVoiceRoomLobby)
			r.Delete("/", h.DeleteVoiceRoomLobby)
			r.Post("/register", h.RegisterVoiceRoom)
		})

		r.Route("/voice-room/{channelId}", func(r chi.Router) {
			r.Get("/", h.GetVoiceRoom)
			r.Patch("/", h.UpdateVoiceRoom)
			r.Delete("/", h.DeleteVoiceRoom)
		})
	})

	r.Route("/v2/guild/{guildId}", func(r chi.Router) {
		r.Get("/activity-leaderboard", h.GetGuildActivityLeaderboard)
	})
}

//	@Router		/v1/guild/{guild_id}/settings [POST]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//
//	@Success	200			{object}	GuildSettingsResponse
//	@Failure	400			{object}	APIError
//
// nolint:staticcheck
func (h *GuildHandler) CreateGuildSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	settings, err := h.uc.RegisterGuild(ctx, guildId)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrGuildSettingsExists.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusConflict)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, GuildSettingsResponse{
		Success: true,
		Data:    *settings,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/settings [GET]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//
//	@Success	200			{object}	GuildSettingsResponse
//	@Failure	400			{object}	APIError
//	@Failure	404			{object}	APIError
//
// nolint:staticcheck
func (h *GuildHandler) GetGuildSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	settings, err := h.uc.GetGuildSettings(ctx, guildId)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrGuildNotFound.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, GuildSettingsResponse{
		Success: true,
		Data:    *settings,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/settings/activity  [PATCH]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string							true	"The guild ID."
//	@Param		settings	body		GuildActivitySettingsUpdateBody	true	"The activity settings."
//
//	@Success	200			{object}	GuildSettingsResponse
//	@Failure	400			{object}	APIError
//	@Failure	404			{object}	APIError
//
// nolint:staticcheck
func (h *GuildHandler) UpdateGuildActivitySettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	var updateBody *GuildActivitySettingsUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidRequestBody.Error(),
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		}

		return
	}
	if err := updateBody.Validate(); err != nil {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: err.Error(),
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		}

		return
	}

	settings, err := h.uc.UpdateGuildActivitySettings(ctx, guildId, u.UpdateAcitivtySettings{
		ChatActivity: updateBody.ChatActivity,
	})

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrGuildNotFound.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, GuildSettingsResponse{
		Success: true,
		Data:    *settings,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/settings/activity-roles [POST]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string						true	"The guild ID."
//	@Param		role		body		GuildActivityRoleCreateBody	true	"The activity settings."
//
//	@Success	200			{object}	GuildSettingsResponse
//	@Failure	400			{object}	APIError
//	@Failure	404			{object}	APIError
//
// nolint:staticcheck
func (h *GuildHandler) CreateActivityRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	var createBody *GuildActivityRoleCreateBody
	if err := json.NewDecoder(r.Body).Decode(&createBody); err != nil {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidRequestBody.Error(),
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		}

		return
	}

	_, err := h.uc.CreateActivityRole(ctx, guildId, createBody.ActivityType, createBody.RoleID, createBody.RequiredPoints)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrActivityRoleExists.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusConflict)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[any]{
		Success: true,
		Data:    nil,
	}, http.StatusCreated)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/settings/message-embeds [POST]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path	string	true	"The guild ID."
//
// nolint:staticcheck
func (h *GuildHandler) UpdateGuildMessageEmbedSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	var body *GuildMessageEmbedSettingsUpdateBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidRequestBody.Error(),
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		}

		return
	}

	if err := body.Validate(); err != nil {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: err.Error(),
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		}

		return
	}

	settings, err := h.uc.UpdateMessageEmbedSettings(ctx, guildId, u.UpdateMessageEmbedSettingsOpts{
		IsEnabled:             body.IsEnabled,
		AddDisabledChannel:    body.AddDisabledChannel,
		RemoveDisabledChannel: body.RemoveDisabledChannel,
		AddIgnoredChannel:     body.AddIgnoredChannel,
		RemoveIgnoredChannel:  body.RemoveIgnoredChannel,
		AddIgnoredRole:        body.AddIgnoredRole,
		RemoveIgnoredRole:     body.RemoveIgnoredRole,
	})

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrGuildNotFound.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, GuildSettingsResponse{
		Success: true,
		Data:    *settings,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router	/v1/guild/{guild_id}/activity-leaderboard-card [GET]
//	@Tags	Guilds
//
//	@Deprecated
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id		path	string	true	"The guild ID."
//	@Param		activity_type	query	string	true	"The activity type."	Enum(chat, voice)
//	@Param		time_period		query	string	true	"The time period."		Enum(weekly, monthly, all)
//
// nolint:staticcheck
func (h *GuildHandler) GenerateGuildActivityLeaderboardCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	activityType := httpx.GetQueryParam(r, "activity_type", "chat")
	timePeriod := httpx.GetQueryParam(r, "time_period", "all")

	card, err := h.uc.GenerateGuildActivityLeaderboardCard(ctx, guildId, activityType, timePeriod, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := card.Render(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//	@Router		/v2/guild/{guild_id}/activity-leaderboard-card [GET]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id		path	string	true	"The guild ID."
//	@Param		activity_type	query	string	true	"The activity type."	Enum(chat, voice)
//	@Param		time_period		query	string	true	"The time period."		Enum(weekly, monthly, all)
//
// nolint:staticcheck
func (h *GuildHandler) GetGuildActivityLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	referer := r.Header.Get("Referer")
	if referer == "" {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: "Referer header is required.",
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		}

		return
	}

	guildId := chi.URLParam(r, "guildId")
	activityType := httpx.GetQueryParam(r, "activity_type", "chat")
	timePeriod := httpx.GetQueryParam(r, "time_period", "all")
	pageStr := httpx.GetQueryParam(r, "page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	leaderboard, err := h.uc.GetGuildActivityLeaderboard(ctx, referer, guildId, activityType, timePeriod, page)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrGuildNotFound.Code:
				fallthrough
			case u.ErrLeaderboardNoRows.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[u.GuildLeaderboard]{
		Success: true,
		Data:    *leaderboard,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/voice-room-lobby/{origin_channel_id} [POST]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id			path	string	true	"The guild ID."
//	@Param		origin_channel_id	path	string	true	"The channel ID for the lobby origin."
//
// nolint:staticcheck
func (h *GuildHandler) CreateVoiceRoomLobby(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	originChannelId := chi.URLParam(r, "originChannelId")
	var body *VoiceRoomLobbySettings
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidRequestBody.Error(),
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		}

		return
	}

	lobby, err := h.uc.CreateVoiceRoomLobby(ctx, guildId, originChannelId, u.VoiceRoomLobbySettings{
		UserLimit:      body.UserLimit,
		CanRename:      body.CanRename,
		CanLock:        body.CanLock,
		CanAdjustLimit: body.CanAdjustLimit,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrVoiceRoomLobbyExists.Code:
				fallthrough
			case u.ErrVoiceRoomLobbyIsVoiceRoom.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusConflict)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[u.VoiceRoomLobby]{
		Success: true,
		Data:    *lobby,
	}, http.StatusCreated)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/voice-room-lobby/{origin_channel_id} [GET]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id			path	string	true	"The guild ID."
//	@Param		origin_channel_id	path	string	true	"The channel ID for the lobby origin."
//
// nolint:staticcheck
func (h *GuildHandler) GetVoiceRoomLobby(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	originChannelId := chi.URLParam(r, "originChannelId")

	lobby, err := h.uc.GetVoiceRoomLobby(ctx, guildId, originChannelId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrVoiceRoomLobbyNotFound.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}
		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[u.VoiceRoomLobby]{
		Success: true,
		Data:    *lobby,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/voice-room-lobby/{origin_channel_id} [PATCH]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id			path	string	true	"The guild ID."
//	@Param		origin_channel_id	path	string	true	"The channel ID for the lobby origin."
//
// nolint:staticcheck
func (h *GuildHandler) UpdateVoiceRoomLobby(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	originChannelId := chi.URLParam(r, "originChannelId")
	var body *VoiceRoomLobbySettings
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidRequestBody.Error(),
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		}

		return
	}

	lobby, err := h.uc.UpdateVoiceRoomLobby(ctx, guildId, originChannelId, u.VoiceRoomLobbySettings{
		UserLimit:      body.UserLimit,
		CanRename:      body.CanRename,
		CanLock:        body.CanLock,
		CanAdjustLimit: body.CanAdjustLimit,
	})

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrVoiceRoomLobbyNotFound.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[u.VoiceRoomLobby]{
		Success: true,
		Data:    *lobby,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/voice-room-lobby/{origin_channel_id} [DELETE]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id			path	string	true	"The guild ID."
//	@Param		origin_channel_id	path	string	true	"The channel ID for the lobby origin."
//
// nolint:staticcheck
func (h *GuildHandler) DeleteVoiceRoomLobby(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	originChannelId := chi.URLParam(r, "originChannelId")

	err := h.uc.DeleteVoiceRoomLobby(ctx, guildId, originChannelId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrVoiceRoomLobbyNotFound.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[any]{
		Success: true,
		Data:    nil,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/voice-room-lobby/{origin_channel_id}/register [POST]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id			path	string	true	"The guild ID."
//	@Param		origin_channel_id	path	string	true	"The channel ID for the lobby origin."
//
// nolint:staticcheck
func (h *GuildHandler) RegisterVoiceRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	originChannelId := chi.URLParam(r, "originChannelId")
	var body *VoiceRoomRegisterBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidRequestBody.Error(),
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		}

		return
	}

	registeredRoom, err := h.uc.RegisterVoiceRoom(ctx, guildId, originChannelId, body.ChannelId, body.CreatorId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrVoiceRoomExists.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusConflict)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[u.VoiceRoom]{
		Success: true,
		Data:    *registeredRoom,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/voice-room/{channel_id} [GET]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path	string	true	"The guild ID."
//	@Param		channel_id	path	string	true	"The guild ID."
//
// nolint:staticcheck
func (h *GuildHandler) GetVoiceRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	channelId := chi.URLParam(r, "channelId")

	room, err := h.uc.GetVoiceRoom(ctx, guildId, channelId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrVoiceRoomNotFound.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[u.VoiceRoom]{
		Success: true,
		Data:    *room,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/voice-room/{channel_id} [PATCH]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path	string	true	"The guild ID."
//	@Param		channel_id	path	string	true	"The guild ID."
//
// nolint:staticcheck
func (h *GuildHandler) UpdateVoiceRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	channelId := chi.URLParam(r, "channelId")
	var body *VoiceRoomModifyBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err := httpx.WriteJSON(w, APIError{
			Success: false,
			Message: ErrInvalidRequestBody.Error(),
		}, http.StatusBadRequest)

		if err != nil {
			log.Error(err)
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		}

		return
	}

	room, err := h.uc.UpdateVoiceRoom(ctx, guildId, channelId, u.VoiceRoomModify{
		CurrentOwnerId: body.CurrentOwnerId,
		IsLocked:       body.IsLocked,
	})

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrVoiceRoomNotFound.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[u.VoiceRoom]{
		Success: true,
		Data:    *room,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/voice-room/{channel_id} [DELETE]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path	string	true	"The guild ID."
//	@Param		channel_id	path	string	true	"The guild ID."
//
// nolint:staticcheck
func (h *GuildHandler) DeleteVoiceRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	channelId := chi.URLParam(r, "channelId")

	err := h.uc.DeleteVoiceRoom(ctx, guildId, channelId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		var ueErr u.UsecaseError
		if errors.As(err, &ueErr) {
			var writeErr error

			switch ueErr.Code {
			case u.ErrVoiceRoomNotFound.Code:
				writeErr = httpx.WriteJSON(w, APIError{
					Success: false,
					Code:    ueErr.Code,
					Message: ueErr.Message,
				}, http.StatusNotFound)
			}

			if writeErr != nil {
				log.Error(writeErr)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[any]{
		Success: true,
		Data:    nil,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}
