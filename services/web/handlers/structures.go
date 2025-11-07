package handlers

import u "github.com/typical-developers/discord-bot-backend/internal/usecase"

// --- Response Generics
type APIResponse[T any] struct {
	Data T `json:"data"`
}

type APIError struct {
	Code    string `json:"code"`
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

type GuildMessageEmbedSettingsUpdateBody u.UpdateMessageEmbedSettingsOpts

func (u GuildMessageEmbedSettingsUpdateBody) Validate() error {
	if u.IsEnabled == nil && u.AddDisabledChannel == nil && u.AddIgnoredChannel == nil && u.AddIgnoredRole == nil && u.RemoveDisabledChannel == nil && u.RemoveIgnoredChannel == nil && u.RemoveIgnoredRole == nil {
		return ErrInvalidRequestBody
	}

	return nil
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
