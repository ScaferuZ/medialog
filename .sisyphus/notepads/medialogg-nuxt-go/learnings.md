2026-04-13: Backend scaffold uses a minimal Fiber app with /health and Viper env defaults; keep future additions isolated in internal/.
2026-04-13: Database scaffolding should stay migration/query-only until later phases; sqlc config points at migrations/ and internal/db/queries/.
- Nuxt init in this environment required creating frontend/ first, then running npx nuxi init inside it.
- Tailwind module worked with assets/css/tailwind.css and nuxt.config.ts css entry.
2026-04-13: sqlc generation can be run via Docker when local sqlc/go binaries are unavailable; gopls is not installed here, so LSP verification is blocked until the Go toolchain is present.
