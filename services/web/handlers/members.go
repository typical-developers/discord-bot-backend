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

type MemberHandler struct {
	uc u.MemberUsecase
}

func NewMemberHandler(r *chi.Mux, uc u.MemberUsecase) {
	h := MemberHandler{uc: uc}

	r.Route("/v1/guild/{guildId}/member/{memberId}", func(r chi.Router) {
		r.Post("/", h.CreateMemberProfile)
		r.Get("/", h.GetMemberProfile)
		r.Get("/profile-card", h.GenerateMemberProfileCard)
		r.Patch("/chat-activity", h.IncrementMemberChatActivityPoints)
		r.Post("/migrate", h.MigrateMemberProfile)
	})
}

//	@Router	/v1/guild/{guild_id}/member/{member_id} [POST]
//	@Tags	Members
//
//	@Param	guild_id	path	string	true	"The guild ID."
//	@Param	member_id	path	string	true	"The member ID."
//
// nolint:staticcheck
func (h *MemberHandler) CreateMemberProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	guildId := chi.URLParam(r, "guildId")
	memberId := chi.URLParam(r, "memberId")

	profile, err := h.uc.CreateMemberProfile(ctx, guildId, memberId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		if errors.Is(err, u.ErrMemberProfileExists) {
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

		if errors.Is(err, u.ErrMemberNotInGuild) {
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

	err = httpx.WriteJSON(w, MemberProfileResponse{
		Success: true,
		Data:    *profile,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router	/v1/guild/{guild_id}/member/{member_id} [GET]
//	@Tags	Members
//
//	@Param	guild_id	path	string	true	"The guild ID."
//	@Param	member_id	path	string	true	"The member ID."
//
// nolint:staticcheck
func (h *MemberHandler) GetMemberProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	guildId := chi.URLParam(r, "guildId")
	memberId := chi.URLParam(r, "memberId")

	profile, err := h.uc.GetMemberProfile(ctx, guildId, memberId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		if errors.Is(err, u.ErrMemberProfileNotFound) {
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

		if errors.Is(err, u.ErrMemberNotInGuild) {
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

	err = httpx.WriteJSON(w, MemberProfileResponse{
		Success: true,
		Data:    *profile,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router	/v1/guild/{guild_id}/member/{member_id}/profile-card [GET]
//	@Tags	Members
//
//	@Param	guild_id	path	string	true	"The guild ID."
//	@Param	member_id	path	string	true	"The member ID."
//
// nolint:staticcheck
func (h *MemberHandler) GenerateMemberProfileCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	guildId := chi.URLParam(r, "guildId")
	memberId := chi.URLParam(r, "memberId")

	card, err := h.uc.GenerateMemberProfileCard(ctx, guildId, memberId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		if errors.Is(err, u.ErrMemberProfileNotFound) {
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

		if errors.Is(err, u.ErrMemberNotInGuild) {
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := card.Render(w); err != nil {
		http.Error(w, "Failed to render profile card", http.StatusInternalServerError)
		return
	}
}

//	@Router	/v1/guild/{guild_id}/member/{member_id}/chat-activity [PATCH]
//	@Tags	Members
//
//	@Param	guild_id	path	string	true	"The guild ID."
//	@Param	member_id	path	string	true	"The member ID."
//
// nolint:staticcheck
func (h *MemberHandler) IncrementMemberChatActivityPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	guildId := chi.URLParam(r, "guildId")
	memberId := chi.URLParam(r, "memberId")

	profile, err := h.uc.IncrementMemberChatActivityPoints(ctx, guildId, memberId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		if errors.Is(err, u.ErrMemberProfileNotFound) {
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

		if errors.Is(err, u.ErrMemberOnGrantCooldown) {
			err := httpx.WriteJSON(w, APIError{
				Success: false,
				Message: err.Error(),
			}, http.StatusTooManyRequests)

			if err != nil {
				log.Error(err)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		if errors.Is(err, u.ErrChatActivityTrackingDisabled) {
			err := httpx.WriteJSON(w, APIError{
				Success: false,
				Message: err.Error(),
			}, http.StatusForbidden)

			if err != nil {
				log.Error(err)
				http.Error(w, ErrInternalError.Error(), http.StatusInternalServerError)
			}

			return
		}

		if errors.Is(err, u.ErrMemberNotInGuild) {
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

	err = httpx.WriteJSON(w, MemberProfileResponse{
		Success: true,
		Data:    *profile,
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}

//	@Router	/v1/guild/{guild_id}/member/{member_id}/migrate [POST]
//	@Tags	Members
//
//	@Param	guild_id	path	string						true	"The guild ID."
//	@Param	member_id	path	string						true	"The member ID."
//
//	@Param	body		body	MigrateMemberProfileBody	true	"The migration body."
//
// nolint:staticcheck
func (h *MemberHandler) MigrateMemberProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	guildId := chi.URLParam(r, "guildId")
	memberId := chi.URLParam(r, "memberId")

	var body *MigrateMemberProfileBody
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

	if memberId == body.ToMemberId {
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

	err := h.uc.MigrateMemberProfile(ctx, guildId, memberId, body.ToMemberId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, ErrGatewayTimeout.Error(), http.StatusGatewayTimeout)
			return
		}

		if errors.Is(err, u.ErrMemberProfileNotFound) {
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

	err = httpx.WriteJSON(w, APIResponse[struct{}]{
		Success: true,
		Data:    struct{}{},
	}, http.StatusOK)
	if err != nil {
		log.Error(err)
	}
}
