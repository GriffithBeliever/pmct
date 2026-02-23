# CLAUDE.md — Personal Media Collection Tracker

## Project Overview

A full-stack Personal Media Collection Tracker supporting movies, music, and games.
Users can manage their collections, track status, search/filter, share profiles, and
leverage AI-powered features (recommendations, auto-metadata, mood matching).

**Stack:** Go (backend API) + Next.js 14 App Router (frontend)

---

## Architecture

```
ems/
├── backend/           # Go microservice (REST + optional gRPC)
│   ├── cmd/server/    # Entry point
│   ├── internal/
│   │   ├── auth/      # JWT + session management
│   │   ├── media/     # Media CRUD domain
│   │   ├── ai/        # AI feature integrations
│   │   ├── search/    # Full-text search logic
│   │   └── db/        # Database layer (PostgreSQL)
│   ├── go.mod
│   └── go.sum
├── frontend/          # Next.js 14 App Router
│   ├── app/           # RSC-first pages and layouts
│   ├── components/    # Shared UI components
│   ├── lib/           # API clients, utilities
│   └── package.json
├── CLAUDE.md
└── README.md
```

---

## Go Backend Rules (golang-pro)

### MUST DO
- Use **Go 1.21+** features (slog, slices, maps packages).
- Add `context.Context` as the **first parameter** to every blocking/IO function.
- Handle **all errors explicitly** — never discard with `_` without a comment.
- Wrap errors with `fmt.Errorf("operation: %w", err)` for stack context.
- Write **table-driven tests** with `t.Run` subtests; always run with `-race`.
- Document all exported types, functions, and packages with GoDoc comments.
- Use **functional options** pattern for service/server configuration.
- Use **small, focused interfaces** — prefer `io.Reader`/`io.Writer` composition.
- Use generics (`X | Y` constraints) for reusable collection utilities.
- Format with `gofmt`; lint with `golangci-lint`.

### MUST NOT DO
- No `panic` for normal error handling.
- No goroutines without a clear lifecycle (use `errgroup`, `context` cancellation).
- No hardcoded config — use env vars or a config struct with functional options.
- No naked returns.
- No reflection without performance justification.

### Concurrency Patterns
- Use `golang.org/x/sync/errgroup` for parallel API/DB calls.
- Use buffered channels with explicit capacity; document the rationale.
- Use `sync.Once` for singleton initialization (DB pool, AI client).
- Always select on `ctx.Done()` in long-running goroutines.

### Example: Parallel metadata fetch
```go
// Fetch AI metadata and DB record in parallel
g, ctx := errgroup.WithContext(ctx)

var aiMeta *ai.Metadata
var item *media.Item

g.Go(func() error {
    var err error
    aiMeta, err = ai.FetchMetadata(ctx, title)
    return fmt.Errorf("ai metadata: %w", err)
})

g.Go(func() error {
    var err error
    item, err = repo.GetByID(ctx, id)
    return fmt.Errorf("get media: %w", err)
})

if err := g.Wait(); err != nil {
    return nil, err
}
```

---

## Next.js / React Frontend Rules (vercel-react-best-practices)

### Data Fetching — Eliminate Waterfalls (CRITICAL)
- **Prefer Server Components (RSC)** for data fetching — avoid client-side waterfalls.
- Use `Promise.all()` for independent parallel fetches in Server Components.
- Use `React.cache()` for per-request deduplication of repeated fetches.
- Use `<Suspense>` boundaries to stream content progressively.
- Start promises early, `await` late (defer-await pattern).

```tsx
// Good — parallel, no waterfall
const [media, profile] = await Promise.all([
  fetchMedia(id),
  fetchProfile(userId),
]);
```

### Bundle Size (CRITICAL)
- Import directly from source, **never from barrel files** (`index.ts`).
- Use `next/dynamic` with `ssr: false` for heavy client components (editors, charts).
- Load analytics/AI SDK scripts **after hydration** with `next/script strategy="lazyOnload"`.
- Load feature modules conditionally (e.g., AI recommendation panel only when opened).
- Preload on hover/focus for perceived speed (`router.prefetch` on hover).

```tsx
// Good — dynamic import for heavy AI recommendation panel
const RecommendationPanel = dynamic(() => import('@/components/RecommendationPanel'), {
  ssr: false,
  loading: () => <Skeleton />,
});
```

