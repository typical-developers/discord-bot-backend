package api_structures

type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ActivityRole struct {
	RoleID         string `json:"role_id"`
	RequiredPoints int    `json:"required_points"`
}

type Activity struct {
	IsEnabled       bool           `json:"is_enabled"`
	GrantAmount     int            `json:"grant_amount"`
	CooldownSeconds int            `json:"cooldown_seconds"`
	ActivityRoles   []ActivityRole `json:"activity_roles"`
}

type GuildSettings struct {
	ChatActivity Activity `json:"chat_activity"`
}
