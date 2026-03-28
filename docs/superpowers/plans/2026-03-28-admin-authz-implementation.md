# Admin Bootstrap And Dashboard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add default super-admin bootstrap, a usable `/admin` dashboard API, and a frontend admin login/dashboard flow that can log in successfully end-to-end.

**Architecture:** Keep the existing user-facing auth flow unchanged, then extend the existing admin auth stack with a bootstrap initializer that seeds a default super-admin from config or fallback defaults. Reuse the current Vue frontend project, but add a separate admin auth store, admin API client helpers, and protected admin routes so the admin login token never collides with the user token.

**Tech Stack:** Go, Gin, Gorm, PostgreSQL, Casbin, Vue 3, Pinia, Vue Router, Naive UI, Vitest

---

### Task 1: Add Bootstrap Admin Configuration And Seeding

**Files:**
- Modify: `internal/config/config.go`
- Modify: `internal/bootstrap/app.go`
- Modify: `internal/repository/admin_user.go`
- Test: `internal/bootstrap/app_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestEnsureBootstrapAdminCreatesFallbackSuperAdmin(t *testing.T) {
	repo := &stubBootstrapAdminUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*model.AdminUser, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	cfg := &config.Config{}
	err := ensureBootstrapAdmin(context.Background(), repo, cfg)

	if err != nil {
		t.Fatalf("ensureBootstrapAdmin() error = %v", err)
	}
	if repo.createdUser == nil {
		t.Fatal("expected bootstrap admin to be created")
	}
	if repo.createdUser.Email != "admin@go-oj.dev" {
		t.Fatalf("expected fallback email, got %q", repo.createdUser.Email)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/bootstrap -run TestEnsureBootstrapAdminCreatesFallbackSuperAdmin -v`
Expected: FAIL with undefined bootstrap initializer

- [ ] **Step 3: Write minimal implementation**

```go
email := strings.TrimSpace(strings.ToLower(cfg.AdminBootstrap.Email))
if email == "" {
	email = "admin@go-oj.dev"
}
password := strings.TrimSpace(cfg.AdminBootstrap.Password)
if password == "" {
	password = "Admin@123456"
}
displayName := strings.TrimSpace(cfg.AdminBootstrap.DisplayName)
if displayName == "" {
	displayName = "Super Admin"
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/bootstrap -run TestEnsureBootstrapAdminCreatesFallbackSuperAdmin -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/config/config.go internal/bootstrap/app.go internal/repository/admin_user.go internal/bootstrap/app_test.go
git commit -m "feat: seed bootstrap super admin"
```

### Task 2: Return A Real Admin Dashboard Payload

**Files:**
- Modify: `internal/handler/admin_stub.go`
- Test: `internal/router/admin_router_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestAdminDashboardReturnsAdminContext(t *testing.T) {
	authz, _ := service.NewInMemoryAdminAuthorizer()
	_ = authz.AssignRole("1", service.AdminRoleAdmin)

	r := newAdminTestRouter(t, authz)
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !strings.Contains(rec.Body.String(), "\"email\":\"admin@go-oj.dev\"") {
		t.Fatalf("dashboard response missing admin email: %s", rec.Body.String())
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/router -run TestAdminDashboardReturnsAdminContext -v`
Expected: FAIL because dashboard only returns `{ "ok": true }`

- [ ] **Step 3: Write minimal implementation**

```go
func (h *AdminHandler) Dashboard(c *gin.Context) {
	resp := response.Success(gin.H{
		"title": "Admin Dashboard",
		"admin_user": gin.H{
			"id":           c.GetUint("admin_user_id"),
			"email":        c.GetString("admin_user_email"),
			"display_name": c.GetString("admin_user_display_name"),
		},
	})
	resp.RequestID = c.GetString("request_id")
	c.JSON(http.StatusOK, resp)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/router -run TestAdminDashboardReturnsAdminContext -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/handler/admin_stub.go internal/router/admin_router_test.go internal/router/router.go
git commit -m "feat: add admin dashboard payload"
```

### Task 3: Add Frontend Admin Auth Store, Routes, And Pages

**Files:**
- Create: `frontend/src/features/admin-auth/store.ts`
- Create: `frontend/src/features/admin-auth/api.ts`
- Create: `frontend/src/pages/AdminLoginPage.vue`
- Create: `frontend/src/pages/AdminDashboardPage.vue`
- Modify: `frontend/src/router/index.ts`
- Modify: `frontend/src/shared/types/api.ts`
- Modify: `frontend/src/shared/lib/http.ts`
- Test: `frontend/src/router/router.test.ts`
- Test: `frontend/src/features/admin-auth/store.test.ts`

