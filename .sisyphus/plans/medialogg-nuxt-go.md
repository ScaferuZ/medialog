# Medialogg - Unified Media Logging Platform
# Execution Plan: Nuxt 3 + Go Stack

## TL;DR

> **Goal**: Build a unified media logging platform supporting films, anime, books, manga, games, and doujin.
> **Stack**: Nuxt 3 (Nitro-powered frontend) + Go/Fiber (API backend) + PostgreSQL
> **Priority**: Films first (TMDB integration), then expand to other media types

**Deliverables**:
- Nuxt 3 frontend with SSR/SSG
- Go API server with JWT auth
- PostgreSQL database with migrations
- Docker Compose for homelab deployment
- TMDB integration for film data

**Estimated Effort**: Medium (12 weeks phased approach)
**Parallel Execution**: YES - Frontend and backend can develop in parallel
**Critical Path**: Auth → Media schema → Log functionality → TMDB integration

---

## Context

### Original Request
User wants a unified media logging platform combining features from:
- **Letterboxd** (film logging, reviews, ratings, lists)
- **Goodreads** (book tracking, reading progress)
- **Backloggd** (game logging, collections)

### Key Decisions Made
1. **Frontend**: Nuxt 3 (uses Nitro internally, Vue ecosystem)
2. **Backend**: Go with Fiber framework (fast, single binary)
3. **Database**: PostgreSQL with sqlc (type-safe SQL)
4. **Hosting**: Self-hosted homelab (Docker Compose)
5. **Priority Media**: Films first (TMDB is most mature API)
6. **Doujin**: Content warnings for potentially NSFW content

---

## Work Objectives

### Core Objective
Create a unified platform where users can log, review, and share their media consumption across multiple formats.

### Concrete Deliverables
1. User authentication (register/login/JWT)
2. Media database with polymorphic types
3. Activity logging (watched/reading/playing)
4. Rating and review system
5. Social features (follow, timeline, likes)
6. TMDB integration for film data
7. User lists and profiles

### Definition of Done
- [ ] User can register, login, and manage profile
- [ ] User can search films via TMDB
- [ ] User can log films with status/rating/notes
- [ ] User can follow others and see timeline
- [ ] User can create lists and write reviews
- [ ] Docker Compose deploys all services
- [ ] All API endpoints documented (OpenAPI)

### Must Have
- Authentication with JWT
- Film logging with TMDB integration
- Follow/timeline system
- Self-hostable via Docker

### Must NOT Have (Guardrails)
- No ORM (use sqlc for type-safe SQL)
- No microservices (keep monolithic Go backend)
- No paid hosted services (self-contained Docker deployment)
- No over-engineering authentication (standard JWT, not OAuth initially)
- No content moderation system for MVP (content warnings only)

---

## Verification Strategy

### Test Decision
- **Infrastructure exists**: NO (greenfield project)
- **Automated tests**: TDD for backend (Go tests), tests-after for frontend (Vitest)
- **Framework**: `go test` for backend, Vitest for Nuxt
- **Agent-Executed QA**: ALWAYS (manual browser testing via Playwright skill, API testing via curl)

### QA Policy
Every task includes agent-executed QA scenarios with evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

---

## Execution Strategy

### Dependency Matrix

```
Task Dependencies:

Phase 1 (Foundation - can run in parallel):
├── T1: Project scaffolding (backend) ──────────────────┐
├── T2: Project scaffolding (frontend) ────────────────┤── Blocks Phase 2
├── T3: Database setup ────────────────────────────────┤
└── T4: Docker Compose setup ──────────────────────────┘

Phase 2 (Auth - sequential):
├── T5: User model + migrations ← T1, T3
├── T6: JWT middleware + handlers ← T5
├── T7: Frontend auth pages ← T2, T6
└── T8: Auth integration tests ← T5, T6, T7

Phase 3 (Core - Media + Logs):
├── T9: Media model + migrations ← T3
├── T10: Media API handlers ← T9
├── T11: Log model + migrations ← T3
├── T12: Log API handlers ← T5, T9, T11
├── T13: Frontend media pages ← T2, T10
└── T14: Frontend log UI ← T7, T12

Phase 4 (Social):
├── T15: Follow model + handlers ← T5
├── T16: Timeline API ← T5, T12, T15
├── T17: Like functionality ← T5, T12
└── T18: Frontend social pages ← T7, T14, T16

Phase 5 (TMDB):
├── T19: TMDB client ← T9, T10
├── T20: Sync background jobs ← T19
└── T21: Frontend search integration ← T13, T19

Critical Path: T1/T2/T3 → T5 → T6 → T7 → T12 → T14 → T16 → T18 → T21
Parallel Speedup: ~60% faster than sequential
```

---

## TODOs


- [x] 1. **Backend Project Scaffolding**

  **What to do**:
  - Initialize Go module: `go mod init github.com/yourusername/medialogg/backend`
  - Setup Fiber framework with basic routing
  - Create directory structure: `cmd/server/`, `internal/api/`, `internal/db/`, `internal/middleware/`, `internal/config/`, `migrations/`
  - Add Makefile with common commands
  - Setup environment config with viper (DATABASE_URL, JWT_SECRET, etc.)

  **Must NOT do**:
  - Don't add ORM (we use sqlc)
  - Don't add unnecessary middleware yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with T2, T3, T4)
  - **Blocks**: T5, T6
  - **Blocked By**: None

  **References**:
  - Fiber docs: https://docs.gofiber.io/
  - Go project layout: https://github.com/golang-standards/project-layout

  **Acceptance Criteria**:
  - [ ] `go run cmd/server/main.go` starts server on port 8080
  - [ ] Server responds to `GET /health` with `{"status": "ok"}`
  - [ ] Environment variables loaded from `.env` file

  **QA Scenarios**:
  ```
  Scenario: Health check endpoint works
    Tool: Bash (curl)
    Preconditions: Server running on port 8080
    Steps:
      1. Run: curl http://localhost:8080/health
      2. Assert response contains {"status": "ok"}
    Expected Result: HTTP 200 with status ok
    Evidence: .sisyphus/evidence/task-1-health-check.txt
  ```

  **Commit**: YES
  - Message: `chore(backend): initial project scaffolding`

