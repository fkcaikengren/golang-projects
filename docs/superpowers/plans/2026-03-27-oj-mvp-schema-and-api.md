# OJ MVP Schema And API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first usable OJ MVP backend covering user email/password auth, problem sets, problems with tags, submission history, and per-user problem status on top of the existing Gin + GORM service.

**Architecture:** Extend the current layered structure (`model` -> `repository` -> `service` -> `handler` -> `router`) instead of introducing a new framework. Keep MVP state in Postgres via GORM auto-migration, expose read APIs for problem browsing, and write APIs for registration/login and submissions; submission judging remains synchronous placeholder logic for now, while `user_problem_stats` is updated transactionally.

**Tech Stack:** Go 1.26, Gin, GORM, PostgreSQL, existing `internal/pkg/response`, standard library crypto packages, `golang.org/x/crypto/bcrypt`, `github.com/golang-jwt/jwt/v5` if token auth is added.

---

## File Structure

**Create**

- `internal/model/user.go`
- `internal/model/problem_set.go`
- `internal/model/problem.go`
- `internal/model/tag.go`
- `internal/model/problem_tag.go`
- `internal/model/problem_set_problem.go`
- `internal/model/submission.go`
- `internal/model/user_problem_stat.go`
- `internal/repository/user.go`
- `internal/repository/problem_set.go`
- `internal/repository/problem.go`
- `internal/repository/submission.go`
- `internal/service/auth.go`
- `internal/service/problem_set.go`
- `internal/service/problem.go`
- `internal/service/submission.go`
- `internal/handler/auth.go`
- `internal/handler/problem_set.go`
- `internal/handler/problem.go`
- `internal/handler/submission.go`
- `internal/service/auth_test.go`
- `internal/service/submission_test.go`
- `internal/router/router_api_test.go`

**Modify**

- `go.mod`
- `internal/bootstrap/app.go`
- `internal/router/router.go`
- `internal/config/config.go`
- `internal/config/config_test.go`

## Task 1: Add MVP Data Models And Auto-Migration

**Files:**
- Create: `internal/model/user.go`
- Create: `internal/model/problem_set.go`
- Create: `internal/model/problem.go`
- Create: `internal/model/tag.go`
- Create: `internal/model/problem_tag.go`
- Create: `internal/model/problem_set_problem.go`
- Create: `internal/model/submission.go`
- Create: `internal/model/user_problem_stat.go`
- Modify: `internal/bootstrap/app.go`

- [ ] **Step 1: Write the failing bootstrap migration test expectation**

Document the target migration set before code changes. The app should migrate these models in `InitDB`: `SystemSetting`, `User`, `ProblemSet`, `Problem`, `Tag`, `ProblemTag`, `ProblemSetProblem`, `Submission`, `UserProblemStat`.

- [ ] **Step 2: Add the model files using `BaseModel` conventions**

Create structs with the exact fields from the approved spec. Example for `internal/model/user.go`:

```go
package model

type User struct {
	BaseModel
	Email        string `gorm:"size:128;not null;uniqueIndex" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Nickname     string `gorm:"size:64;not null" json:"nickname"`
	Status       string `gorm:"size:32;not null;default:active;index" json:"status"`
	LastLoginAt  int64  `gorm:"not null;default:0" json:"last_login_at"`
}
```

Follow the same pattern for the other entities with the approved indexes.

- [ ] **Step 3: Update auto-migration**

Modify `internal/bootstrap/app.go`:

```go
	if err := db.AutoMigrate(
		&model.SystemSetting{},
		&model.User{},
		&model.ProblemSet{},
		&model.Problem{},
		&model.Tag{},
		&model.ProblemTag{},
		&model.ProblemSetProblem{},
		&model.Submission{},
		&model.UserProblemStat{},
	); err != nil {
		return nil, err
	}
```

- [ ] **Step 4: Run the model/build verification**

Run: `go test ./internal/bootstrap ./internal/model/...`

Expected: packages compile; no undefined model references.

- [ ] **Step 5: Commit**

```bash
git add internal/model internal/bootstrap/app.go
git commit -m "feat: add oj mvp data models"
```

## Task 2: Add Auth Configuration And User Repository

**Files:**
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`
- Create: `internal/repository/user.go`

- [ ] **Step 1: Write the failing config test**

Add or extend a config test that asserts auth settings can be loaded, for example:

```go
func TestLoadConfigIncludesJWTSecret(t *testing.T) {
	cfg, err := Load()
	require.NoError(t, err)
	require.NotEmpty(t, cfg.Auth.JWTSecret)
}
```

- [ ] **Step 2: Extend config with auth settings**

Add an auth section:

```go
type AuthConfig struct {
	JWTSecret string
	TokenTTL  int
}

type Config struct {
	AppName  string
	HTTPPort string
	Database DatabaseConfig
	Auth     AuthConfig
}
```

Bind from env/defaults in the existing load path.

- [ ] **Step 3: Add user repository methods**

Create `internal/repository/user.go` with focused queries:

```go
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository { ... }

func (r *UserRepository) Create(ctx context.Context, user *model.User) error { ... }
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) { ... }
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) { ... }
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uint, ts int64) error { ... }
```

- [ ] **Step 4: Run tests**

Run: `go test ./internal/config ./internal/repository/...`

Expected: config tests pass and repository package compiles.

- [ ] **Step 5: Commit**

```bash
git add internal/config internal/repository/user.go
git commit -m "feat: add auth config and user repository"
```

## Task 3: Implement Registration And Login Service

**Files:**
- Modify: `go.mod`
- Create: `internal/service/auth.go`
- Create: `internal/service/auth_test.go`

- [ ] **Step 1: Write the failing auth service tests**

Create tests for the minimal service contract:

```go
func TestAuthServiceRegister(t *testing.T) {
	// verifies password is hashed and duplicate email is rejected
}

func TestAuthServiceLogin(t *testing.T) {
	// verifies correct password returns token payload and wrong password fails
}
```

Use a small fake/stub user repository in the test file instead of a real DB.

- [ ] **Step 2: Add dependencies**

Run: `go get golang.org/x/crypto/bcrypt github.com/golang-jwt/jwt/v5`

Expected: `go.mod` and `go.sum` updated.

- [ ] **Step 3: Implement auth service**

Create `internal/service/auth.go` with:

```go
type AuthService struct {
	users     userRepo
	jwtSecret []byte
	tokenTTL  time.Duration
}

type RegisterInput struct {
	Email    string
	Password string
	Nickname string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	Token string `json:"token"`
	User  *model.User `json:"user"`
}
```

Methods:

- `Register(ctx, input)` validates email/password, checks duplicate email, hashes password, stores the user
- `Login(ctx, input)` loads by email, compares bcrypt hash, updates `last_login_at`, returns JWT

- [ ] **Step 4: Run tests**

Run: `go test ./internal/service -run 'TestAuthService' -v`

