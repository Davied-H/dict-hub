# CLAUDE.md - Dict-Hub AI Assistant Guide

This file provides guidance for AI assistants working with the Dict-Hub codebase.

## Project Overview

Dict-Hub is a modern MDX dictionary query platform with:
- **Go backend** (Gin framework) serving API and embedded frontend
- **React frontend** (TypeScript + Vite) for the web UI
- **SQLite database** for persistence (history, word frequency, dictionary metadata)
- **MDX format support** for dictionary files (MDict format)

The application is designed as a self-hosted dictionary service supporting multiple dictionaries, search history, word frequency tracking, and audio pronunciation.

## Tech Stack

### Backend (Go 1.24)
- **Framework:** Gin web framework
- **ORM:** GORM with SQLite driver
- **Config:** Viper (YAML + environment variables)
- **Entry point:** `backend/cmd/server/main.go`

### Frontend (React 19 + TypeScript)
- **Build tool:** Vite 5.4
- **State management:** Zustand
- **Data fetching:** TanStack Query (React Query)
- **UI components:** HeroUI
- **Styling:** Tailwind CSS 3.4
- **Animations:** Framer Motion

### Official Documentation Site
- Located in `official-web/`
- Uses MDX (Markdown + React) for documentation
- Deployed to GitHub Pages at `/dict-hub/`

## Repository Structure

```
dict-hub/
├── backend/                    # Go backend service
│   ├── cmd/server/main.go     # Application entry point
│   ├── configs/config.yaml    # Default configuration
│   ├── internal/              # Core business logic
│   │   ├── handler/           # HTTP request handlers
│   │   ├── service/           # Business logic (mdx/, audio/, etc.)
│   │   ├── model/             # Database models
│   │   ├── database/          # DB initialization
│   │   ├── config/            # Config management
│   │   ├── router/            # Route definitions
│   │   ├── middleware/        # CORS, recovery, logging
│   │   └── cache/             # In-memory caching
│   ├── pkg/                   # Reusable packages
│   │   ├── mdict/             # MDX format handling
│   │   └── response/          # Response utilities
│   ├── thirdparty/mdx/        # MDX parser implementation
│   └── web/dist/              # Embedded frontend (build output)
│
├── frontend/                  # React application
│   ├── src/
│   │   ├── api/              # API client methods
│   │   ├── components/       # Reusable UI components
│   │   ├── pages/            # Page components (Home, Settings, History, Word)
│   │   ├── stores/           # Zustand stores (theme, search)
│   │   ├── hooks/            # Custom hooks (useSearch, useSuggest, etc.)
│   │   └── types/            # TypeScript definitions
│   ├── vite.config.ts        # Builds to ../backend/web/dist
│   └── package.json
│
├── official-web/             # Documentation site (MDX)
│   └── src/docs/             # MDX documentation files
│
├── scripts/                  # Alfred & Raycast integrations
├── docker-compose.yml        # Docker Compose config
├── Dockerfile               # Multi-stage Docker build
├── Makefile                 # Build commands
└── start.sh                 # One-command startup script
```

## Development Workflow

### Quick Start (Development)

```bash
# One-command startup (recommended)
./start.sh start

# Or manually in separate terminals:
# Terminal 1 - Backend (port 8080)
cd backend && go run cmd/server/main.go

# Terminal 2 - Frontend (port 3000, proxies /api to 8080)
cd frontend && npm install && npm run dev
```

### Build Commands (Makefile)

```bash
make build          # Full build (frontend + backend)
make frontend       # Build frontend only
make backend        # Build backend only
make clean          # Clean build artifacts
make docker-build   # Build Docker image
make release VERSION=x.x.x  # Build and push release
```

### Frontend Development

```bash
cd frontend
npm run dev         # Dev server at localhost:3000
npm run build       # Production build (outputs to backend/web/dist)
npm run lint        # ESLint check
```

### Backend Development

```bash
cd backend
go run cmd/server/main.go    # Run development server
go build -o dict-hub ./cmd/server  # Build binary
```

## Key Architecture Patterns

### Backend Layering
```
HTTP Request → Router → Handler → Service → Model/Database
```
- **Handlers** (`internal/handler/`): Parse requests, call services, return responses
- **Services** (`internal/service/`): Business logic, dictionary operations
- **Models** (`internal/model/`): GORM database models

### Frontend Data Flow
```
Component → Custom Hook → React Query → Axios → API
```
- **Stores** (Zustand): Global state (theme, selected dictionaries)
- **React Query**: Server state caching and synchronization
- **Custom hooks**: Encapsulate API calls with loading/error states

### Embedded Frontend
The frontend builds to `backend/web/dist/` and is embedded into the Go binary using `go:embed`. In production, a single binary serves both API and static files.

## API Endpoints

Base URL: `http://localhost:8080`

