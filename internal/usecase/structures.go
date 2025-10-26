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

	OpenedRooms []string `json:"opened_rooms"`
}

type MessageEmbeds struct {
	IsEnabled        bool     `json:"is_enabled"`
	DisabledChannels []string `json:"disabled_channels"`
	IgnoredChannels  []string `json:"ignored_channels"`
	IgnoredRoles     []string `json:"ignored_roles"`
}

type GuildSettings struct {
	ChatActivityTracking GuildActivityTracking `json:"chat_activity"`
	MessageEmbeds        MessageEmbeds         `json:"message_embeds"`
	VoiceRoomLobbies     []VoiceRoomLobby      `json:"voice_room_lobbies"`
}

type UpdateActivitySettingsOpts struct {
	IsEnabled       *bool  `json:"is_enabled"`
	GrantAmount     *int32 `json:"grant_amount"`
	CooldownSeconds *int32 `json:"cooldown"`
}

type UpdateMessageEmbedSettingsOpts struct {
	IsEnabled             *bool   `json:"is_enabled"`
	AddDisabledChannel    *string `json:"add_disabled_channel"`
	RemoveDisabledChannel *string `json:"remove_disabled_channel"`
	AddIgnoredChannel     *string `json:"add_ignored_channel"`
	RemoveIgnoredChannel  *string `json:"remove_ignored_channel"`
	AddIgnoredRole        *string `json:"add_ignored_role"`
	RemoveIgnoredRole     *string `json:"remove_ignored_role"`
}

type UpdateAcitivtySettings struct {
	ChatActivity *UpdateActivitySettingsOpts `json:"chat_activity"`
}

type VoiceRoomRegister struct {
	CreatorId string `json:"creator_id"`
	ChannelId string `json:"channel_id"`
}

type VoiceRoomModify struct {
	CurrentOwnerId *string `json:"current_owner_id"`
	IsLocked       *bool   `json:"is_locked"`
}

type VoiceRoomLobbySettings struct {
	UserLimit      *int32 `json:"user_limit"`
	CanRename      *bool  `json:"can_rename"`
	CanLock        *bool  `json:"can_lock"`
	CanAdjustLimit *bool  `json:"can_adjust_limit"`
}

type VoiceRoom struct {
	OriginChannelId string `json:"origin_channel_id"`
	CreatorId       string `json:"creator_id"`
	CurrentOwnerId  string `json:"current_owner_id"`
	IsLocked        bool   `json:"is_locked"`

	Settings VoiceRoomLobbySettings `json:"settings"`
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
