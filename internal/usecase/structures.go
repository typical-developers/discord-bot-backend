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

type MemberActivityRole struct {
	RoleID         string `json:"role_id"`
	Accent         string `json:"accent"`
	Name           string `json:"name"`
	RequiredPoints int32  `json:"required_points"`
}

type MemberActivityProgress struct {
	CurrentProgress  int32 `json:"current_progress"`
	RequiredProgress int32 `json:"required_progress"`
}

type MemberActivity struct {
	Rank         int32 `json:"rank"`
	Points       int32 `json:"total_points"`
	IsOnCooldown bool  `json:"is_on_cooldown"`

	CurrentActivityRoleIds []string                `json:"current_activity_role_ids"`
	CurrentActivityRole    *MemberActivityRole     `json:"current_activity_role"`
	NextActivityRole       *MemberActivityProgress `json:"next_activity_role"`
}

type MemberProfile struct {
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar_url"`

	CardStyle    int32          `json:"card_style"`
	ChatActivity MemberActivity `json:"chat_activity"`
}

type MigrateMemberProfile struct {
	ToMemberId string `json:"to_member_id"`
}
