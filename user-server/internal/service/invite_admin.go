package service

import (
	"context"
	"errors"
	"time"

	"hivemtk-user/internal/model"
	"hivemtk-user/internal/repository"
)

// invite_admin.go 邀请码业务（SaaS P2 准入；L3 分层合规——controller 不触 repository）。

// InviteAdminService 邀请码管理。
type InviteAdminService struct {
	repo *repository.InviteCodeRepo
}

func NewInviteAdminService() *InviteAdminService {
	r := repository.NewInviteCodeRepo()
	if !r.Available() {
		return &InviteAdminService{}
	}
	return &InviteAdminService{repo: r}
}

func (s *InviteAdminService) Ready() bool { return s.repo != nil }

// Generate 生成 12hex 邀请码（expiresDays<=0 = 永久）。
func (s *InviteAdminService) Generate(ctx context.Context, teamID uint, maxUses, expiresDays int, createdBy *uint) (*model.InviteCode, error) {
	if teamID == 0 {
		return nil, errors.New("无效团队 id")
	}
	if maxUses < 0 {
		maxUses = 0
	}
	var expiresAt *time.Time
	if expiresDays > 0 {
		t := time.Now().AddDate(0, 0, expiresDays)
		expiresAt = &t
	}
	// 生成逻辑（随机 12hex + 入账）在 repo.Generate——service 只做参数整形与归一。
	return s.repo.Generate(ctx, teamID, maxUses, expiresAt, createdBy)
}

func (s *InviteAdminService) ListByTeam(ctx context.Context, teamID uint) ([]model.InviteCode, error) {
	return s.repo.ListByTeam(ctx, teamID)
}

// Disable 停用（幂等）。
func (s *InviteAdminService) Disable(ctx context.Context, id uint) error {
	return s.repo.Disable(ctx, id)
}

// OwnsCode 团队管理员归属校验：码必须属于本团队。
func (s *InviteAdminService) OwnsCode(ctx context.Context, teamID, codeID uint) bool {
	rows, err := s.repo.ListByTeam(ctx, teamID)
	if err != nil {
		return false
	}
	for _, r := range rows {
		if r.ID == codeID {
			return true
		}
	}
	return false
}
