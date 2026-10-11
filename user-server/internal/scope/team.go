// team.go 团队维度隔离 scope（多租户 P0）。
//
// 与 tenant.go 的 owner 维度正交：owner 管"谁创建的"，team 管"哪个租户的"。
// 数据流：JWT 中间件把 team_id/role 写进请求 ctx（WithTeamInfo），
// repository 用 TeamScoped(ctx) 追加过滤，Create 由 db 包的租户回调回填。
//
// 旁路开关：TENANT_ISOLATION=off/0/false 时全部直通（回滚口径，见
// scripts/rollback-tenant-p0.sh）——关闭后系统回到 2026-10-11 之前的全共享行为。
package scope

import (
	"context"
	"os"

	"gorm.io/gorm"
)

// teamInfoKey 租户上下文在 ctx 里的单一承载键（struct 值，presence 即"有身份"）。
const teamInfoKey = "team_info"

// tenantIsolationEnabled 可测试替换（seam）；默认读环境开关。
// 缺省（未配置）= 开启：对存量部署安全——回填钩子已把所有人归入默认团队，
// 过滤结果与升级前一致；而"忘了配开关导致裸奔"的代价远大于"多配一个开关"。
var tenantIsolationEnabled = func() bool {
	v := os.Getenv("TENANT_ISOLATION")
	return v != "off" && v != "0" && v != "false"
}

// SetTenantIsolationForTest 仅供测试切换开关（成对还原）。
func SetTenantIsolationForTest(enabled bool) (restore func()) {
	old := tenantIsolationEnabled
	tenantIsolationEnabled = func() bool { return enabled }
	return func() { tenantIsolationEnabled = old }
}

// TeamInfo 请求侧租户上下文。
type TeamInfo struct {
	// TeamID 当前用户归属团队（0=平台侧/未归属）
	TeamID uint
	// TenantID 当前用户归属租户（多租户 P1）：0=平台侧/项目超管；1+ 归属租户。
	// tenant_admin 角色的数据范围 = 本租户全部团队。
	TenantID uint
	// Platform 平台管理员直通：role==admin（超级管理员跨团队全量，需求 3）
	Platform bool
}

// WithTeamInfo 把租户上下文注入 ctx（JWT 中间件/MCP guard 调用）。
// 注入即"在场"：之后的 TeamScoped/回调按 info 过滤与回填。
func WithTeamInfo(ctx context.Context, info TeamInfo) context.Context {
	if ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, ginCtxKey(teamInfoKey), info)
}

// TeamInfoFromCtx 取租户上下文；present=false 表示**没有任何身份注入**——
// 这是系统上下文（定时任务 / ingest 管线 / bridge 等平台机器流量），
// 与"已注入的非平台身份"（team 过滤）和"已注入的平台管理员"（直通）是三种不同语义：
// 系统上下文读直通、写落默认团队；非平台身份读写都锁本团队（fail-closed）。
func TeamInfoFromCtx(ctx context.Context) (TeamInfo, bool) {
	if ctx == nil {
		return TeamInfo{}, false
	}
	if v, ok := fromGinContext(ctx, teamInfoKey); ok {
		if info, ok2 := v.(TeamInfo); ok2 {
			return info, true
		}
	}
	// 兼容链路：scope.SetGinValuesToCtx 只注入了 user_id/role（老调用方）——
	// 有 user 即有身份：按"非平台、无团队"处理（fail-closed，只可能少看）。
	if _, ok := fromGinContext(ctx, ginUserKey); ok {
		return TeamInfo{}, true
	}
	return TeamInfo{}, false
}

