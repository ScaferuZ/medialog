2026-04-13: Kept the scaffold minimal: Fiber v2 HTTP server, Viper config loader, placeholder db/api packages, and no ORM/auth wiring.
2026-04-13: Added pgx/v5 pool setup with retry logic plus sqlc/migrate scaffolding, but left schema baseline empty per phase guidance.
- Used Nuxt 3 SSR with App Router, Tailwind module, minimal global layout, and a small API composable around $fetch.
