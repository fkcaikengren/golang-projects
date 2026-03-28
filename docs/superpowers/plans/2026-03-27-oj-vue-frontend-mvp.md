# OJ Vue Frontend MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a standalone Vite + Vue 3 frontend under `frontend/` that provides the OJ MVP user-facing shell with auth, problem-set browsing, problem browsing, problem detail, and submissions pages backed by mock-friendly API/query layers.

**Architecture:** Create a dedicated `frontend/` app using Vue Router for route structure, Pinia for auth state, TanStack Vue Query for server-state access, and Axios for API modules that mirror the existing Gin routes under `/api/v1`. Keep route pages thin, move reusable query and API code into feature modules, and allow mock data to satisfy queries before full backend integration.

**Tech Stack:** Vite, Vue 3, TypeScript, Vue Router, Pinia, Naive UI, @tanstack/vue-query, Axios, Vitest, Vue Test Utils

---

## File Structure

**Create**

- `frontend/package.json`
- `frontend/tsconfig.json`
- `frontend/tsconfig.app.json`
- `frontend/tsconfig.node.json`
- `frontend/vite.config.ts`
- `frontend/index.html`
- `frontend/src/main.ts`
- `frontend/src/App.vue`
- `frontend/src/app/providers.ts`
- `frontend/src/app/styles.css`
- `frontend/src/router/index.ts`
- `frontend/src/shared/lib/http.ts`
- `frontend/src/shared/config/env.ts`
- `frontend/src/shared/types/api.ts`
- `frontend/src/shared/mocks/data.ts`
- `frontend/src/shared/mocks/server.ts`
- `frontend/src/shared/ui/AppShell.vue`
- `frontend/src/features/auth/api.ts`
- `frontend/src/features/auth/store.ts`
- `frontend/src/features/problem-sets/api.ts`
- `frontend/src/features/problems/api.ts`
- `frontend/src/features/submissions/api.ts`
- `frontend/src/features/problem-sets/queries.ts`
- `frontend/src/features/problems/queries.ts`
- `frontend/src/features/submissions/queries.ts`
- `frontend/src/pages/HomePage.vue`
- `frontend/src/pages/AuthPage.vue`
- `frontend/src/pages/ProblemSetsPage.vue`
- `frontend/src/pages/ProblemSetDetailPage.vue`
- `frontend/src/pages/ProblemsPage.vue`
- `frontend/src/pages/ProblemDetailPage.vue`
- `frontend/src/pages/SubmissionsPage.vue`
- `frontend/src/pages/ProgressPage.vue`
- `frontend/src/widgets/home/HeroPanel.vue`
- `frontend/src/widgets/home/FeaturedSets.vue`
- `frontend/src/widgets/home/RecentProblems.vue`
- `frontend/src/components.d.ts`
- `frontend/src/vite-env.d.ts`
- `frontend/src/test/setup.ts`
- `frontend/src/features/auth/store.test.ts`
- `frontend/src/router/router.test.ts`
- `frontend/src/pages/HomePage.test.ts`

**Modify**

- `.gitignore`

## Task 1: Scaffold the frontend workspace

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/tsconfig.json`
- Create: `frontend/tsconfig.app.json`
- Create: `frontend/tsconfig.node.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/index.html`
- Create: `frontend/src/main.ts`
- Create: `frontend/src/App.vue`
- Create: `frontend/src/vite-env.d.ts`
- Modify: `.gitignore`

- [ ] **Step 1: Write the failing package/build expectation**

Document the initial contract in `frontend/package.json`: scripts `dev`, `build`, `preview`, `test`, `test:run`. Dependencies must include `vue`, `vue-router`, `pinia`, `naive-ui`, `@tanstack/vue-query`, `axios`. Dev dependencies must include `vite`, `typescript`, `vitest`, `@vitejs/plugin-vue`, `@vue/test-utils`, `jsdom`.

- [ ] **Step 2: Create the Vite Vue TypeScript app files**

Use Vite scaffolding semantics but commit explicit files instead of generated boilerplate. `frontend/src/main.ts` should mount `App.vue` and install shared providers from `src/app/providers.ts`.

- [ ] **Step 3: Add app shell entry and root HTML**

`frontend/src/App.vue` should render `RouterView`. `frontend/index.html` should contain a root `#app` node and title for the OJ frontend.

