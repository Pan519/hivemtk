package model

import "time"

// Tenant 租户（多租户 P1 实体）：项目可部署多个租户，每个租户管理多个团队。
// 权限层级：项目超管（admin）= 全量；租户超管（tenant_admin）= 本租户全部团队；
// 团队管理员（team_admin）= 本团队。
type Tenant struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Code string `gorm:"column:code;type:varchar(64);uniqueIndex" json:"code"`
	Name string `gorm:"type:varchar(128);not null" json:"name"`
	// Status 生命周期：active / disabled（disabled 后名下团队业务写入面在 P2 收紧）
	Status string `gorm:"type:varchar(16);default:'active';index" json:"status"`
	// OwnerUserID 租户主管理员（system_users.id，可 NULL=暂无）
	OwnerUserID *uint `gorm:"column:owner_user_id;index" json:"owner_user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TenantStatusActive / Disabled 租户生命周期常量。
const (
	TenantStatusActive   = "active"
	TenantStatusDisabled = "disabled"
)

// TableName 显式表名。
func (*Tenant) TableName() string { return "tenants" }