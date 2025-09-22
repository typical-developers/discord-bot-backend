package usecase

import "errors"

var (
	ErrGuildNotFound      = errors.New("guild does not exist")
	ErrGuildSettingsExist = errors.New("guild settings already exist")
)
