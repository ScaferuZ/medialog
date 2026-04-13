# Medialogg - Unified Media Logging Platform

## Executive Summary

Medialogg is a unified platform for logging, reviewing, and sharing media consumption across multiple formats (films, anime, books, manga, games, doujin). It combines the social features of Letterboxd, Goodreads, and Backloggd into a single cohesive experience.

---

## 1. Tech Stack Recommendations

### Option A: Modern Full-Stack TypeScript (Recommended)

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| **Frontend** | Next.js 14+ (App Router) | SSR/SSG support, API routes, excellent DX |
| **Styling** | Tailwind CSS + shadcn/ui | Rapid UI development, accessible components |
| **Backend** | Next.js API Routes + tRPC | Type-safe APIs, colocated with frontend |
| **Database** | PostgreSQL + Prisma | Relational data fits media relationships well |
| **Auth** | NextAuth.js or Lucia | Flexible auth with social providers |
| **Search** | Meilisearch or Algolia | Fast, typo-tolerant media search |
| **File Storage** | Cloudflare R2 or AWS S3 | User uploads, media images |
| **Cache** | Redis (Upstash) | Sessions, rate limiting, trending data |

### Option B: Separated Backend (For Mobile Apps Later)

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| **Frontend** | Next.js 14+ or Nuxt 3 | Modern React/Vue framework |
| **Backend** | Fastify (Node) or Fiber (Go) | High-performance API server |
| **API** | GraphQL (Pothos) or REST | Flexible data fetching |
| **Database** | PostgreSQL + Prisma/Drizzle | Same as Option A |
| **Auth** | Custom JWT or Auth0 | Token-based for mobile support |

### Recommendation

**Start with Option A** (Next.js full-stack). It provides:
- Faster initial development
- Type safety across frontend/backend
- Easy deployment to Vercel
- Can extract to separate backend later if needed

---

## 2. Database Schema Design

### Core Entities