Expected: PASS for register/login cases.

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum internal/service/auth.go internal/service/auth_test.go
git commit -m "feat: implement auth service"
```

## Task 4: Expose Auth Handlers And Routes

**Files:**
- Create: `internal/handler/auth.go`
- Modify: `internal/router/router.go`
- Modify: `internal/bootstrap/app.go`

- [ ] **Step 1: Write the failing router test**

Add a router test for:

```go
func TestRouterRegistersAuthRoutes(t *testing.T) {
	// POST /api/v1/auth/register
	// POST /api/v1/auth/login
}
```

Expected route status is not `404`.

- [ ] **Step 2: Implement auth handler**

Create `internal/handler/auth.go`:

```go
type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler { ... }
func (h *AuthHandler) Register(c *gin.Context) { ... }
func (h *AuthHandler) Login(c *gin.Context) { ... }
```

Request payloads:

```go
type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required,min=2,max=32"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
```

- [ ] **Step 3: Wire app and router**

Update `internal/bootstrap/app.go` to construct `UserRepository`, `AuthService`, and `AuthHandler`.

Update `internal/router/router.go` signature to accept `authHandler` and register:

```go
authGroup := apiV1.Group("/auth")
{
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./internal/router ./internal/handler/... -v`

Expected: auth routes exist; handler package compiles.

- [ ] **Step 5: Commit**

```bash
git add internal/handler/auth.go internal/router/router.go internal/bootstrap/app.go
git commit -m "feat: add auth endpoints"
```

## Task 5: Implement Problem Set And Problem Read Repositories

**Files:**
- Create: `internal/repository/problem_set.go`
- Create: `internal/repository/problem.go`

- [ ] **Step 1: Define the required query methods**

Problem set repository must support:

- `ListActive(ctx context.Context) ([]model.ProblemSet, error)`
- `GetBySlug(ctx context.Context, slug string) (*model.ProblemSet, error)`
- `ListProblems(ctx context.Context, setID uint) ([]model.Problem, error)`

Problem repository must support:

- `List(ctx context.Context, filter ProblemFilter) ([]model.Problem, error)`
- `GetBySlug(ctx context.Context, slug string) (*model.Problem, error)`
- `ListTags(ctx context.Context, problemID uint) ([]model.Tag, error)`

- [ ] **Step 2: Implement repository code**

Use explicit joins for the list queries. Example shape:

```go
type ProblemFilter struct {
	Difficulty string
	TagSlug    string
	Keyword    string
}
```

Filter only on non-empty fields; keep sorting simple for MVP.

- [ ] **Step 3: Run repository compilation**

Run: `go test ./internal/repository/...`

Expected: compile success with no missing imports or types.

- [ ] **Step 4: Commit**

```bash
git add internal/repository/problem_set.go internal/repository/problem.go
git commit -m "feat: add problem query repositories"
```

## Task 6: Implement Problem Browsing Services And Handlers

**Files:**
- Create: `internal/service/problem_set.go`
- Create: `internal/service/problem.go`
- Create: `internal/handler/problem_set.go`
- Create: `internal/handler/problem.go`
- Modify: `internal/bootstrap/app.go`
- Modify: `internal/router/router.go`

- [ ] **Step 1: Write the failing router API test**

Add test cases for:

- `GET /api/v1/problem-sets`
- `GET /api/v1/problem-sets/:slug`
- `GET /api/v1/problems`
- `GET /api/v1/problems/:slug`

Expected: route exists and returns JSON envelope.

- [ ] **Step 2: Implement services**

`ProblemSetService` should return active problem sets and problem set detail with embedded problem list.

`ProblemService` should return filtered problem lists and problem detail with tags.

- [ ] **Step 3: Implement handlers**

Add query bindings for filters:

```go
type problemListQuery struct {
	Difficulty string `form:"difficulty"`
	Tag        string `form:"tag"`
	Keyword    string `form:"keyword"`
}
```

Use existing response helper consistently.

- [ ] **Step 4: Wire app and router**

Add constructors in `internal/bootstrap/app.go` and routes:

```go
apiV1.GET("/problem-sets", problemSetHandler.List)
apiV1.GET("/problem-sets/:slug", problemSetHandler.Detail)
apiV1.GET("/problems", problemHandler.List)
apiV1.GET("/problems/:slug", problemHandler.Detail)
```

- [ ] **Step 5: Run tests**

Run: `go test ./internal/router ./internal/handler/... ./internal/service/... -v`

Expected: PASS for route registration and package tests.

- [ ] **Step 6: Commit**

```bash
git add internal/service/problem_set.go internal/service/problem.go internal/handler/problem_set.go internal/handler/problem.go internal/bootstrap/app.go internal/router/router.go
git commit -m "feat: add problem browsing endpoints"
```

## Task 7: Implement Submission Repository And Stats Upsert

**Files:**
- Create: `internal/repository/submission.go`

- [ ] **Step 1: Write the failing service-oriented repository expectation**

Document the repository contract:

- create a submission
- list submissions by user
- list submissions by problem and user
- update or create `user_problem_stats` in one transaction

- [ ] **Step 2: Implement repository methods**

Create:

```go
type SubmissionRepository struct {
	db *gorm.DB
}

func (r *SubmissionRepository) Create(ctx context.Context, submission *model.Submission) error { ... }
func (r *SubmissionRepository) ListByUser(ctx context.Context, userID uint) ([]model.Submission, error) { ... }
func (r *SubmissionRepository) ListByUserAndProblem(ctx context.Context, userID, problemID uint) ([]model.Submission, error) { ... }
func (r *SubmissionRepository) SaveWithStats(ctx context.Context, submission *model.Submission) error { ... }
```

`SaveWithStats` should:

- insert the submission
- fetch or initialize `user_problem_stats`
- increment `submit_count`
- if status is `accepted`, increment `accepted_count` and set `first_accepted_at` when empty
- set `last_submit_at`, `last_submission_id`
- derive aggregate status: `solved` if accepted count > 0, else `attempted`

- [ ] **Step 3: Run repository compilation**

Run: `go test ./internal/repository/...`

Expected: compile success.

- [ ] **Step 4: Commit**

```bash
git add internal/repository/submission.go
git commit -m "feat: add submission repository and stats upsert"
```

## Task 8: Implement Submission Service With Synchronous MVP Judge

**Files:**
- Create: `internal/service/submission.go`
- Create: `internal/service/submission_test.go`

- [ ] **Step 1: Write the failing submission service tests**

Add tests for:

```go
func TestSubmissionServiceSubmitAccepted(t *testing.T) {}
func TestSubmissionServiceSubmitWrongAnswer(t *testing.T) {}
func TestSubmissionServiceListUserSubmissions(t *testing.T) {}
```

Use stub repositories. The MVP judge can be deterministic placeholder logic so tests stay unit-level.

- [ ] **Step 2: Implement the service**

Create:

```go
type SubmissionService struct {
	submissions submissionRepo
	problems    problemLookupRepo
}

type SubmitCodeInput struct {
	UserID    uint
	ProblemID uint
	Language  string
	Code      string
}
```

MVP judge rule:

- empty code -> validation error
- if code contains the substring `"return"` then mark `accepted`
- otherwise mark `wrong_answer`

Persist through `SaveWithStats`.

- [ ] **Step 3: Run tests**

Run: `go test ./internal/service -run 'TestSubmissionService' -v`

Expected: PASS for accepted, wrong answer, and list cases.

- [ ] **Step 4: Commit**

```bash
git add internal/service/submission.go internal/service/submission_test.go
git commit -m "feat: add mvp submission service"
```

## Task 9: Expose Submission APIs

**Files:**
- Create: `internal/handler/submission.go`
- Modify: `internal/bootstrap/app.go`
- Modify: `internal/router/router.go`
- Create: `internal/router/router_api_test.go`

- [ ] **Step 1: Write the failing router tests**

Add cases for:

- `POST /api/v1/submissions`
- `GET /api/v1/submissions`
- `GET /api/v1/problems/:slug/submissions`

Expected: route exists and request validation runs.

- [ ] **Step 2: Implement the handler**

Create:

```go
type SubmissionHandler struct {
	service *service.SubmissionService
}

type submitRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	ProblemID uint  `json:"problem_id" binding:"required"`
	Language string `json:"language" binding:"required"`
	Code     string `json:"code" binding:"required"`
}
```

Handlers:

- `Submit`
- `ListMySubmissions`
- `ListProblemSubmissions`

For MVP, accept `user_id` via request/query instead of auth middleware.

- [ ] **Step 3: Wire routes**

Update router:

```go
apiV1.POST("/submissions", submissionHandler.Submit)
apiV1.GET("/submissions", submissionHandler.ListMySubmissions)
apiV1.GET("/problems/:slug/submissions", submissionHandler.ListProblemSubmissions)
```

- [ ] **Step 4: Run tests**

Run: `go test ./internal/router ./internal/handler/... -v`

Expected: PASS for submission route registration tests.

- [ ] **Step 5: Commit**

```bash
git add internal/handler/submission.go internal/bootstrap/app.go internal/router/router.go internal/router/router_api_test.go
git commit -m "feat: add submission endpoints"
```

## Task 10: Final Verification And Seed Data Notes

**Files:**
- Modify: `docs/superpowers/specs/2026-03-27-oj-mvp-design.md` if implementation decisions changed

- [ ] **Step 1: Run the full test suite**

Run: `go test ./...`

Expected: all tests pass.

- [ ] **Step 2: Run the server locally**

Run: `go run ./cmd/server`

Expected: server starts, DB connects, routes are available.

- [ ] **Step 3: Smoke test the APIs**

Run:

```bash
curl -X POST http://127.0.0.1:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"secret123","nickname":"alice"}'
```

Expected: `200` response with created user and token or success payload.

Run:

```bash
curl 'http://127.0.0.1:8080/api/v1/problem-sets'
```

Expected: `200` JSON envelope, possibly empty list before seeding.

- [ ] **Step 4: Commit**

```bash
git add .
git commit -m "test: verify oj mvp backend flow"
```

## Self-Review

Spec coverage check:

- 用户邮箱注册、登录: covered by Tasks 2-4
- 题库: covered by Tasks 5-6
- 题目标签: covered by Tasks 1, 5, 6
- 题目: covered by Tasks 1, 5, 6
- 用户做题记录: covered by Tasks 1, 7, 8, 9

Placeholder scan:

- No `TODO`, `TBD`, or “similar to previous task” placeholders remain.
- Each code-writing task names exact files and concrete method shapes.

Type consistency check:

- Model names match the approved spec and migration list.
- Route names are consistent across bootstrap, handler, and router tasks.
- `user_problem_stats` remains aggregate-only; submission detail stays in `submissions`.
