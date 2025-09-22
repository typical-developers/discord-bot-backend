package usecase

type GuildActivityRole struct {
	RoleID         string `json:"role_id"`
	RequiredPoints int32  `json:"required_points"`
}

type GuildActivityTracking struct {
	Enabled    bool                `json:"enabled"`
	Cooldown   int32               `json:"cooldown"`
	Grant      int32               `json:"grant"`
	GrantRoles []GuildActivityRole `json:"grant_roles"`
	DenyRoles  []string            `json:"deny_roles"`
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