```prisma
// Media Types (polymorphic base)
model Media {
  id          String   @id @default(cuid())
  type        MediaType // FILM, ANIME, BOOK, MANGA, GAME, DOUJIN
  title       String
  originalTitle String? // For anime/manga in Japanese
  description String?  @db.Text
  coverImage  String?
  releaseDate DateTime?
  genres      Genre[]
  
  // External IDs for syncing
  tmdbId      String?  // For films
  malId       Int?     // MyAnimeList for anime/manga
  googleBooksId String? // For books
  igdbId      Int?     // For games
  
  // Type-specific fields (JSON for flexibility)
  metadata    Json?    // Runtime, episodes, pages, platform, etc.
  
  createdAt   DateTime @default(now())
  updatedAt   DateTime @updatedAt
  
  // Relations
  logs        Log[]
  reviews     Review[]
  
  @@index([type])
  @@index([title])
  @@index([tmdbId])
  @@index([malId])
}

model Genre {
  id    String @id @default(cuid())
  name  String @unique
  media Media[]
}

// Users
model User {
  id            String    @id @default(cuid())
  username      String    @unique
  email         String    @unique
  passwordHash  String?   // For credentials auth
  displayName   String?
  avatar        String?
  bio           String?   @db.Text
  
  // Profile settings
  isPublic      Boolean   @default(true)
  
  createdAt     DateTime  @default(now())
  updatedAt     DateTime  @updatedAt
  
  // Relations
  logs          Log[]
  reviews       Review[]
  followers     Follow[]  @relation("following")
  following     Follow[]  @relation("followers")
  lists         List[]
  likes         Like[]
}

model Follow {
  id          String @id @default(cuid())
  followerId  String
  followingId String
  follower    User   @relation("followers", fields: [followerId], references: [id], onDelete: Cascade)
  following   User   @relation("following", fields: [followingId], references: [id], onDelete: Cascade)
  createdAt   DateTime @default(now())
  
  @@unique([followerId, followingId])
  @@index([followerId])
  @@index([followingId])
}

// Activity Logs (The core feature)
model Log {
  id          String    @id @default(cuid())
  userId      String
  mediaId     String
  
  // Log details
  status      LogStatus // COMPLETED, IN_PROGRESS, PLANNED, DROPPED
  rating      Float?    // 0.5-5 or 0-10 based on preference
  startedAt   DateTime?
  completedAt DateTime?
  rewatch     Int       @default(0) // Rewatch/reread/replay count
  
  // Progress tracking
  progress    Int?      // Pages read, episodes watched, hours played
  total       Int?      // Total pages/episodes (denormalized for performance)
  
  // User content
  note        String?   @db.Text // Quick note (not a full review)
  containsSpoilers Boolean @default(false)
  
  createdAt   DateTime  @default(now())
  updatedAt   DateTime  @updatedAt
  
  // Relations
  user        User      @relation(fields: [userId], references: [id], onDelete: Cascade)
  media       Media     @relation(fields: [mediaId], references: [id], onDelete: Cascade)
  review      Review?   // Optional linked review
  likes       Like[]
  
  @@unique([userId, mediaId]) // One log per media per user
  @@index([userId, createdAt])
  @@index([mediaId])
}

enum LogStatus {
  PLANNED
  IN_PROGRESS
  COMPLETED
  DROPPED
}

enum MediaType {
  FILM
  ANIME
  BOOK
  MANGA
  GAME
  DOUJIN
}

// Reviews (Separate from logs for flexibility)
model Review {
  id          String   @id @default(cuid())
  userId      String
  mediaId     String
  logId       String?  @unique // Optional link to a log
  
  // Content
  title       String?
  content     String   @db.Text
  rating      Float
  containsSpoilers Boolean @default(false)
  
  // Engagement
  likes       Like[]
  
  createdAt   DateTime @default(now())
  updatedAt   DateTime @updatedAt
  
  // Relations
  user        User     @relation(fields: [userId], references: [id], onDelete: Cascade)
  media       Media    @relation(fields: [mediaId], references: [id], onDelete: Cascade)
  log         Log?     @relation(fields: [logId], references: [id])
  
  @@index([userId, createdAt])
  @@index([mediaId, createdAt])
}

// Social Features
model Like {
  id        String   @id @default(cuid())
  userId    String
  logId     String?
  reviewId  String?
  createdAt DateTime @default(now())
  
  user      User     @relation(fields: [userId], references: [id], onDelete: Cascade)
  log       Log?     @relation(fields: [logId], references: [id], onDelete: Cascade)
  review    Review?  @relation(fields: [reviewId], references: [id], onDelete: Cascade)
  
  @@unique([userId, logId])
  @@unique([userId, reviewId])
}

// User Lists ("Best Horror Movies", "Summer Reading", etc.)
model List {
  id          String     @id @default(cuid())
  userId      String
  title       String
  description String?    @db.Text
  isPublic    Boolean    @default(true)
  createdAt   DateTime   @default(now())
  updatedAt   DateTime   @updatedAt
  
  user        User       @relation(fields: [userId], references: [id], onDelete: Cascade)
  items       ListItem[]
  
  @@index([userId])
}

model ListItem {
  id        String   @id @default(cuid())
  listId    String
  mediaId   String
  order     Int      @default(0)
  note      String?
  addedAt   DateTime @default(now())
  
  list      List     @relation(fields: [listId], references: [id], onDelete: Cascade)
  
  @@unique([listId, mediaId])
}
```

---

## 3. API Structure

### Authentication Endpoints
```
POST   /api/auth/register
POST   /api/auth/login
POST   /api/auth/logout
POST   /api/auth/refresh
GET    /api/auth/me
```

### User Endpoints
```
GET    /api/users/:username          // Public profile
PUT    /api/users/me                 // Update profile
GET    /api/users/me/stats           // Stats (total watched, etc.)
POST   /api/users/:username/follow
DELETE /api/users/:username/follow
GET    /api/users/:username/followers
GET    /api/users/:username/following
```

### Media Endpoints
```
GET    /api/media                    // Search with filters
GET    /api/media/:id
GET    /api/media/:id/reviews
GET    /api/media/tmdb/:tmdbId       // Find by external ID
GET    /api/media/mal/:malId
POST   /api/media                    // Admin: add media
PUT    /api/media/:id                // Admin: update media
```

### Log Endpoints (Core Activity)
```
GET    /api/logs                     // Timeline (following)
GET    /api/logs/me                  // My logs
GET    /api/logs/me/:mediaId         // Log for specific media
POST   /api/logs                     // Create/update log
PUT    /api/logs/:id
DELETE /api/logs/:id
POST   /api/logs/:id/like
```

### Review Endpoints
```
GET    /api/reviews                  // Recent reviews
GET    /api/reviews/me
GET    /api/reviews/:id
POST   /api/reviews
PUT    /api/reviews/:id
DELETE /api/reviews/:id
POST   /api/reviews/:id/like
```