### Core Search
- `GET /api/v1/search` - Search word across dictionaries
- `GET /api/v1/search/suggest` - Get word suggestions

### Dictionary Management
- `GET /api/v1/dictionaries` - List all dictionaries
- `POST /api/v1/dictionaries` - Add new dictionary
- `PUT /api/v1/dictionaries/:id/toggle` - Enable/disable dictionary
- `PUT /api/v1/dictionaries/reorder` - Reorder dictionaries
- `DELETE /api/v1/dictionaries/:id` - Delete dictionary

### History & Stats
- `GET /api/v1/history` - Get search history
- `DELETE /api/v1/history` - Clear history
- `POST /api/v1/wordfreq/import` - Import word frequency data

### Audio
- `GET /api/v1/audio/:word` - Get pronunciation audio
- `GET /api/v1/audio/:word/availability` - Check audio availability

### Health
- `GET /health` - Health check endpoint

## Configuration

### Backend Config (`backend/configs/config.yaml`)

```yaml
server:
  port: 8080
  mode: development  # or production

database:
  driver: sqlite
  path: ./data/dictionary.db

mdx:
  dict_dir: ./dicts
  source_dir: ./dicts/source
  sound_dir: ./dicts/sound
  auto_load: true
```

### Environment Variables (override config.yaml)

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | 8080 | Server port |
| `SERVER_MODE` | development | Mode (development/production) |
| `DATABASE_PATH` | ./data/dictionary.db | SQLite database path |
| `MDX_DICT_DIR` | ./dicts | Dictionary root directory |
| `MDX_SOURCE_DIR` | ./dicts/source | MDX source files |
| `MDX_SOUND_DIR` | ./dicts/sound | Audio files |
| `MDX_AUTO_LOAD` | false | Auto-load dictionaries on startup |

## Code Conventions

### Go (Backend)
- Use layered architecture (handler → service → model)
- Explicit error returns (no panic for business logic)
- Configuration via Viper with env var overrides
- Use GORM for database operations

### TypeScript (Frontend)
- Functional components with hooks
- Zustand for global state management
- TanStack Query for server state
- Path alias: `@/` maps to `src/`
- Tailwind CSS for styling (utility-first)

### File Naming
- Go: `snake_case.go`
- TypeScript: `PascalCase.tsx` for components, `camelCase.ts` for utilities
- CSS: Use Tailwind classes, avoid custom CSS files

## Docker Deployment

### Production Build (Multi-stage)
1. **Stage 1:** Node 20 - Build React frontend
2. **Stage 2:** Go 1.24 - Compile backend with embedded frontend
3. **Stage 3:** Alpine 3.20 - Minimal runtime image

```bash
# Build and run with Docker Compose
docker-compose up -d

# Or manually
docker build -t dict-hub .
docker run -d -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/dicts:/app/dicts \
  dict-hub
```

### Data Persistence
Mount these volumes for data persistence:
- `./data` - SQLite database
- `./dicts` - Dictionary files (MDX/MDD)
- `./dicts/source` - MDX source files
- `./dicts/sound` - Audio pronunciation files

## Common Tasks

### Adding a New API Endpoint
1. Define route in `backend/internal/router/router.go`
2. Create handler in `backend/internal/handler/`
3. Implement service logic in `backend/internal/service/`
4. Add frontend API client in `frontend/src/api/`

### Adding a New Page
1. Create page component in `frontend/src/pages/`
2. Add route in `frontend/src/App.tsx`
3. Create any needed hooks in `frontend/src/hooks/`

### Adding Documentation
1. Create MDX file in `official-web/src/docs/`
2. Add navigation link in documentation layout

## Important Files

| File | Purpose |
|------|---------|
| `backend/cmd/server/main.go` | Application entry point |
| `backend/internal/router/router.go` | API route definitions |
| `backend/configs/config.yaml` | Default configuration |
| `frontend/src/App.tsx` | React app with routing |
| `frontend/vite.config.ts` | Vite config (build output, proxy) |
| `docker-compose.yml` | Docker deployment config |
| `Makefile` | Build automation |
| `start.sh` | Development startup script |

## Testing

Currently no automated tests are configured. When adding tests:
- **Go:** Use standard `testing` package, place tests in `*_test.go` files
- **React:** Consider Vitest (Vite-native) or Jest with React Testing Library

## Git Workflow

- Main branch: `main`
- Feature branches: `feature/description`
- Use conventional commit messages
- CI/CD via GitHub Actions (`.github/workflows/`)

## Notes for AI Assistants

1. **Understand the full stack**: Changes often span backend and frontend
2. **Check existing patterns**: Follow established code style in each directory
3. **Build process matters**: Frontend builds into backend for embedding
4. **Configuration layers**: Environment vars override YAML config
5. **Chinese comments**: Some code/docs contain Chinese - this is normal
6. **MDX format**: Core feature - understand dictionary file handling in `backend/thirdparty/mdx/`