- [ ] **Step 1: Write the failing test**

```ts
it('redirects anonymous admin visitors to /admin/login', async () => {
  setActivePinia(createPinia())
  const router = createAppRouter()

  await router.push('/admin')

  expect(router.currentRoute.value.fullPath).toBe('/admin/login')
})
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd frontend && npm run test:run -- src/router/router.test.ts`
Expected: FAIL because admin routes do not exist yet

- [ ] **Step 3: Write minimal implementation**

```ts
if (to.path.startsWith('/admin') && !adminAuthStore.isAuthenticated) {
  next('/admin/login')
  return
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd frontend && npm run test:run -- src/router/router.test.ts src/features/admin-auth/store.test.ts`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add frontend/src/features/admin-auth/store.ts frontend/src/features/admin-auth/api.ts frontend/src/pages/AdminLoginPage.vue frontend/src/pages/AdminDashboardPage.vue frontend/src/router/index.ts frontend/src/shared/types/api.ts frontend/src/shared/lib/http.ts frontend/src/features/admin-auth/store.test.ts frontend/src/router/router.test.ts
git commit -m "feat: add admin frontend auth flow"
```

### Task 4: Verify End-To-End Admin Login Flow

**Files:**
- Modify: `.env.example`
- Modify: `docs/project-prd.md`

- [ ] **Step 1: Run backend tests**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 2: Run frontend tests**

Run: `cd frontend && npm run test:run`
Expected: PASS

- [ ] **Step 3: Run frontend build**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 4: Manual login smoke**

Run: `go run ./cmd/server` and `cd frontend && npm run dev`
Expected: Logging in at `/admin/login` with `admin@go-oj.dev / Admin@123456` redirects to `/admin` and renders dashboard data from the protected backend endpoint.

- [ ] **Step 5: Commit**

```bash
git add .env.example docs/project-prd.md
git commit -m "docs: document bootstrap admin login flow"
```

### Task 1: Add Admin Domain Models And Persistence

**Files:**
- Create: `internal/model/admin_user.go`
- Create: `internal/model/operation_log.go`
- Create: `internal/repository/admin_user.go`
- Modify: `internal/bootstrap/app.go`
- Test: `internal/bootstrap/app_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestInitDBFailure(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "127.0.0.1",
			Port:     1,
			Name:     "oj",
			User:     "postgres",
			Password: "Postgres",
			SSLMode:  "disable",
		},
	}

	if _, err := InitDB(cfg); err == nil {
		t.Fatal("InitDB() expected error, got nil")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/bootstrap -run TestInitDBFailure -v`
Expected: PASS now, then add table-migration assertions in a new focused test that fails until `admin_users` and `operation_logs` join the migration list.

- [ ] **Step 3: Write minimal implementation**

```go
type AdminUser struct {
	BaseModel
	Email        string `gorm:"size:128;not null;uniqueIndex" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	DisplayName  string `gorm:"size:64;not null" json:"display_name"`
	Status       string `gorm:"size:32;not null;default:active;index" json:"status"`
	LastLoginAt  int64  `gorm:"not null;default:0" json:"last_login_at"`
}

type OperationLog struct {
	BaseModel
	AdminUserID uint   `gorm:"not null;index" json:"admin_user_id"`
	Resource    string `gorm:"size:64;not null" json:"resource"`
	Action      string `gorm:"size:32;not null" json:"action"`
	TargetType  string `gorm:"size:64;not null" json:"target_type"`
	TargetID    string `gorm:"size:64;not null" json:"target_id"`
	RequestID   string `gorm:"size:64;not null" json:"request_id"`
	DetailJSON  string `gorm:"type:text;not null" json:"detail_json"`
	IP          string `gorm:"size:64;not null" json:"ip"`
	UserAgent   string `gorm:"size:255;not null" json:"user_agent"`
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/bootstrap -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/model/admin_user.go internal/model/operation_log.go internal/repository/admin_user.go internal/bootstrap/app.go internal/bootstrap/app_test.go
git commit -m "feat: add admin auth persistence models"
```

### Task 2: Add Admin Login Service And JWT Parsing

**Files:**
- Create: `internal/service/admin_auth.go`
- Create: `internal/handler/admin_auth.go`
- Modify: `internal/service/auth_test.go`
- Test: `internal/service/admin_auth_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestAdminAuthServiceLogin(t *testing.T) {
	t.Run("disabled admin cannot login", func(t *testing.T) {
		repo := &stubAdminUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*model.AdminUser, error) {
				return &model.AdminUser{
					BaseModel:    model.BaseModel{ID: 7},
					Email:        email,
					PasswordHash: mustHash(t, "password123"),
					DisplayName:  "ops",
					Status:       "disabled",
				}, nil
			},
		}

		svc := NewAdminAuthService(repo, []byte("secret"), time.Hour)
		_, err := svc.Login(context.Background(), AdminLoginInput{Email: "ops@example.com", Password: "password123"})
		if !errors.Is(err, ErrUserDisabled) {
			t.Fatalf("expected ErrUserDisabled, got %v", err)
		}
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/service -run TestAdminAuthServiceLogin -v`
Expected: FAIL with undefined admin auth types

- [ ] **Step 3: Write minimal implementation**

```go
claims := jwt.MapClaims{
	"admin_user_id": user.ID,
	"email":         user.Email,
	"scope":         "admin",
	"exp":           time.Now().Add(tokenTTL).Unix(),
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/service -run TestAdminAuthServiceLogin -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/service/admin_auth.go internal/handler/admin_auth.go internal/service/admin_auth_test.go
git commit -m "feat: add admin login service"
```

### Task 3: Add Casbin Authorization Service

**Files:**
- Create: `internal/service/admin_authorizer.go`
- Test: `internal/service/admin_authorizer_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestAdminAuthorizerAssistantCannotWriteSettings(t *testing.T) {
	authz, err := NewAdminAuthorizerWithMemoryPolicies()
	if err != nil {
		t.Fatalf("NewAdminAuthorizerWithMemoryPolicies() error = %v", err)
	}

	ok, err := authz.Enforce("2", "admin", "settings", "write")
	if err != nil {
		t.Fatalf("Enforce() error = %v", err)
	}
	if ok {
		t.Fatal("assistant should not have settings:write")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/service -run TestAdminAuthorizerAssistantCannotWriteSettings -v`
Expected: FAIL with undefined authorizer

- [ ] **Step 3: Write minimal implementation**

```go
const adminDomain = "admin"

var defaultPolicyRules = [][]string{
	{"admin", adminDomain, "dashboard", "read"},
	{"admin", adminDomain, "dashboard", "write"},
	{"assistant", adminDomain, "problems", "read"},
	{"assistant", adminDomain, "problems", "write"},
	{"assistant", adminDomain, "submissions", "read"},
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/service -run TestAdminAuthorizer -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/service/admin_authorizer.go internal/service/admin_authorizer_test.go go.mod go.sum
git commit -m "feat: add admin casbin authorizer"
```

### Task 4: Protect Admin Routes With Middleware And Permission Mapping

**Files:**
- Create: `internal/handler/admin_stub.go`
- Modify: `internal/router/router.go`
- Modify: `internal/router/router_test.go`
- Modify: `internal/router/router_api_test.go`
- Test: `internal/router/admin_router_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestAdminRoutesRequireAuthorization(t *testing.T) {
	r := newAdminTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/admin/settings", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("GET /admin/settings status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/router -run TestAdminRoutesRequireAuthorization -v`
Expected: FAIL with 404 until admin routes and middleware exist

- [ ] **Step 3: Write minimal implementation**

```go
admin := r.Group("/admin")
admin.POST("/login", adminAuthHandler.Login)
admin.Use(adminAuthMiddleware(...))
admin.GET("/settings", adminStubHandler.Settings)
admin.PATCH("/settings", adminStubHandler.Settings)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/router -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/handler/admin_stub.go internal/router/router.go internal/router/router_test.go internal/router/router_api_test.go internal/router/admin_router_test.go
git commit -m "feat: protect admin routes with casbin authz"
```

### Task 5: Wire App Bootstrap And Seed Default Policies

**Files:**
- Modify: `internal/bootstrap/app.go`
- Modify: `internal/bootstrap/app_test.go`
- Test: `internal/bootstrap/app_test.go`

- [ ] **Step 1: Write the failing test**

```go
func TestBuildAdminAuthorizerRejectsInvalidConfig(t *testing.T) {
	_, err := buildAdminAuthorizer(nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/bootstrap -run TestBuildAdminAuthorizerRejectsInvalidConfig -v`
Expected: FAIL with undefined helper

- [ ] **Step 3: Write minimal implementation**

```go
authorizer, err := service.NewAdminAuthorizer(db)
if err != nil {
	return nil, fmt.Errorf("new admin authorizer: %w", err)
}
if err := authorizer.SeedPolicies(context.Background()); err != nil {
	return nil, fmt.Errorf("seed admin policies: %w", err)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/bootstrap -v && go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/bootstrap/app.go internal/bootstrap/app_test.go
git commit -m "feat: wire admin authz into app bootstrap"
```
