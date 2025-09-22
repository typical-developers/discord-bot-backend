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

type GuildActivityRoleCreateBody struct {
	ActivityType   string `json:"activity_type"`
	RoleID         string `json:"role_id"`
	RequiredPoints int32  `json:"required_points"`
}