- [ ] **Step 4: Ignore frontend build artifacts**

Add `frontend/node_modules`, `frontend/dist`, and `frontend/.vitest` to `.gitignore` if they are not already ignored.

- [ ] **Step 5: Verify the scaffold builds**

Run: `cd frontend && npm install && npm run build`
Expected: Vite build succeeds with no missing entry file errors.

## Task 2: Establish providers, routing, and base layout

**Files:**
- Create: `frontend/src/app/providers.ts`
- Create: `frontend/src/app/styles.css`
- Create: `frontend/src/router/index.ts`
- Create: `frontend/src/shared/ui/AppShell.vue`
- Create: `frontend/src/pages/HomePage.vue`
- Create: `frontend/src/pages/AuthPage.vue`
- Create: `frontend/src/pages/ProblemSetsPage.vue`
- Create: `frontend/src/pages/ProblemSetDetailPage.vue`
- Create: `frontend/src/pages/ProblemsPage.vue`
- Create: `frontend/src/pages/ProblemDetailPage.vue`
- Create: `frontend/src/pages/SubmissionsPage.vue`
- Create: `frontend/src/pages/ProgressPage.vue`

- [ ] **Step 1: Write the failing router test**

Create `frontend/src/router/router.test.ts` asserting the router resolves these paths: `/`, `/login`, `/register`, `/problem-sets`, `/problem-sets/:slug`, `/problems`, `/problems/:slug`, `/submissions`, `/progress`.

- [ ] **Step 2: Run the router test to verify it fails**

Run: `cd frontend && npm run test:run -- src/router/router.test.ts`
Expected: FAIL because router module or routes do not exist yet.

- [ ] **Step 3: Implement the router and page stubs**

Create one page component per route and register them in `src/router/index.ts`. Keep page components minimal but routeable.

- [ ] **Step 4: Implement shared providers and app shell**

`src/app/providers.ts` should export an installer for Pinia, Vue Query, Router, and Naive UI config. `AppShell.vue` should provide the top navigation and content slot used by route pages.

- [ ] **Step 5: Run the router test and build verification**

Run: `cd frontend && npm run test:run -- src/router/router.test.ts && npm run build`
Expected: Router test passes and the app still builds.

## Task 3: Add typed config, HTTP client, and mock data source

**Files:**
- Create: `frontend/src/shared/config/env.ts`
- Create: `frontend/src/shared/lib/http.ts`
- Create: `frontend/src/shared/types/api.ts`
- Create: `frontend/src/shared/mocks/data.ts`
- Create: `frontend/src/shared/mocks/server.ts`

- [ ] **Step 1: Write the failing Home page render test**

Create `frontend/src/pages/HomePage.test.ts` that mounts `HomePage` with app providers and expects the hero title plus at least one featured problem set card sourced through the mock-backed query layer.

- [ ] **Step 2: Run the Home page test to verify it fails**

Run: `cd frontend && npm run test:run -- src/pages/HomePage.test.ts`
Expected: FAIL because API/query mocks are not wired.

- [ ] **Step 3: Define API response types and mock fixtures**

Create stable front-end types for auth user, problem set summary, problem summary, problem detail, and submission item. Populate mock fixtures that match the backend MVP fields and route semantics.

- [ ] **Step 4: Implement env and HTTP client with mock fallback**

`env.ts` should expose `apiBaseUrl` and `useMock`. `http.ts` should create a shared Axios instance. `server.ts` should expose async functions that return mock payloads shaped like API responses when `useMock` is true.

- [ ] **Step 5: Re-run the Home page test**

Run: `cd frontend && npm run test:run -- src/pages/HomePage.test.ts`
Expected: FAIL moves forward only if the page itself is still unimplemented; mock data layer should no longer be the blocker.

## Task 4: Implement auth state and API modules

**Files:**
- Create: `frontend/src/features/auth/api.ts`
- Create: `frontend/src/features/auth/store.ts`
- Create: `frontend/src/features/auth/store.test.ts`
- Create: `frontend/src/features/problem-sets/api.ts`
- Create: `frontend/src/features/problems/api.ts`
- Create: `frontend/src/features/submissions/api.ts`