---

- [x] 2. **Frontend Project Scaffolding (Nuxt 3)**

  **What to do**:
  - Initialize Nuxt 3: `npx nuxi@latest init frontend`
  - Setup Tailwind CSS: `npm install -D @nuxtjs/tailwindcss`
  - Create directory structure: `pages/`, `components/`, `composables/`, `plugins/`, `server/`, `assets/`
  - Setup API client composable with $fetch
  - Add basic layout component with header/footer
  - Configure `nuxt.config.ts` with environment variables

  **Must NOT do**:
  - Don't add unnecessary Nuxt modules
  - Don't implement auth yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with T1, T3, T4)
  - **Blocks**: T7, T13
  - **Blocked By**: None

  **References**:
  - Nuxt 3 docs: https://nuxt.com/docs/getting-started/introduction
  - Nuxt directory structure: https://nuxt.com/docs/guide/directory-structure

  **Acceptance Criteria**:
  - [ ] `npm run dev` starts dev server
  - [ ] Homepage renders with layout
  - [ ] Tailwind classes work correctly

  **QA Scenarios**:
  ```
  Scenario: Homepage renders correctly
    Tool: Playwright
    Preconditions: Dev server running on port 3000
    Steps:
      1. Navigate to http://localhost:3000
      2. Assert page title contains "Medialogg"
      3. Assert layout components visible
    Expected Result: Page loads with header and footer
    Evidence: .sisyphus/evidence/task-2-homepage.png
  ```

  **Commit**: YES
  - Message: `chore(frontend): initial Nuxt 3 scaffolding`

---

- [x] 3. **Database Setup (PostgreSQL + Migrations)**

  **What to do**:
  - Create initial migration files using golang-migrate
  - Setup sqlc configuration: `sqlc.yaml`
  - Create `queries/` directory for SQL queries
  - Add Makefile targets: `make migrate-up`, `make migrate-down`, `make sqlc-generate`
  - Setup database connection pool (pgx)
  - Create initial empty migration for baseline

  **Must NOT do**:
  - Don't add tables yet (separate tasks)
  - Don't use ORM

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with T1, T2, T4)
  - **Blocks**: T5, T9, T11
  - **Blocked By**: None

  **References**:
  - sqlc docs: https://docs.sqlc.dev/en/stable/
  - golang-migrate: https://github.com/golang-migrate/migrate

  **Acceptance Criteria**:
  - [ ] `make migrate-up` runs successfully
  - [ ] `make sqlc-generate` runs without errors
  - [ ] Database connection works with test query

  **QA Scenarios**:
  ```
  Scenario: Database migrations work
    Tool: Bash
    Preconditions: PostgreSQL running (docker-compose)
    Steps:
      1. Run: make migrate-up
      2. Assert output contains "success" or no errors
      3. Run: make migrate-down
      4. Assert migrations rolled back
    Expected Result: Migrations run without errors
    Evidence: .sisyphus/evidence/task-3-migrations.txt
  ```

  **Commit**: YES
  - Message: `chore(db): setup PostgreSQL + sqlc + migrations`

---

- [x] 4. **Docker Compose Setup**

  **What to do**:
  - Create `docker-compose.yml` with services: postgres, backend, frontend
  - Create `Dockerfile.backend` for Go server (multi-stage build)
  - Create `Dockerfile.frontend` for Nuxt (node alpine)
  - Add `.dockerignore` files
  - Setup health checks for each service
  - Add volume mounts for persistence
  - Create development vs production compose overrides

  **Must NOT do**:
  - Don't optimize for production yet
  - Don't add unnecessary services (redis, etc.)

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with T1, T2, T3)
  - **Blocks**: All integration testing
  - **Blocked By**: None

  **References**:
  - Docker Compose docs: https://docs.docker.com/compose/
  - Go Docker best practices: https://docs.docker.com/language/golang/build-images/

  **Acceptance Criteria**:
  - [ ] `docker-compose up` starts all services
  - [ ] Backend health check passes
  - [ ] Frontend accessible on port 3000
  - [ ] PostgreSQL accessible on port 5432

  **QA Scenarios**:
  ```
  Scenario: Docker Compose starts all services
    Tool: Bash
    Preconditions: Docker running
    Steps:
      1. Run: docker-compose up -d
      2. Run: docker-compose ps
      3. Assert all services show "running"
      4. Run: curl http://localhost:8080/health
      5. Run: curl http://localhost:3000
    Expected Result: All services healthy
    Evidence: .sisyphus/evidence/task-4-docker.txt
  ```

  **Commit**: YES
  - Message: `chore: add Docker Compose setup`

---

