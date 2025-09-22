package usecase

import "errors"

var (
	ErrGuildNotFound       = errors.New("guild does not exist")
	ErrGuildSettingsExists = errors.New("guild settings already exists")
	ErrActivityRoleExists  = errors.New("activity role already exists")
)