// TeamScoped 团队隔离 scope（P0 红线：核心业务表查询必须经过它）。
//
// 语义：
//   - TENANT_ISOLATION=off → 直通（回滚口径）；
//   - 平台管理员（role=admin）→ 直通（跨团队全量管理，需求 3）；
//   - 其余 → team_id = ctx 团队；无团队身份兜底只看 team_id=0（平台侧行），
//     这是 fail-closed：漏注入 ctx 的调用方看不到别人的数据，只可能少看。
func TeamScoped(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if db == nil || !tenantIsolationEnabled() {
			return db
		}
		passthrough, tenantID, teamID := teamFilterDecision(ctx)
		if passthrough {
			return db
		}
		if tenantID > 0 {
			// 租户超管：team_id ∈ 本租户全部团队（子查询 teams 表，scope 不 import model 用字面表名）
			return db.Where("team_id IN (?)", db.Table("teams").Where("tenant_id = ?", tenantID).Select("id"))
		}
		return db.Where("team_id = ?", teamID)
	}
}

// TenantTeamScoped 团队表的租户范围过滤（多租户 P1）：按 ctx 身份层级过滤 teams 表。
//   - 系统上下文 / 平台管理员（admin）→ 全量直通；
//   - 租户超管（tenant_admin）→ WHERE tenant_id = ctx 租户；
//   - 其余（team_admin / 成员）→ WHERE id = 本团队。
func TenantTeamScoped(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if db == nil || !tenantIsolationEnabled() {
			return db
		}
		info, present := TeamInfoFromCtx(ctx)
		if !present {
			return db
		}
		if info.Platform || isAdmin(ctx) {
			return db
		}
		if isTenantAdmin(ctx) {
			return db.Where("tenant_id = ?", info.TenantID)
		}
		return db.Where("id = ?", info.TeamID)
	}
}

// teamFilterDecision 租户过滤的四态判定（纯函数，供 TeamScoped 与单测复用）：
//   - (true, 0, 0)：直通——系统上下文（无身份）或平台管理员（admin）；
//   - (false, tenantID>0, 0)：租户超管（tenant_admin）——锁本租户团队集合；
//   - (false, 0, teamID)：其余——锁本团队（fail-closed，漏注入等于少看不是多看）。
func teamFilterDecision(ctx context.Context) (passthrough bool, tenantID, teamID uint) {
	info, present := TeamInfoFromCtx(ctx)
	if !present {
		return true, 0, 0 // 系统上下文（无身份）：平台机器流量直通（读全量）
	}
	if info.Platform || isAdmin(ctx) {
		return true, 0, 0 // 平台管理员：跨团队全量（需求 3）
	}
	if isTenantAdmin(ctx) {
		return false, info.TenantID, 0 // 租户超管：本租户全部团队
	}
	return false, 0, info.TeamID
}

// isTenantAdmin 判断 ctx 身份是否为租户超管（tenant_admin）。
func isTenantAdmin(ctx context.Context) bool {
	if v, ok := fromGinContext(ctx, roleKey); ok {
		if s, ok2 := v.(string); ok2 {
			return s == "tenant_admin"
		}
	}
	return false
}

// roleKey ctx 中的角色键（JWT 中间件 Set("role", ...)）。
const roleKey = "role"

// FillTeamIDOnCreate 建行前决定 team_id（db 包 GORM 回调调用）：
//   - 隔离关闭 → 不动（全 0，回滚口径）；
//   - 已注入身份：平台管理员 → 0（平台侧行）；普通身份 → 其团队；
//   - 系统上下文（无身份）→ resolveDefault()（默认团队），保证 ingest 等
//     平台管线新建的业务行对团队成员可见，而不是沉进谁也看不见的 team 0。
//
// 显式指定过 team_id 的调用方不受影响（回调层判零值）。
func FillTeamIDOnCreate(ctx context.Context, resolveDefault func() uint, setField func(uint) error) error {
	if !tenantIsolationEnabled() {
		return nil
	}
	info, present := TeamInfoFromCtx(ctx)
	if present {
		if info.Platform {
			return nil
		}
		return setField(info.TeamID)
	}
	if resolveDefault == nil {
		return nil
	}
	return setField(resolveDefault())
}