- [x] 5. **User Model + Migrations**

  **What to do**:
  - Create migration `002_users.up.sql`:
    - `users` table with: id, username, email, password_hash, display_name, avatar_url, bio, is_public, created_at, updated_at
    - Indexes on username, email
  - Create migration `002_users.down.sql` for rollback
  - Add SQL queries to `queries/users.sql`:
    - CreateUser, GetUserByUsername, GetUserByEmail, UpdateUser, DeleteUser
  - Run `make sqlc-generate` to generate Go types
  - Create Go model file `internal/db/models.go` if needed

  **Must NOT do**:
  - Don't add authentication logic yet (T6)
  - Don't add profile settings yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 2 (sequential)
  - **Blocks**: T6, T8
  - **Blocked By**: T1, T3

  **References**:
  - sqlc queries: https://docs.sqlc.dev/en/stable/howto/select.html

  **Acceptance Criteria**:
  - [ ] `make migrate-up` creates users table
  - [ ] `make sqlc-generate` produces `db/queries.go` with user functions
  - [ ] Can insert and query user (via sql query)

  **QA Scenarios**:
  ```
  Scenario: User table exists and queries work
    Tool: Bash (psql)
    Preconditions: Database migrated
    Steps:
      1. Insert test user via SQL
      2. Query users table
      3. Assert user exists with correct schema
    Expected Result: User table with all columns
    Evidence: .sisyphus/evidence/task-5-users.txt
  ```

  **Commit**: YES
  - Message: `feat(db): add users table and queries`

---

- [x] 6. **JWT Authentication Handler + Middleware**

  **What to do**:
  - Install `golang-jwt/jwt/v5` and `bcrypt`
  - Create `internal/middleware/auth.go` with JWT validation middleware
  - Create `internal/api/auth.go` with handlers:
    - `POST /api/auth/register` - creates user, returns token
    - `POST /api/auth/login` - validates credentials, returns token
    - `POST /api/auth/refresh` - refreshes token
    - `GET /api/auth/me` - returns current user (protected)
  - Create `internal/config/jwt.go` for token generation/validation
  - Add password hashing with bcrypt
  - Add input validation with `go-playground/validator`

  **Must NOT do**:
  - Don't add OAuth providers yet
  - Don't add email verification

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 2 (sequential after T5)
  - **Blocks**: T7, T8
  - **Blocked By**: T5

  **References**:
  - golang-jwt: https://github.com/golang-jwt/jwt
  - bcrypt: https://pkg.go.dev/golang.org/x/crypto/bcrypt

  **Acceptance Criteria**:
  - [ ] Register creates user and returns valid JWT
  - [ ] Login returns valid JWT for existing user
  - [ ] Middleware rejects requests without token
  - [ ] Middleware allows requests with valid token
  - [ ] Password is hashed with bcrypt

  **QA Scenarios**:
  ```
  Scenario: User can register and login
    Tool: Bash (curl)
    Preconditions: Backend running
    Steps:
      1. Register: curl -X POST http://localhost:8080/api/auth/register -d '{"username":"testuser","email":"test@test.com","password":"password123"}'
      2. Assert response contains token
      3. Login: curl -X POST http://localhost:8080/api/auth/login -d '{"username":"testuser","password":"password123"}'
      4. Assert response contains token
      5. Get me: curl -H "Authorization: Bearer <token>" http://localhost:8080/api/auth/me
      6. Assert response contains user data
    Expected Result: Auth flow works end-to-end
    Evidence: .sisyphus/evidence/task-6-auth-flow.txt

  Scenario: Auth rejects invalid credentials
    Tool: Bash (curl)
    Preconditions: Backend running
    Steps:
      1. Login with wrong password
      2. Assert HTTP 401
      3. Access protected route without token
      4. Assert HTTP 401
    Expected Result: Auth properly rejects invalid requests
    Evidence: .sisyphus/evidence/task-6-auth-reject.txt
  ```

  **Commit**: YES
  - Message: `feat(auth): implement JWT authentication`

---

