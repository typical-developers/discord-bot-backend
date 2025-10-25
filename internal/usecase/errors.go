package usecase

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")

	ErrGuildSettingsExists = errors.New("guild settings already exists")
	ErrGuildNotFound       = errors.New("guild does not exist")

	ErrMemberNotInGuild             = errors.New("member is not in guild")
	ErrMemberProfileNotFound        = errors.New("member profile not found")
	ErrMemberProfileExists          = errors.New("member profile already exists")
	ErrMemberOnGrantCooldown        = errors.New("member on grant cooldown")
	ErrChatActivityTrackingDisabled = errors.New("chat activity tracking is disabled")
	ErrActivityRoleExists           = errors.New("activity role already exists")

	ErrVoiceRoomLobbyExists   = errors.New("voice room lobby already exists")
	ErrVoiceRoomLobbyNotFound = errors.New("voice room lobby does not exist")
	ErrVoiceRoomExists        = errors.New("voice room already exists")
	ErrVoiceRoomNotFound      = errors.New("voice room does not exist")
)
