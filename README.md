# URL Shortener

A fast URL shortening service built with Go and Redis. Part of [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-url-shortener).

## What It Does

Converts long URLs into short, memorable links:
- `https://example.com/very/long/url` → `http://localhost:8080/abc123`
- Fast Redis-powered redirects
- Automatic collision handling

## Quick Start

```bash
# Start everything with Docker
docker-compose up --build

# Service runs on http://localhost:8080
```

## Usage

**Shorten a URL:**
```bash
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}'

# Response: {"key":"xY9pQ2", "short_url":"http://localhost:8080/xY9pQ2"}
```

**Use the short URL:**
```bash
curl -L http://localhost:8080/xY9pQ2
# Redirects to https://github.com
```

## Configuration

Edit `.env` to customize:
```bash
SHORT_URL_BASE=https://yourdomain.com  # Your domain
APP_PORT=8080                           # Server port
REDIS_PASSWORD=                         # Redis password
```

## How It Works

1. Generates hash using xxHash algorithm
2. Encodes to Base62 (0-9, A-Z, a-z) for short codes
3. Stores in Redis for instant lookups
4. Handles collisions with random suffixes

## Troubleshooting

```bash
# View logs
docker-compose logs

# Restart services
docker-compose restart

# Clear all data
docker-compose down -v
```

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/shorten` | Create short URL |
| GET | `/{key}` | Redirect to original URL |

## Tech Stack

- **Go** - Fast, concurrent backend
- **Redis** - In-memory data store
- **Docker** - Easy deployment