- [x] 7. **Frontend Auth Pages (Nuxt)**

  **What to do**:
  - Create `pages/login.vue` with login form
  - Create `pages/register.vue` with registration form
  - Create `composables/useAuth.ts`:
    - `login(username, password)` - calls API, stores token
    - `register(username, email, password)` - calls API, stores token
    - `logout()` - clears token
    - `user` - reactive user state
    - `isAuthenticated` - computed auth state
  - Store token in `useCookie('token')` or localStorage
  - Add auth middleware plugin to auto-add Authorization header
  - Create auth guard for protected routes

  **Must NOT do**:
  - Don't implement other pages yet
  - Don't add profile management

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`/frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 2 (sequential after T6)
  - **Blocks**: T8, T14
  - **Blocked By**: T2, T6

  **References**:
  - Nuxt composables: https://nuxt.com/docs/getting-started/state-management
  - Nuxt useCookie: https://nuxt.com/docs/api/composables/use-cookie

  **Acceptance Criteria**:
  - [x] Register page creates account and redirects to dashboard
  - [x] Login page authenticates and redirects to dashboard
  - [x] Auth state persists on page reload (via cookies)
  - [x] Protected routes redirect to login when not authenticated

  **QA Scenarios**:
  ```
  Scenario: User can register via frontend
    Tool: Playwright
    Preconditions: Backend running, frontend running
    Steps:
      1. Navigate to http://localhost:3000/register
      2. Fill form: username "playwright_user", email "pw@test.com", password "password123"
      3. Click submit
      4. Assert redirected to dashboard
      5. Assert user state populated
    Expected Result: Registration works, user logged in
    Evidence: .sisyphus/evidence/task-7-register.png

  Scenario: User can login via frontend
    Tool: Playwright
    Preconditions: Backend running, user exists from previous scenario
    Steps:
      1. Navigate to http://localhost:3000/login
      2. Fill form: username "playwright_user", password "password123"
      3. Click submit
      4. Assert redirected to dashboard
    Expected Result: Login works, user redirected
    Evidence: .sisyphus/evidence/task-7-login.png
  ```

  **Commit**: YES
  - Message: `feat(frontend): add auth pages and composable`

---

- [ ] 8. **Auth Integration Tests**

  **What to do**:
  - Create `backend/internal/api/auth_test.go`
  - Write Go tests for:
    - Registration success/failure (duplicate username, invalid email)
    - Login success/failure (wrong password, non-existent user)
    - Token refresh success/failure
    - Protected endpoint access
  - Create `frontend/tests/auth.spec.ts` with Playwright tests
  - Add test commands to Makefile
  - Ensure all tests pass

  **Must NOT do**:
  - Don't test features not implemented yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 2 (sequential after T5-7)
  - **Blocks**: Phase 3
  - **Blocked By**: T5, T6, T7

  **References**:
  - Go testing: https://go.dev/doc/tutorial/add-a-test
  - Playwright for Nuxt: https://nuxt.com/docs/getting-started/testing

  **Acceptance Criteria**:
  - [ ] `go test ./internal/api/... -v` passes all auth tests
  - [ ] `npm run test` passes frontend auth tests
  - [ ] All edge cases covered (duplicate, invalid, missing)

  **QA Scenarios**:
  ```
  Scenario: All auth tests pass
    Tool: Bash
    Preconditions: Code complete
    Steps:
      1. Run: cd backend && go test ./internal/api/... -v
      2. Assert all tests pass
      3. Run: cd frontend && npm run test
      4. Assert all tests pass
    Expected Result: 100% test pass rate
    Evidence: .sisyphus/evidence/task-8-tests.txt
  ```

  **Commit**: YES
  - Message: `test: add auth integration tests`

---

- [ ] 9. **Media Model + Migrations**

  **What to do**:
  - Create migration `003_media.up.sql`:
    - `media` table: id, type (FILM, ANIME, BOOK, MANGA, GAME, DOUJIN), title, original_title, description, cover_image, release_date, metadata (JSONB)
    - External IDs: tmdb_id, mal_id, google_books_id, igdb_id
    - Indexes on type, title, external IDs
  - Create migration `003_media.down.sql`
  - Create migration `004_genres.up.sql`:
    - `genres` table: id, name
    - `media_genres` junction table
  - Add SQL queries to `queries/media.sql`:
    - CreateMedia, GetMediaById, GetMediaByTmdbId, SearchMedia, GetMediaByType
  - Run migrations and sqlc generate

  **Must NOT do**:
  - Don't add TMDB sync yet (T19)
  - Don't add log functionality yet

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T11)
  - **Parallel Group**: Phase 3 (can run with T11)
  - **Blocks**: T10, T12, T19
  - **Blocked By**: T3

  **References**:
  - PostgreSQL JSONB: https://www.postgresql.org/docs/current/datatype-json.html
  - sqlc queries: https://docs.sqlc.dev/en/stable/howto/select.html

  **Acceptance Criteria**:
  - [ ] Tables created with correct schema
  - [ ] sqlc generates `CreateMedia`, `GetMediaById`, etc.
  - [ ] Can insert media with JSONB metadata

  **QA Scenarios**:
  ```
  Scenario: Media table accepts JSONB metadata
    Tool: Bash (psql)
    Preconditions: Database migrated
    Steps:
      1. Insert media with metadata: {"runtime": 120, "director": "Test Director"}
      2. Query media by id
      3. Assert metadata field contains correct JSON
    Expected Result: JSONB stored and retrieved correctly
    Evidence: .sisyphus/evidence/task-9-media-jsonb.txt
  ```

  **Commit**: YES
  - Message: `feat(db): add media and genres tables`

---

- [ ] 10. **Media API Handlers**

  **What to do**:
  - Create `internal/api/media.go`:
    - `GET /api/media` - search with pagination (query params: search, type, page, limit)
    - `GET /api/media/:id` - get by ID
    - `POST /api/media` - create (admin only for MVP)
    - `GET /api/media/:id/reviews` - get reviews for media
  - Add pagination helper
  - Add filtering by media type
  - Create response DTOs in `internal/api/dto/`

  **Must NOT do**:
  - Don't add TMDB integration yet
  - Don't add user-generated content yet

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 3 (sequential after T9)
  - **Blocks**: T13, T19
  - **Blocked By**: T9

  **References**:
  - Fiber routing: https://docs.gofiber.io/guide/routing

  **Acceptance Criteria**:
  - [ ] Search returns paginated results
  - [ ] Get by ID returns full media details
  - [ ] Create requires authentication
  - [ ] Reviews endpoint returns empty array initially

  **QA Scenarios**:
  ```
  Scenario: Media search returns results
    Tool: Bash (curl)
    Preconditions: Backend running, test media inserted
    Steps:
      1. Search: curl "http://localhost:8080/api/media?search=test"
      2. Assert response contains media array
      3. Assert pagination fields present
    Expected Result: Search returns paginated results
    Evidence: .sisyphus/evidence/task-10-media-search.txt

  Scenario: Get media by ID
    Tool: Bash (curl)
    Preconditions: Backend running, media exists
    Steps:
      1. Get: curl http://localhost:8080/api/media/{id}
      2. Assert response contains media details with metadata
    Expected Result: Single media object returned
    Evidence: .sisyphus/evidence/task-10-media-get.txt
  ```

  **Commit**: YES
  - Message: `feat(api): add media endpoints`

---

- [ ] 11. **Log Model + Migrations**

  **What to do**:
  - Create migration `005_logs.up.sql`:
    - `logs` table: id, user_id, media_id, status (PLANNED, IN_PROGRESS, COMPLETED, DROPPED), rating, started_at, completed_at, rewatch_count, progress, total, note, contains_spoilers, created_at, updated_at
    - Unique constraint on (user_id, media_id)
    - Indexes on user_id, media_id, created_at
  - Create migration `005_logs.down.sql`
  - Add SQL queries to `queries/logs.sql`:
    - CreateLog, GetLogById, GetUserLogs, GetUserLogForMedia, UpdateLog, DeleteLog, GetLogsByMediaId

  **Must NOT do**:
  - Don't add social features yet (likes, timeline)

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T9)
  - **Parallel Group**: Phase 3
  - **Blocks**: T12, T14
  - **Blocked By**: T3, T5

  **References**:
  - sqlc queries: https://docs.sqlc.dev/en/stable/howto/select.html

  **Acceptance Criteria**:
  - [ ] Logs table created with correct schema
  - [ ] sqlc generates all log query functions
  - [ ] Unique constraint prevents duplicate logs

  **QA Scenarios**:
  ```
  Scenario: User can only have one log per media
    Tool: Bash (psql)
    Preconditions: Database migrated, user exists
    Steps:
      1. Insert log for user_id=X, media_id=Y
      2. Attempt to insert another log for same user and media
      3. Assert unique constraint violation
    Expected Result: Duplicate rejected with error
    Evidence: .sisyphus/evidence/task-11-unique-log.txt
  ```

  **Commit**: YES
  - Message: `feat(db): add logs table`

---

- [ ] 12. **Log API Handlers**

  **What to do**:
  - Create `internal/api/logs.go`:
    - `GET /api/logs/me` - get current user's logs (paginated)
    - `GET /api/logs/me/:mediaId` - get log for specific media
    - `POST /api/logs` - create/update log
    - `PUT /api/logs/:id` - update log
    - `DELETE /api/logs/:id` - delete log
  - Add validation:
    - Rating: 0-10 with 0.5 increments
    - Status: enum validation
  - Create response DTOs
  - Add upsert logic (if log exists for user+media, update it)

  **Must NOT do**:
  - Don't add likes yet
  - Don't add timeline yet

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 3 (after T5, T9, T11)
  - **Blocks**: T14, T16
  - **Blocked By**: T5, T9, T11

  **References**:
  - Fiber routing: https://docs.gofiber.io/guide/routing

  **Acceptance Criteria**:
  - [ ] Create log returns new log with ID
  - [ ] Update log modifies existing log
  - [ ] Delete log removes entry
  - [ ] Get logs returns paginated results
  - [ ] Validation rejects invalid ratings/statuses

  **QA Scenarios**:
  ```
  Scenario: User can create and update log
    Tool: Bash (curl)
    Preconditions: Backend running, user authenticated, media exists
    Steps:
      1. Create: curl -X POST -H "Authorization: Bearer <token>" http://localhost:8080/api/logs -d '{"media_id":"<id>","status":"in_progress","rating":8.5,"note":"Great so far!"}'
      2. Assert log created with ID
      3. Update: curl -X PUT http://localhost:8080/api/logs/<id> -d '{"status":"completed","rating":9.0}'
      4. Assert log updated
      5. Get: curl -H "Authorization: Bearer <token>" http://localhost:8080/api/logs/me
      6. Assert log in results
    Expected Result: Full CRUD for logs
    Evidence: .sisyphus/evidence/task-12-log-crud.txt

  Scenario: Rating validation works
    Tool: Bash (curl)
    Preconditions: Backend running
    Steps:
      1. Create log with rating: 15 (invalid)
      2. Assert HTTP 400 with validation error
      3. Create log with rating: 9.5 (valid)
      4. Assert HTTP 201
    Expected Result: Validation rejects out-of-range rating
    Evidence: .sisyphus/evidence/task-12-validation.txt
  ```

  **Commit**: YES
  - Message: `feat(api): add log endpoints`

---

- [ ] 13. **Frontend Media Pages**

  **What to do**:
  - Create `pages/media/index.vue` - search/browse page
  - Create `pages/media/[id].vue` - media detail page
  - Create `components/MediaCard.vue` - media preview card
  - Create `components/MediaGrid.vue` - grid layout for media
  - Create `composables/useMedia.ts`:
    - `searchMedia(query, type)`
    - `getMedia(id)`
    - `getMediaByType(type)`
  - Add loading states and error handling
  - Style with Tailwind CSS

  **Must NOT do**:
  - Don't add log UI yet (T14)
  - Don't add TMDB search yet (T21)

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`/frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T12)
  - **Parallel Group**: Phase 3
  - **Blocks**: T21
  - **Blocked By**: T2, T10

  **Acceptance Criteria**:
  - [ ] Search page displays media grid
  - [ ] Detail page shows media info
  - [ ] MediaCard shows title, cover, rating badge
  - [ ] Responsive on mobile/desktop

  **QA Scenarios**:
  ```
  Scenario: Media search page works
    Tool: Playwright
    Preconditions: Frontend running, backend running, test media exists
    Steps:
      1. Navigate to http://localhost:3000/media
      2. Type "test" in search box
      3. Assert results display in grid
      4. Click on a media card
      5. Assert navigation to detail page
      6. Assert media details visible
    Expected Result: Search and detail pages functional
    Evidence: .sisyphus/evidence/task-13-media-pages.png
  ```

  **Commit**: YES
  - Message: `feat(frontend): add media pages`