### Server-Side Performance (HIGH)
- Authenticate Server Actions the same way as API routes (verify session server-side).
- Use `React.cache()` for DB/API calls shared across a request tree.
- Minimize data serialized from Server → Client components (pass only needed props).
- Use `after()` (Next.js 15) or `waitUntil` for non-blocking post-response work (analytics, cache warm).

### Re-render Optimization (MEDIUM)
- Hoist static non-primitive props (objects/arrays) outside components or use `useMemo`.
- Use primitive values (not objects) as `useEffect`/`useMemo` dependencies.
- Use `useRef` for transient frequently-updated values (scroll position, animation frame).
- Use `startTransition` for non-urgent state updates (search filtering, tab switches).
- Derive state during render rather than syncing with `useEffect`.
- Use `useTransition` for loading states over `useState` boolean flags.

### Rendering (MEDIUM)
- Use `content-visibility: auto` for long media lists (CSS, not JS).
- Extract static JSX (icons, empty states) outside component functions.
- Use ternary (`condition ? <A/> : <B/>`) not `&&` for conditional rendering.
- Suppress expected hydration mismatches with `suppressHydrationWarning` (e.g., timestamps).

### JavaScript Performance (LOW-MEDIUM)
- Build `Map` indexes for repeated lookups (genre → items, creator → items).
- Use `Set` for O(1) membership checks (owned IDs, wishlist IDs).
- Hoist `RegExp` outside loops for search filtering.
- Use `toSorted()` / `toReversed()` for immutable array operations.
- Early-return from functions to reduce nesting.

---

## Feature Roadmap

### Core (MVP)
- [ ] User auth (JWT, OAuth via Google/GitHub)
- [ ] User profiles (public/private, shareable URL)
- [ ] Media CRUD: movies, music, games
  - Fields: title, creator, genre, release date, cover art, notes, rating
- [ ] Status tracking: `owned` | `wishlist` | `currently_using` | `completed`
- [ ] Full-text search by title, creator, genre

### AI Features
- [ ] **Auto-metadata enrichment** — fetch title/creator/genre/cover from TMDB/MusicBrainz/IGDB via AI-assisted parsing
- [ ] **Smart recommendations** — Claude API suggests similar items based on collection
- [ ] **Mood-based discovery** — "I want something chill tonight" → filtered suggestions
- [ ] **Natural language search** — "action movies from the 90s I haven't finished"
- [ ] **Collection insights** — AI summary: top genres, completion rate, listening/watching patterns
- [ ] **Duplicate detection** — AI flags likely duplicates before insert

### Extra Features
- [ ] Collection sharing — public profile page with filterable grid
- [ ] Import from CSV / Letterboxd / Goodreads / Steam
- [ ] Collaborative wishlists (share with friends)
- [ ] Activity feed (recently added, status changes)
- [ ] Stats dashboard (breakdown by type, status, genre, year)
- [ ] Dark mode (CSS variables, no flash on load via inline script)

---

## AI Integration (Claude API)

Use `claude-sonnet-4-6` (latest capable model) for all AI features.

```go
// backend/internal/ai/client.go
// Use context, handle errors, stream where latency matters
func (c *Client) Recommend(ctx context.Context, collection []media.Item) ([]Recommendation, error) {
    // Build prompt from collection summary
    // Call Anthropic API with context propagation
    // Return typed recommendations
}
```

- Stream AI responses to the frontend via Server-Sent Events (SSE) for long operations.
- Cache AI results per user+collection-hash to avoid redundant API calls (LRU cache in Go).
- Never expose the Anthropic API key to the client — all AI calls go through the Go backend.

---

## Environment Variables

```
# Backend
DATABASE_URL=
JWT_SECRET=
ANTHROPIC_API_KEY=
TMDB_API_KEY=
IGDB_CLIENT_ID=
IGDB_CLIENT_SECRET=

# Frontend
NEXT_PUBLIC_API_URL=
```

---

## Development Commands

```bash
# Backend
cd backend && go run ./cmd/server       # dev server
cd backend && go test ./... -race       # tests with race detector
cd backend && golangci-lint run         # lint

# Frontend
cd frontend && npm run dev              # dev server
cd frontend && npm run build            # production build
cd frontend && npm run lint             # ESLint
```

---

## Code Quality Gates

- Go: `golangci-lint` must pass; `go test ./... -race` must pass; 80%+ coverage on domain packages.
- Next.js: `next build` must succeed with no type errors; bundle analyzer checked for regressions.
- No secrets committed; use `.env.local` (gitignored).
- All AI prompts versioned in `backend/internal/ai/prompts/`.
