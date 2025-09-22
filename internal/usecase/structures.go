package usecase

type GuildActivityRole struct {
	RoleID         string `json:"role_id"`
	RequiredPoints int32  `json:"required_points"`
}

type GuildActivityTracking struct {
	IsEnabled       bool                `json:"is_enabled"`
	GrantAmount     int32               `json:"grant_amount"`
	CooldownSeconds int32               `json:"cooldown"`
	ActivityRoles   []GuildActivityRole `json:"activity_roles"`
	DenyRoles       []string            `json:"deny_roles"`
}

type VoiceRoomLobby struct {
	ChannelID      string `json:"channel_id"`
	UserLimit      int32  `json:"user_limit"`
	CanRename      bool   `json:"can_rename"`
	CanLock        bool   `json:"can_lock"`
	CanAdjustLimit bool   `json:"can_adjust_limit"`
}

type GuildSettings struct {
	ChatActivityTracking GuildActivityTracking `json:"chat_activity"`
	VoiceRoomLobbies     []VoiceRoomLobby      `json:"voice_room_lobbies"`
}

type UpdateActivitySettingsOpts struct {
	IsEnabled       *bool  `json:"is_enabled"`
	GrantAmount     *int32 `json:"grant_amount"`
	CooldownSeconds *int32 `json:"cooldown"`
}

type UpdateAcitivtySettings struct {
	ChatActivity UpdateActivitySettingsOpts `json:"chat_activity"`
}