---

- [ ] 14. **Frontend Log UI**

  **What to do**:
  - Create `pages/dashboard.vue` - user dashboard with recent logs
  - Create `components/LogCard.vue` - log display card
  - Create `components/LogForm.vue` - create/edit log form
  - Create `components/LogModal.vue` - modal for quick logging
  - Create `composables/useLogs.ts`:
    - `getMyLogs()`
    - `getLogForMedia(mediaId)`
    - `createLog(data)`
    - `updateLog(id, data)`
    - `deleteLog(id)`
  - Add rating input component (5-star or 10-point toggleable)
  - Add status dropdown (planned, in-progress, completed, dropped)
  - Add progress tracker for episodic media

  **Must NOT do**:
  - Don't add social features yet
  - Don't add reviews yet

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`/frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 3 (after T7, T12)
  - **Blocks**: T18
  - **Blocked By**: T7, T12

  **References**:
  - Nuxt pages: https://nuxt.com/docs/guide/directory-structure/pages

  **Acceptance Criteria**:
  - [ ] Dashboard shows user's recent logs
  - [ ] LogForm creates/updates logs
  - [ ] LogCard displays status, rating, progress
  - [ ] Modal allows quick logging from media detail page

  **QA Scenarios**:
  ```
  Scenario: User can create log from dashboard
    Tool: Playwright
    Preconditions: Frontend running, user logged in, media exists
    Steps:
      1. Navigate to dashboard
      2. Click "Log new media"
      3. Search for media
      4. Select status: "in-progress"
      5. Set rating: 8.5
      6. Add note: "Really enjoying this!"
      7. Submit form
      8. Assert log appears in dashboard
    Expected Result: Log created and displayed
    Evidence: .sisyphus/evidence/task-14-create-log.png

  Scenario: User can update existing log
    Tool: Playwright
    Preconditions: Frontend running, user logged in, log exists
    Steps:
      1. Navigate to dashboard
      2. Click on existing log
      3. Change status to "completed"
      4. Update rating to 9.0
      5. Save changes
      6. Assert log updated
    Expected Result: Log updated successfully
    Evidence: .sisyphus/evidence/task-14-update-log.png
  ```

  **Commit**: YES
  - Message: `feat(frontend): add log UI components`

---

- [ ] 15. **Follow Model + Handlers**

  **What to do**:
  - Create migration `006_follows.up.sql`:
    - `follows` table: id, follower_id, following_id, created_at
    - Unique constraint on (follower_id, following_id)
    - Check: follower_id != following_id
    - Indexes on follower_id, following_id
  - Create migration `006_follows.down.sql`
  - Add SQL queries to `queries/follows.sql`:
    - FollowUser, UnfollowUser, GetFollowers, GetFollowing, IsFollowing
  - Create `internal/api/follows.go`:
    - `POST /api/users/:username/follow`
    - `DELETE /api/users/:username/follow`
    - `GET /api/users/:username/followers`
    - `GET /api/users/:username/following`

  **Must NOT do**:
  - Don't add timeline yet (T16)
  - Don't add notifications

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T17)
  - **Parallel Group**: Wave 4
  - **Blocks**: T16, T18
  - **Blocked By**: T5

  **References**:
  - PostgreSQL foreign keys: https://www.postgresql.org/docs/current/ddl-constraints.html

  **Acceptance Criteria**:
  - [ ] Follow creates relationship
  - [ ] Unfollow removes relationship
  - [ ] Can't follow self
  - [ ] Can't follow twice

  **QA Scenarios**:
  ```
  Scenario: Follow/unfollow works
    Tool: Bash (curl)
    Preconditions: Backend running, two users exist
    Steps:
      1. UserA follows UserB
      2. Assert following list contains UserB
      3. UserA unfollows UserB
      4. Assert following list empty
    Expected Result: Follow toggle works
    Evidence: .sisyphus/evidence/task-15-follow.txt

  Scenario: Cant follow self
    Tool: Bash (curl)
    Preconditions: Backend running, user authenticated
    Steps:
      1. User attempts to follow themselves
      2. Assert HTTP 400 with error
    Expected Result: Self-follow rejected
    Evidence: .sisyphus/evidence/task-15-self-follow.txt
  ```

  **Commit**: YES
  - Message: `feat(api): add follow functionality`

---

- [ ] 16. **Timeline API**

  **What to do**:
  - Create `internal/api/timeline.go`:
    - `GET /api/timeline` - get activity from followed users (paginated)
  - Add SQL queries to `queries/timeline.sql`:
    - GetTimeline (joins logs + users + media + follows)
  - Implement timeline query:
    ```sql
    SELECT l.*, u.username, m.title, m.cover_image
    FROM logs l
    JOIN users u ON l.user_id = u.id
    JOIN media m ON l.media_id = m.id
    WHERE l.user_id IN (SELECT following_id FROM follows WHERE follower_id = $1)
    ORDER BY l.created_at DESC
    LIMIT $2 OFFSET $3
    ```
  - Add caching for timeline (optional, use Redis later if needed)

  **Must NOT do**:
  - Don't add algorithmic ranking (chronological for MVP)
  - Don't add comments yet

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (sequential after T5, T12, T15)
  - **Blocks**: T18
  - **Blocked By**: T5, T12, T15

  **References**:
  - PostgreSQL JOINs: https://www.postgresql.org/docs/current/queries-table-expressions.html

  **Acceptance Criteria**:
  - [ ] Timeline returns followed users' logs
  - [ ] Timeline excludes own logs
  - [ ] Timeline paginated correctly

  **QA Scenarios**:
  ```
  Scenario: Timeline shows followed users activity
    Tool: Bash (curl)
    Preconditions: Backend running, UserA follows UserB, UserB has logs
    Steps:
      1. Get timeline as UserA
      2. Assert UserB's logs appear
      3. Assert proper pagination
    Expected Result: Timeline contains followed users' logs
    Evidence: .sisyphus/evidence/task-16-timeline.txt

  Scenario: Timeline excludes non-followed users
    Tool: Bash (curl)
    Preconditions: Backend running, UserA does NOT follow UserC, UserC has logs
    Steps:
      1. Get timeline as UserA
      2. Assert UserC's logs do NOT appear
    Expected Result: Only followed users' logs
    Evidence: .sisyphus/evidence/task-16-timeline-filter.txt
  ```

  **Commit**: YES
  - Message: `feat(api): add timeline endpoint`

---

- [ ] 17. **Like Functionality**

  **What to do**:
  - Create migration `007_likes.up.sql`:
    - `likes` table: id, user_id, log_id (nullable), review_id (nullable), created_at
    - Check constraint: exactly one of log_id or review_id
    - Unique on (user_id, log_id) and (user_id, review_id)
  - Add SQL queries to `queries/likes.sql`:
    - LikeLog, UnlikeLog, GetLogLikes, IsLogLiked
  - Create `internal/api/likes.go`:
    - `POST /api/logs/:id/like`
    - `DELETE /api/logs/:id/like`
  - Add like count to log responses

  **Must NOT do**:
  - Don't add review likes yet (wait for review feature)
  - Don't add notifications

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T15)
  - **Parallel Group**: Wave 4
  - **Blocks**: None
  - **Blocked By**: T5, T12

  **Acceptance Criteria**:
  - [ ] Like adds to log
  - [ ] Unlike removes from log
  - [ ] Can't like twice
  - [ ] Like count returned in log response

  **QA Scenarios**:
  ```
  Scenario: Like/unlike works
    Tool: Bash (curl)
    Preconditions: Backend running, user authenticated, log exists
    Steps:
      1. Like a log
      2. Assert like count increases
      3. Unlike the log
      4. Assert like count decreases
    Expected Result: Like toggle works
    Evidence: .sisyphus/evidence/task-17-like.txt
  ```

  **Commit**: YES
  - Message: `feat(api): add like functionality`

---

- [ ] 18. **Frontend Social Pages**

  **What to do**:
  - Create `pages/timeline.vue` - shows followed users' activity
  - Create `pages/[username].vue` - public user profile
  - Create `components/FollowButton.vue` - follow/unfollow toggle
  - Create `components/UserCard.vue` - user preview card
  - Create `composables/useFollow.ts`:
    - `followUser(username)`
    - `unfollowUser(username)`
    - `isFollowing(username)`
  - Create `composables/useTimeline.ts`:
    - `getTimeline(page)`
  - Update `LogCard` to show like button
  - Add profile stats: logs count, followers, following

  **Must NOT do**:
  - Don't add profile editing yet
  - Don't add private accounts yet

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`/frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (after T14, T16)
  - **Blocks**: None
  - **Blocked By**: T7, T14, T16

  **Acceptance Criteria**:
  - [ ] Timeline shows followed users' logs
  - [ ] Profile page shows user stats and logs
  - [ ] Follow button toggles follow state
  - [ ] Like button on log cards works

  **QA Scenarios**:
  ```
  Scenario: Timeline page works
    Tool: Playwright
    Preconditions: Frontend running, user logged in, following users with logs
    Steps:
      1. Navigate to timeline
      2. Assert followed users' logs appear
      3. Click like on a log
      4. Assert like state changes
    Expected Result: Timeline functional with likes
    Evidence: .sisyphus/evidence/task-18-timeline-page.png

  Scenario: User profile page works
    Tool: Playwright
    Preconditions: Frontend running, users exist
    Steps:
      1. Navigate to /username
      2. Assert profile shows avatar, stats, logs
      3. Click follow button
      4. Assert button changes to "Following"
      5. Assert follower count updates
    Expected Result: Profile and follow work
    Evidence: .sisyphus/evidence/task-18-profile.png
  ```

  **Commit**: YES
  - Message: `feat(frontend): add social pages`

