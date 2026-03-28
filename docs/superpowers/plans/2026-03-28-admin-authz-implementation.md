# Admin Authz Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first usable backend foundation for `/admin/*` authentication and authorization with separate `admin_users`, Casbin-backed RBAC with domain, and protected admin routes.

**Architecture:** Keep the existing frontend-user auth flow unchanged, then add a parallel admin identity flow with its own repository, service, JWT claims, Gin middleware, and Casbin authorization service. Register a minimal set of admin routes now so routing, token validation, status checks, and permission checks can be tested end-to-end before real admin business handlers are added.

**Tech Stack:** Go, Gin, Gorm, PostgreSQL, Casbin

---

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