- [ ] **Step 1: Write the failing auth store test**

`frontend/src/features/auth/store.test.ts` should verify three behaviors: initial anonymous state, login writes token/user into the store, logout clears them.

- [ ] **Step 2: Run the auth store test to verify it fails**

Run: `cd frontend && npm run test:run -- src/features/auth/store.test.ts`
Expected: FAIL because the auth store is missing.

- [ ] **Step 3: Implement auth API and store**

`auth/api.ts` should expose `login` and `register` calling either Axios or mock server functions. `store.ts` should define a Pinia store with `setSession` and `clearSession`.

- [ ] **Step 4: Implement the remaining domain API modules**

`problem-sets/api.ts`, `problems/api.ts`, and `submissions/api.ts` should wrap the corresponding backend endpoints:

```text
GET /api/v1/problem-sets
GET /api/v1/problem-sets/:slug
GET /api/v1/problems
GET /api/v1/problems/:slug
GET /api/v1/submissions
POST /api/v1/submissions
```

- [ ] **Step 5: Re-run the auth store test**

Run: `cd frontend && npm run test:run -- src/features/auth/store.test.ts`
Expected: PASS

## Task 5: Implement query modules and user-facing pages

**Files:**
- Create: `frontend/src/features/problem-sets/queries.ts`
- Create: `frontend/src/features/problems/queries.ts`
- Create: `frontend/src/features/submissions/queries.ts`
- Create: `frontend/src/widgets/home/HeroPanel.vue`
- Create: `frontend/src/widgets/home/FeaturedSets.vue`
- Create: `frontend/src/widgets/home/RecentProblems.vue`
- Modify: `frontend/src/pages/HomePage.vue`
- Modify: `frontend/src/pages/AuthPage.vue`
- Modify: `frontend/src/pages/ProblemSetsPage.vue`
- Modify: `frontend/src/pages/ProblemSetDetailPage.vue`
- Modify: `frontend/src/pages/ProblemsPage.vue`
- Modify: `frontend/src/pages/ProblemDetailPage.vue`
- Modify: `frontend/src/pages/SubmissionsPage.vue`
- Modify: `frontend/src/pages/ProgressPage.vue`

- [ ] **Step 1: Write the failing page expectations**

Extend `frontend/src/pages/HomePage.test.ts` or add focused tests so the route pages expose the MVP landmarks:

```ts
expect(screen.getByText('开始刷题')).toBeTruthy()
expect(screen.getByText('热门题单')).toBeTruthy()
```

- [ ] **Step 2: Run the page test to verify it fails**

Run: `cd frontend && npm run test:run -- src/pages/HomePage.test.ts`
Expected: FAIL because the widgets and query hooks are not implemented yet.

- [ ] **Step 3: Implement query composables**

Each feature query module should wrap the matching API functions with `useQuery` or `useMutation`, returning typed data for route pages.

- [ ] **Step 4: Implement page UIs**

Build the route pages with Naive UI primitives and a shared OJ-styled shell:

- `HomePage`: hero, featured sets, recent problems
- `AuthPage`: login/register tabs and submit action
- `ProblemSetsPage`: set cards
- `ProblemSetDetailPage`: set summary plus problem list
- `ProblemsPage`: list/table with difficulty/tag/status
- `ProblemDetailPage`: statement plus code input area and submission action
- `SubmissionsPage`: user submission history
- `ProgressPage`: placeholder for future completion stats

- [ ] **Step 5: Re-run page tests and the build**

Run: `cd frontend && npm run test:run -- src/pages/HomePage.test.ts && npm run build`
Expected: PASS and build succeeds.

## Task 6: Full verification

**Files:**
- Verify only

- [ ] **Step 1: Run the full test suite**

Run: `cd frontend && npm run test:run`
Expected: All frontend tests pass.

- [ ] **Step 2: Run production build**

Run: `cd frontend && npm run build`
Expected: Successful production build output in `frontend/dist`.

- [ ] **Step 3: Run a final manual smoke start**

Run: `cd frontend && npm run dev -- --host 0.0.0.0`
Expected: Dev server starts and serves the app for manual inspection.

- [ ] **Step 4: Commit**

```bash
git add .gitignore frontend
git commit -m "feat: add oj vue frontend mvp"
```