---

- [ ] 19. **TMDB Client Integration**

  **What to do**:
  - Create `internal/tmdb/client.go`:
    - `SearchMovies(query string, page int)` - search TMDB
    - `GetMovie(id int)` - get movie details
    - `GetMovieCredits(id int)` - get cast/crew
    - `GetPopularMovies(page int)` - get popular list
  - Add TMDB API configuration
  - Add rate limiting (TMDB has rate limits)
  - Add image URL helper (poster URLs need base URL)
  - Add error handling for API failures
  - Create `internal/tmdb/types.go` for TMDB response structs

  **Must NOT do**:
  - Don't sync all movies (API limits)
  - Don't add other media types yet

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 5 (sequential after T9, T10)
  - **Blocks**: T20, T21
  - **Blocked By**: T9, T10

  **References**:
  - TMDB API: https://developers.themoviedb.org/3
  - Rate limits: 40 requests per 10 seconds

  **Acceptance Criteria**:
  - [ ] SearchMovies returns parsed results
  - [ ] GetMovie returns full details
  - [ ] Rate limiting implemented
  - [ ] Errors handled gracefully

  **QA Scenarios**:
  ```
  Scenario: TMDB search works
    Tool: Bash (go test)
    Preconditions: TMDB_API_KEY set
    Steps:
      1. Run: go test ./internal/tmdb/... -v -run TestSearchMovies
      2. Assert returns movie results
    Expected Result: TMDB search returns valid data
    Evidence: .sisyphus/evidence/task-19-tmdb-search.txt

  Scenario: Rate limiting prevents API abuse
    Tool: Bash (go test)
    Preconditions: TMDB client initialized
    Steps:
      1. Make 50 rapid requests
      2. Assert rate limiting kicks in
      3. Assert no API errors
    Expected Result: Rate limiting works
    Evidence: .sisyphus/evidence/task-19-rate-limit.txt
  ```

  **Commit**: YES
  - Message: `feat(integration): add TMDB client`

