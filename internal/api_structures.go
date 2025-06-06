package models

import (
	"time"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type ActivityRole struct {
	RoleID         string `json:"role_id"`
	RequiredPoints int    `json:"required_points"`
}

type ActivityConfig struct {
	IsEnabled       bool           `json:"is_enabled"`
	GrantAmount     int            `json:"grant_amount"`
	CooldownSeconds int            `json:"cooldown_seconds"`
	ActivityRoles   []ActivityRole `json:"activity_roles"`
}

type GuildSettings struct {
	ChatActivity ActivityConfig `json:"chat_activity"`
}

// ---------------------------------------------------------------------

type CardStyle int

const (
	CardStyleDefault CardStyle = iota
)

type ActivityType string

const (
	ActivityTypeChat ActivityType = "chat"
	// ActivityTypeVoice ActivityType = "voice"
)

type LeaderboardType string

const (
	LeaderboardTypeAllTime LeaderboardType = "all"
	LeaderboardTypeMonthly LeaderboardType = "monthly"
	LeaderboardTypeWeekly  LeaderboardType = "weekly"
)

func (l LeaderboardType) Valid() bool {
	return l == LeaderboardTypeAllTime || l == LeaderboardTypeMonthly || l == LeaderboardTypeWeekly
}

func (a ActivityType) Valid() bool {
	return a == ActivityTypeChat
}

type ActivityRoleProgress struct {
	RoleID         string `json:"role_id"`
	Progress       int    `json:"progress"`
	RequiredPoints int    `json:"required_points"`
}

type MemberRoles struct {
	Next     *ActivityRoleProgress `json:"next"`
	Obtained []ActivityRole        `json:"obtained"`
}

type MemberActivity struct {
	Rank         int         `json:"rank"`
	Points       int         `json:"points"`
	IsOnCooldown bool        `json:"is_on_cooldown"`
	LastGrant    time.Time   `json:"last_grant"`
	Roles        MemberRoles `json:"roles"`
}

type MemberProfile struct {
	CardStyle    CardStyle      `json:"card_style"`
	ChatActivity MemberActivity `json:"chat_activity"`
}

type APIResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data,omitempty"`
}

// ---------------------------------------------------------------------

type ChatActivityRoleQuery struct {
	Title  string `json:"title"`
	Accent string `json:"accent"`
}

// ---------------------------------------------------------------------

type ActivitySettings struct {
	Enabled     *bool  `json:"enabled,omitempty"`
	Cooldown    *int32 `json:"cooldown,omitempty"`
	GrantAmount *int32 `json:"grant_amount,omitempty"`
}

type UpdateActivitySettings struct {
	ChatActivity ActivitySettings `json:"chat_activity"`
}

type AddActivityRole struct {
	GrantType      ActivityType `json:"grant_type"`
	RoleID         string       `json:"role_id"`
	RequiredPoints int          `json:"required_points"`
}
