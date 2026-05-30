# algo*roulette*

A full-stack algorithm practice tool that randomizes your grind. No more cherry-picking easy problems — spin the wheel and face whatever comes up.


**Live at:** [algoroulette.up.railway.app](https://algoroulette.up.railway.app)

---
 
## What it does
 
algoroulette pulls a random algorithm problem from a curated question bank and throws it at you. Filter by topic or difficulty if you dare, or go full random and suffer accordingly.
 
Once you get a question you can mark it solved, skip it, bookmark it for later, and add personal notes. Your progress is tracked per user so every spin builds a history of what you've done and what you've been avoiding.
 
---
 
## Features
 
- **Random question engine** — spin to get a random problem, filtered by topic and difficulty
- **150 NeetCode problems** pre-loaded, importable via JSON
- **Per-user progress tracking** — solved, skipped, bookmarked, notes
- **JWT authentication** — register and login, progress tied to your account
- **Bulk import** — upload a JSON file to seed questions in one shot
- **LeetCode link** — every question links directly to its LeetCode page
- **Auto-running migrations** — schema changes apply on server startup
---
 
## Tech Stack
 
**Backend**
- Go + Gin framework
- PostgreSQL
- JWT auth with bcrypt password hashing
- Repository pattern with `internal/` directory structure
- Rate limiting on auth endpoints
**Frontend**
- Vanilla JS, no framework
- Cormorant Garamond + DM Mono typography
- Dark theme, minimal design
**Deployment**
- Backend + DB on Railway
- Static frontend served by the Go server
---
 
## API Endpoints
 
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/auth/register` | — | Create an account |
| POST | `/auth/login` | — | Login and receive JWT |
| GET | `/question` | ✓ | Get a random question |
| POST | `/question` | ✓ | Add a single question |
| POST | `/questions/import` | ✓ | Bulk import from JSON |
| DELETE | `/questions` | ✓ | Clear all questions |
| PATCH | `/question/:id/status` | ✓ | Mark solved or skipped |
| PATCH | `/question/:id/bookmark` | ✓ | Toggle bookmark |
| PATCH | `/question/:id/notes` | ✓ | Save personal notes |
 
---