---

- [ ] 20. **TMDB Sync Background Jobs**

  **What to do**:
  - Create `internal/jobs/sync.go`:
    - `SyncPopularMovies()` - sync top popular movies weekly
    - `SyncMovieDetails(id int)` - fetch and store full details
  - Add background job scheduler (use `robfig/cron` or simple ticker)
  - Store synced movies in `media` table
  - Add genres from TMDB to `genres` table
  - Add sync status tracking (last sync time)
  - Add CLI command: `make sync-movies`

  **Must NOT do**:
  - Don't sync all movies (storage/API limits)
  - Don't run too frequently (TMDB rate limits)

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 5 (sequential after T19)
  - **Blocks**: T21
  - **Blocked By**: T19

  **Acceptance Criteria**:
  - [ ] Sync stores movies in database
  - [ ] Genres populated from TMDB
  - [ ] Scheduled job runs without blocking
  - [ ] Sync can be triggered manually

  **QA Scenarios**:
  ```
  Scenario: Sync stores movie data
    Tool: Bash
    Preconditions: Backend running, TMDB_API_KEY set
    Steps:
      1. Run: make sync-movies
      2. Query: SELECT COUNT(*) FROM media WHERE type = 'film'
      3. Assert count > 0
      4. Query: SELECT COUNT(*) FROM genres
      5. Assert genres populated
    Expected Result: Movies synced to database
    Evidence: .sisyphus/evidence/task-20-sync.txt
  ```

  **Commit**: YES
  - Message: `feat(jobs): add TMDB sync jobs`

