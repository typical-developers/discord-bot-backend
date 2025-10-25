package handlers

import u "github.com/typical-developers/discord-bot-backend/internal/usecase"

// --- Response Generics
type APIResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type APIError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// --- Guild Settings
type GuildSettingsResponse APIResponse[u.GuildSettings]

type GuildActivitySettingsUpdateBody u.UpdateAcitivtySettings

func (u GuildActivitySettingsUpdateBody) Validate() error {
	if u.ChatActivity == nil {
		return ErrInvalidRequestBody
	}

	return nil
}

type GuildActivityRoleCreateBody struct {
	ActivityType   string `json:"activity_type"`
	RoleID         string `json:"role_id"`
	RequiredPoints int32  `json:"required_points"`
}

// --- Voice Rooms
type VoiceRoomLobbySettings u.VoiceRoomLobbySettings

type VoiceRoomRegisterBody u.VoiceRoomRegister

type VoiceRoomModifyBody u.VoiceRoomModify

// --- Member Profile
type MigrateMemberProfileBody u.MigrateMemberProfile

func (m MigrateMemberProfileBody) Validate() error {
	if m.ToMemberId == "" {
		return ErrInvalidRequestBody
	}

	return nil
}

type MemberProfileResponse APIResponse[u.MemberProfile]
