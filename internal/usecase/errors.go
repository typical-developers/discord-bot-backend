package usecase

import "errors"

var (
	ErrGuildNotFound = errors.New("guild does not exist")

	ErrMemberNotInGuild             = errors.New("member is not in guild")
	ErrMemberProfileNotFound        = errors.New("member profile not found")
	ErrMemberProfileExists          = errors.New("member profile already exists")
	ErrMemberOnGrantCooldown        = errors.New("member on grant cooldown")
	ErrChatActivityTrackingDisabled = errors.New("chat activity tracking is disabled")

	ErrGuildSettingsExists = errors.New("guild settings already exists")
	ErrActivityRoleExists  = errors.New("activity role already exists")
)
