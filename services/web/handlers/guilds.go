package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

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

		r.Get("/activity-leaderboard-card", h.GenerateGuildActivityLeaderboardCard)
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
	settings, err := h.uc.CreateGuildSettings(ctx, guildId)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		if errors.Is(err, u.ErrGuildSettingsExists) {
			err := httpx.WriteJSON(w, APIError{
				Success: false,
				Message: err.Error(),
			}, http.StatusConflict)

			if err != nil {
				log.Error(err)
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

		if errors.Is(err, u.ErrGuildNotFound) {
			err := httpx.WriteJSON(w, APIError{
				Success: false,
				Message: err.Error(),
			}, http.StatusNotFound)

			if err != nil {
				log.Error(err)
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

		if errors.Is(err, u.ErrActivityRoleExists) {
			err := httpx.WriteJSON(w, APIError{
				Success: false,
				Message: err.Error(),
			}, http.StatusConflict)

			if err != nil {
				log.Error(err)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		log.Error(err)
		http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
		return
	}

	err = httpx.WriteJSON(w, APIResponse[struct{}]{
		Success: true,
		Data:    struct{}{},
	}, http.StatusCreated)
	if err != nil {
		log.Error(err)
	}
}

//	@Router		/v1/guild/{guild_id}/activity-leaderboard-card [GET]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id		path	string	true	"The guild ID."
//	@Param		activity_type	query	string	true	"The activity type."	Enum(chat, voice)
//	@Param		time_period		query	string	true	"The time period."		Enum(week, month, all)
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
