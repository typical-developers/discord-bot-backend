package usecase

type UsecaseError struct {
	Code    string
	Message string
}

func NewUsecaseError(code, message string) UsecaseError {
	return UsecaseError{Code: code, Message: message}
}

func (e UsecaseError) Error() string {
	return e.Message
}

var (
	// Guild Setting Errors
	ErrGuildSettingsExists          = NewUsecaseError("GUILD_ALREADY_EXISTS", "the guild already exists.")
	ErrGuildNotFound                = NewUsecaseError("GUILD_NOT_FOUND", "the guild was not found.")
	ErrChatActivityTrackingDisabled = NewUsecaseError("CHAT_ACTIVITY_TRACKING_DISABLED", "chat activity tracking is disabled.")
	ErrActivityRoleExists           = NewUsecaseError("ACTIVITY_ROLE_ALREADY_EXISTS", "the activity role already exists.")

	// Member Errors
	ErrMemberNotInGuild      = NewUsecaseError("MEMBER_NOT_IN_GUILD", "the member is not in the guild.")
	ErrMemberProfileNotFound = NewUsecaseError("MEMBER_NOT_FOUND", "the member profile was not found.")
	ErrMemberProfileExists   = NewUsecaseError("MEMBER_ALREADY_EXISTS", "the member profile already exists.")
	ErrMemberOnGrantCooldown = NewUsecaseError("MEMBER_ON_COOLDOWN", "the member is on cooldown.")

	// Leaderboard Errors
	ErrLeaderboardNoRows = NewUsecaseError("LEADERBOARD_NO_ROWS", "the leaderboard has no rows.")

	// Voice Room Errors
	ErrVoiceRoomLobbyExists      = NewUsecaseError("VOICE_ROOM_LOBBY_ALREADY_EXISTS", "the voice room lobby already exists.")
	ErrVoiceRoomLobbyNotFound    = NewUsecaseError("VOICE_ROOM_LOBBY_NOT_FOUND", "the voice room lobby was not found.")
	ErrVoiceRoomLobbyIsVoiceRoom = NewUsecaseError("VOICE_ROOM_LOBBY_IS_ACTIVE_VOICE_ROOM", "the voice room lobby is already an active voice room.")
	ErrVoiceRoomExists           = NewUsecaseError("VOICE_ROOM_EXISTS", "the voice room already exists.")
	ErrVoiceRoomNotFound         = NewUsecaseError("VOICE_ROOM_NOT_FOUND", "the voice room was not found.")
)
