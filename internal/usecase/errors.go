package usecase

import "errors"

var (
	ErrInvalidRequestBody = errors.New("INVALID_REQUEST")

	ErrGuildSettingsExists = errors.New("GUILD_ALREADY_EXISTS")
	ErrGuildNotFound       = errors.New("GUILD_NOT_FOUND")

	ErrMemberNotInGuild             = errors.New("MEMBER_NOT_IN_GUILD")
	ErrMemberProfileNotFound        = errors.New("MEMBER_NOT_FOUND")
	ErrMemberProfileExists          = errors.New("MEMBER_ALREADY_EXISTS")
	ErrMemberOnGrantCooldown        = errors.New("MEMBER_ON_COOLDOWN")
	ErrChatActivityTrackingDisabled = errors.New("CHAT_ACTIVITY_TRACKING_DISABLED")
	ErrActivityRoleExists           = errors.New("ACTIVITY_ROLE_ALREADY_EXISTS")

	ErrVoiceRoomLobbyExists      = errors.New("VOICE_ROOM_LOBBY_ALREADY_EXISTS")
	ErrVoiceRoomLobbyNotFound    = errors.New("VOICE_ROOM_LOBBY_NOT_FOUND")
	ErrVoiceRoomLobbyIsVoiceRoom = errors.New("VOICE_ROOM_LOBBY_IS_ACTIVE_VOICE_ROOM")
	ErrVoiceRoomExists           = errors.New("VOICE_ROOM_EXISTS")
	ErrVoiceRoomNotFound         = errors.New("VOICE_ROOM_NOT_FOUND")
)
