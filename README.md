# Personal Media Collection Tracker

A full-stack app to track your movies, music, and games — with AI-powered recommendations, mood discovery, and natural language search.

**Stack:** Go 1.21 + chi · Next.js 14 App Router · PostgreSQL · Claude AI

---

## Features

- **Media CRUD** — movies, music, games with cover art, ratings, notes
- **Status tracking** — owned / wishlist / in-progress / completed
- **Full-text search** — PostgreSQL tsvector + trigram indexes
- **Metadata enrichment** — auto-fetch from TMDB, MusicBrainz, IGDB
- **AI recommendations** — Claude suggests similar items based on your collection
- **Mood discovery** — "I want something chill tonight" → personalized suggestions
- **AI insights** — streaming collection analysis via SSE
- **Natural language search** — parse free-text queries into structured filters
- **Public profiles** — shareable collection pages

---

## Quick Start (Docker)

```bash
cp backend/.env.example backend/.env
# Edit backend/.env and set ANTHROPIC_API_KEY, TMDB_API_KEY, etc.

docker compose up --build
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

---

## Local Development

### Backend

```bash
cd backend
cp .env.example .env
# Edit .env with your credentials

# Start PostgreSQL (or use Docker)
docker run -d --name ems-pg -e POSTGRES_USER=ems -e POSTGRES_PASSWORD=ems -e POSTGRES_DB=ems -p 5432:5432 postgres:16-alpine

go run ./cmd/server
```

### Frontend

```bash
cd frontend
cp .env.local.example .env.local
# Edit .env.local: NEXT_PUBLIC_API_URL=http://localhost:8080

npm install
npm run dev
```

---

## Environment Variables

### Backend (`backend/.env`)

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | ✅ | PostgreSQL connection string |
| `JWT_SECRET` | ✅ | Secret for JWT signing (min 32 chars) |
| `ANTHROPIC_API_KEY` | ✅ | Anthropic API key for AI features |
| `TMDB_API_KEY` | ☐ | The Movie Database API key |
| `IGDB_CLIENT_ID` | ☐ | Twitch/IGDB client ID |
| `IGDB_CLIENT_SECRET` | ☐ | Twitch/IGDB client secret |
| `FRONTEND_URL` | ☐ | Frontend URL for CORS (default: http://localhost:3000) |
| `PORT` | ☐ | Server port (default: 8080) |

### Frontend (`frontend/.env.local`)

| Variable | Required | Description |
|----------|----------|-------------|
| `NEXT_PUBLIC_API_URL` | ✅ | Backend API URL |

---

## API Routes

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/register` | Register |
| POST | `/api/auth/login` | Login (returns JWT) |
| GET | `/api/auth/me` | Get current user |
| GET | `/api/media` | List media (paginated, filterable) |
| POST | `/api/media` | Create media item |
| GET | `/api/media/:id` | Get media item |
| PUT | `/api/media/:id` | Update media item |
| DELETE | `/api/media/:id` | Delete media item |
| PATCH | `/api/media/:id/status` | Update status |
| GET | `/api/search?q=` | Full-text search |
| POST | `/api/metadata/search` | External metadata lookup |
| GET | `/api/ai/recommendations` | AI recommendations |
| GET | `/api/ai/insights` | Streaming AI insights (SSE) |
| POST | `/api/ai/nl-search` | Natural language → filters |
| POST | `/api/ai/mood` | Mood-based discovery |
| POST | `/api/ai/duplicates` | Duplicate detection |
| GET | `/api/profile/:username` | Public profile |
| PUT | `/api/profile` | Update profile |
| GET | `/api/activity` | Activity feed |

---

## Deployment

### Backend → Railway

1. Push to GitHub
2. New Railway project → Deploy from GitHub
3. Add PostgreSQL service (Railway provides `DATABASE_URL`)
4. Set env vars: `JWT_SECRET`, `ANTHROPIC_API_KEY`, etc.
5. Set `FRONTEND_URL` to your Vercel URL

### Frontend → Vercel

1. Import repo on Vercel, set root to `frontend/`
2. Set `NEXT_PUBLIC_API_URL` to your Railway backend URL
3. Deploy

---

## Project Structure

```
ems/
├── backend/
│   ├── cmd/server/main.go        # Entry point
│   └── internal/
│       ├── auth/                 # JWT auth
│       ├── media/                # Media CRUD
│       ├── ai/                   # AI features (Anthropic)
│       ├── metadata/             # TMDB/MusicBrainz/IGDB
│       ├── search/               # Full-text search
│       ├── profile/              # User profiles
│       ├── activity/             # Activity feed
│       ├── db/                   # PostgreSQL + migrations
│       └── httputil/             # HTTP helpers
└── frontend/
    ├── app/                      # Next.js App Router pages
    ├── components/               # React components
    └── lib/                      # API clients, hooks, utils
```
