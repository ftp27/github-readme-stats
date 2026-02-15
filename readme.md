# Go GitHub Readme Stats

This is a Go reimplementation of the [GitHub Readme Stats](https://github.com/anuraghazra/github-readme-stats).

# Motivation

The original project is not working well lately and selfhosted solutions stopped to me as well. So I decided to vibecode the original project in Go to have a more stable solution with lower memory footprint. I hope I find time to reimplement most of the code myself, but it works now and hope it will be useful for some people.

## Requirements

- Go 1.25+
- A GitHub token in `PAT_1` (required)

## Quick start

```bash
cd app
go run ./cmd/server
```

Server starts on port `9000` by default.

## Configuration

The server reads config via Viper from the following sources (in order):

1. Environment variables
2. `.env` in `app/` (optional)
3. `config.yaml`, `config.yml`, or `config.json` in `app/` (optional)

Example `.env` file:

```bash
PORT=9000
NODE_ENV=production
PAT_1=your_github_token_here
CACHE_SECONDS=
FETCH_MULTI_PAGE_STARS=false
WHITELIST=
GIST_WHITELIST=
EXCLUDE_REPO=
```

## API

Base path: `/api`

- `/api?username=...`
- `/api/pin?username=...&repo=...`
- `/api/top-langs?username=...`
- `/api/wakatime?username=...`
- `/api/gist?id=...`
- `/api/status/up`
- `/api/status/pat-info`

## Docker

Build:

```bash
docker build -t github-readme-stats ./app
```

Run:

```bash
docker run --rm -p 9000:9000 \
  -e PAT_1=your_github_token_here \
  github-readme-stats
```

## Notes

- Only the stats card (`/api`) is fully implemented in Go right now.
- Other card endpoints return a placeholder SVG until they are ported.
