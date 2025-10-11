package usecase

import (
	"context"

	"maragu.dev/gomponents"
)

type MemberUsecase interface {
	CreateMemberProfile(ctx context.Context, guildId string, userId string) (*MemberProfile, error)
	GetMemberProfile(ctx context.Context, guildId string, userId string) (*MemberProfile, error)
	IncrementMemberChatActivityPoints(ctx context.Context, guildId string, userId string) (*MemberProfile, error)
	GenerateMemberProfileCard(ctx context.Context, guildId string, userId string) (gomponents.Node, error)
}