### List Endpoints
```
GET    /api/lists                    // Public lists
GET    /api/lists/me
GET    /api/lists/:id
POST   /api/lists
PUT    /api/lists/:id
DELETE /api/lists/:id
POST   /api/lists/:id/items
DELETE /api/lists/:id/items/:itemId
```

---

## 4. Implementation Phases

### Phase 1: Foundation (Week 1-2)
- [ ] Project setup (Next.js + Prisma + PostgreSQL)
- [ ] Authentication system (NextAuth.js)
- [ ] Basic user profiles
- [ ] Database schema implementation
- [ ] Deployment pipeline (Vercel + Neon/Railway)

### Phase 2: Core Logging (Week 3-4)
- [ ] Media database (manual entry for MVP)
- [ ] Log creation/editing UI
- [ ] Status tracking (planned, watching, completed, dropped)
- [ ] Rating system
- [ ] Basic profile pages with stats

### Phase 3: Social Features (Week 5-6)
- [ ] Follow/unfollow system
- [ ] Activity timeline
- [ ] Like functionality
- [ ] User search
- [ ] Public/private profiles

### Phase 4: Reviews & Lists (Week 7-8)
- [ ] Review system (separate from logs)
- [ ] User lists functionality
- [ ] Rich text editor for reviews
- [ ] Spoiler tagging

### Phase 5: External APIs & Discovery (Week 9-10)
- [ ] TMDB integration (films)
- [ ] MyAnimeList API (anime/manga)
- [ ] Google Books API (books)
- [ ] IGDB API (games)
- [ ] Auto-complete search
- [ ] Media discovery/browse

### Phase 6: Polish & Advanced (Week 11-12)
- [ ] Advanced filters
- [ ] Statistics/analytics page
- [ ] Import from Letterboxd/Goodreads
- [ ] Dark mode
- [ ] Mobile responsiveness polish
- [ ] Performance optimization

---

## 5. External API Integrations

### Film Data: TMDB (The Movie Database)
- **Free tier**: 40 requests/10 seconds
- **Features**: Movie/TV data, images, cast, ratings
- **Use case**: Auto-populate film details, posters, metadata

### Anime/Manga Data: MyAnimeList API
- **Free tier**: 100 requests/day (basic)
- **Features**: Anime/manga database, rankings
- **Use case**: Anime/manga metadata

### Alternative: AniList API (GraphQL)
- **Free tier**: 90 requests/minute
- **Features**: More detailed anime data, user lists
- **Use case**: Better for anime-focused features

### Book Data: Google Books API
- **Free tier**: 1000 requests/day
- **Features**: Book search, covers, descriptions
- **Use case**: Book metadata and discovery

### Game Data: IGDB (via Twitch API)
- **Free tier**: API key required, reasonable limits
- **Features**: Game database, covers, platforms, ratings
- **Use case**: Video game metadata

---

## 6. Key Technical Decisions

### Media Type Handling
Use a single `Media` table with a `type` enum rather than separate tables. This allows:
- Unified search across all media types
- Single log/review table
- Easier to add new media types later

### Rating System
Support both 5-star (0.5 increments) and 10-point scales. Store normalized 0-100 internally:
```typescript
// User preference: ratingScale: '5star' | '10point'
// Display: convert from 0-100 to preferred scale
// Store: always 0-100
```

### Timeline Algorithm
Simple chronological for MVP:
```sql
SELECT * FROM Log 
WHERE userId IN (SELECT followingId FROM Follow WHERE followerId = ?)
ORDER BY createdAt DESC
LIMIT 50
```

Later: Add weighted scoring (likes, comments, recency).

### Image Storage
- Use external CDN URLs when available (TMDB, MAL, etc.)
- Upload user avatars to R2/S3
- Image optimization via Next.js Image component

---

## 7. Open Questions

1. **Doujin content**: How to handle potentially NSFW content? Content warnings? Age verification?

2. **Import functionality**: Should we support importing from Letterboxd/Goodreads/Backloggd exports?

3. **Monetization**: Free with premium features? Ads? Completely free?

4. **Moderation**: User reports, review moderation, content flags?

5. **Mobile app**: Is a native mobile app planned, or PWA sufficient?

---

## 8. Next Steps

1. **Confirm tech stack** - Do you prefer Next.js or another framework?
2. **Set up project** - Initialize repo with chosen stack
3. **Database setup** - Provision PostgreSQL instance
4. **Start Phase 1** - Authentication and basic structure

Let me know your preferences and we can start implementing!