---

- [ ] 21. **Frontend TMDB Search Integration**

  **What to do**:
  - Create `composables/useTmdbSearch.ts`:
    - `searchMovies(query, page)`
    - `getPopularMovies(page)`
  - Update `pages/media/index.vue`:
    - Add search input
    - Show TMDB results when searching
    - Show local database when browsing
  - Create `components/MediaSearch.vue`:
    - Debounced search input
    - Results grid with "Log this" button
  - Add "Log this movie" flow:
    - Click result → open log modal → quick log
  - Add loading states and empty states

  **Must NOT do**:
  - Don't implement advanced filters yet
  - Don't add other media types yet

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
  - **Skills**: [`/frontend-ui-ux`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Phase 5 (sequential after T13, T19)
  - **Blocks**: Final integration
  - **Blocked By**: T13, T19

  **References**:
  - TMDB Search API: https://developers.themoviedb.org/3/search/search-movies

  **Acceptance Criteria**:
  - [ ] Search returns TMDB results
  - [ ] Clicking result opens log modal
  - [ ] Can log movie not in database (creates on-demand)
  - [ ] Debounced search prevents API spam

  **QA Scenarios**:
  ```
  Scenario: User can search TMDB and log
    Tool: Playwright
    Preconditions: Frontend running, backend running, TMDB_API_KEY set
    Steps:
      1. Navigate to /media
      2. Type "Inception" in search box
      3. Assert TMDB results appear
      4. Click on "Inception" result
      5. Assert log modal opens with movie data
      6. Set rating, status, submit
      7. Assert log created
    Expected Result: TMDB search and log flow works
    Evidence: .sisyphus/evidence/task-21-tmdb-search.png
  ```

  **Commit**: YES
  - Message: `feat(frontend): integrate TMDB search`

---

## Final Verification Wave

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Read plan end-to-end. Verify all "Must Have" implemented. Check no "Must NOT Have" patterns exist. Verify evidence files in .sisyphus/evidence/.
  Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT: APPROVE/REJECT`

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Run `go test ./...` + `go vet ./...` for backend. Run `npm run test` + `npm run lint` for frontend. Review for common issues.
  Output: `Backend Tests [PASS/FAIL] | Frontend Tests [PASS/FAIL] | Lint [PASS/FAIL] | VERDICT`

- [ ] F3. **Integration QA** — `unspecified-high`
  Start Docker Compose. Execute curl tests against API. Test frontend in browser. Save evidence to `.sisyphus/evidence/final-qa/`.
  Output: `API Endpoints [N/N] | Frontend Pages [N/N] | Integration [PASS/FAIL] | VERDICT`

- [ ] F4. **Scope Fidelity Check** — `deep`
  Compare implemented features vs plan. No scope creep. No missing requirements.
  Output: `Features [N/N] | Scope Creep [NONE/N issues] | VERDICT`

---

## Commit Strategy

- **Phase 1 (T1-T4)**: `chore: initial project scaffolding`
- **Phase 2 (T5-T8)**: `feat: authentication system`
- **Phase 3 (T9-T14)**: `feat: core media logging`
- **Phase 4 (T15-T18)**: `feat: social features`
- **Phase 5 (T19-T21)**: `feat: TMDB integration`

---

## Success Criteria

### Verification Commands
```bash
# Backend tests
cd backend && go test ./... -v

# Frontend tests
cd frontend && npm run test

# Integration
docker-compose up -d
curl http://localhost:8080/health
curl http://localhost:3000

# Auth flow
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@test.com","password":"password123"}'
```

### Final Checklist
- [ ] User can register and login
- [ ] User can search films via TMDB
- [ ] User can log films with status/rating
- [ ] User can follow other users
- [ ] Timeline shows followed users' activity
- [ ] Docker Compose starts all services
- [ ] All tests pass
- [ ] All evidence captured
