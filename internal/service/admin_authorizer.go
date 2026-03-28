package service

import (
	"context"
	"errors"

	"github.com/casbin/casbin/v3"
	casbinmodel "github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	AdminDomain        = "admin"
	AdminRoleAdmin     = "admin"
	AdminRoleAssistant = "assistant"
)

const adminRBACModel = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
`

var defaultAdminPolicies = [][]string{
	{AdminRoleAdmin, AdminDomain, "dashboard", "read"},
	{AdminRoleAdmin, AdminDomain, "dashboard", "write"},
	{AdminRoleAdmin, AdminDomain, "problems", "read"},
	{AdminRoleAdmin, AdminDomain, "problems", "write"},
	{AdminRoleAdmin, AdminDomain, "problem_sets", "read"},
	{AdminRoleAdmin, AdminDomain, "problem_sets", "write"},
	{AdminRoleAdmin, AdminDomain, "tags", "read"},
	{AdminRoleAdmin, AdminDomain, "tags", "write"},
	{AdminRoleAdmin, AdminDomain, "test_cases", "read"},
	{AdminRoleAdmin, AdminDomain, "test_cases", "write"},
	{AdminRoleAdmin, AdminDomain, "judge_configs", "read"},
	{AdminRoleAdmin, AdminDomain, "judge_configs", "write"},
	{AdminRoleAdmin, AdminDomain, "submissions", "read"},
	{AdminRoleAdmin, AdminDomain, "submissions", "write"},
	{AdminRoleAdmin, AdminDomain, "users", "read"},
	{AdminRoleAdmin, AdminDomain, "users", "write"},
	{AdminRoleAdmin, AdminDomain, "settings", "read"},
	{AdminRoleAdmin, AdminDomain, "settings", "write"},
	{AdminRoleAssistant, AdminDomain, "dashboard", "read"},
	{AdminRoleAssistant, AdminDomain, "problems", "read"},
	{AdminRoleAssistant, AdminDomain, "problems", "write"},
	{AdminRoleAssistant, AdminDomain, "problem_sets", "read"},
	{AdminRoleAssistant, AdminDomain, "problem_sets", "write"},
	{AdminRoleAssistant, AdminDomain, "tags", "read"},
	{AdminRoleAssistant, AdminDomain, "tags", "write"},
	{AdminRoleAssistant, AdminDomain, "test_cases", "read"},
	{AdminRoleAssistant, AdminDomain, "test_cases", "write"},
	{AdminRoleAssistant, AdminDomain, "submissions", "read"},
}

type AdminAuthorizer struct {
	enforcer *casbin.Enforcer
}

func NewInMemoryAdminAuthorizer() (*AdminAuthorizer, error) {
	modelInstance, err := casbinmodel.NewModelFromString(adminRBACModel)
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(modelInstance)
	if err != nil {
		return nil, err
	}

	authz := &AdminAuthorizer{enforcer: enforcer}
	if err := authz.SeedPolicies(context.Background()); err != nil {
		return nil, err
	}
	return authz, nil
}

func NewAdminAuthorizer(db *gorm.DB) (*AdminAuthorizer, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}

	modelInstance, err := casbinmodel.NewModelFromString(adminRBACModel)
	if err != nil {
		return nil, err
	}
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", "casbin_rule")
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(modelInstance, adapter)
	if err != nil {
		return nil, err
	}
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return &AdminAuthorizer{enforcer: enforcer}, nil
}

func (a *AdminAuthorizer) SeedPolicies(ctx context.Context) error {
	if a == nil || a.enforcer == nil {
		return errors.New("authorizer is required")
	}

	if _, err := a.enforcer.AddPolicies(defaultAdminPolicies); err != nil {
		return err
	}
	_ = ctx
	return nil
}

func (a *AdminAuthorizer) AssignRole(userID string, role string) error {
	if a == nil || a.enforcer == nil {
		return errors.New("authorizer is required")
	}

	_, err := a.enforcer.AddRoleForUserInDomain(userID, role, AdminDomain)
	return err
}

func (a *AdminAuthorizer) Enforce(sub string, dom string, obj string, act string) (bool, error) {
	if a == nil || a.enforcer == nil {
		return false, errors.New("authorizer is required")
	}

	return a.enforcer.Enforce(sub, dom, obj, act)
}
