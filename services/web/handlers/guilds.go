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

	r.Route("/guild", func(r chi.Router) {
		r.Post("/{guildId}/settings/create", h.CreateGuildSettings)
		r.Get("/{guildId}/settings", h.GetGuildSettings)
		r.Patch("/{guildId}/settings/update/activity", h.UpdateGuildActivitySettings)
	})
}

//	@Router		/guild/{guild_id}/settings/create [post]
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

		if errors.Is(err, u.ErrGuildSettingsExist) {
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

//	@Router		/guild/{guild_id}/settings [get]
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

//	@Router		/guild/{guild_id}/settings/update/activity  [patch]
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
