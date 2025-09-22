package handlers

import u "github.com/typical-developers/discord-bot-backend/internal/usecase"

type APIResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type APIError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GuildSettingsResponse APIResponse[u.GuildSettings]
