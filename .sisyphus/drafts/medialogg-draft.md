# Draft: Medialogg Planning Session

## Requirements (confirmed)

- **Platform**: Unified media logging combining Letterboxd, Goodreads, Backloggd
- **Media types**: Films, anime, books, manga, games, doujin
- **Core features**: Log, review, rate, follow, timeline
- **Priority media**: Films first (TMDB integration)
- **Hosting**: Self-hosted homelab
- **Doujin handling**: Content warnings (no moderation system for MVP)

## Technical Decisions

- **Frontend**: Nuxt 3 (uses Nitro internally, Vue ecosystem)
- **Backend**: Go with Fiber framework
- **Database**: PostgreSQL with sqlc (type-safe SQL, no ORM)
- **API Style**: REST with OpenAPI spec
- **Auth**: JWT (standard, no OAuth initially)
- **Search**: TMDB for films (API limits: 40 req/10s)

## Technical Decisions Rationale

### Why Nuxt 3 over SvelteKit?
- User preference for NitroJS approach
- Nuxt 3 uses Nitro internally (same engine)
- Vue ecosystem familiar to many developers
- File-based routing, SSR/SSG built-in

### Why Go over Node/Python?
- Single binary deployment (perfect for homelab)
- Low memory usage
- Better concurrency
- Type safety without TypeScript complexity

### Why sqlc over ORM?
- Performance (no ORM overhead)
- Type safety from SQL
- More control
- Go-idiomatic

## Open Questions

## Scope Boundaries

- **INCLUDE**:
  - Authentication (JWT)
  - Media logging (all statuses)
  - Follow/timeline system
  - TMDB integration
  - User profiles
  - Lists functionality
  
- **EXCLUDE** (MVP):
  - OAuth providers (add later)
  - Content moderation
  - Advanced search/filters
  - Import from other platforms
  - Notifications
  - Mobile app
  - Other media types (anime, books, games - add sequentially)
