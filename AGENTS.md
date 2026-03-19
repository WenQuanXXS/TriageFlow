# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

## Project Overview

TriageFlow is an LLM-assisted outpatient triage and intelligent queue management system. It is **not** a diagnostic system — it assists with pre-screening and queue routing only. High-risk scenarios must always go through rule engine safety checks, never relying solely on LLM output.

Three user roles: Patient (submits chief complaint), Nurse/Triage desk (reviews AI suggestions, overrides if needed), Doctor (views queue, handles consultations).

## Tech Stack

- **Frontend**: React, React Router, Ant Design, Axios, Zustand or Redux Toolkit, WebSocket
- **Backend**: Go, Gin, GORM, MySQL, WebSocket
- **LLM**: Eino (workflow orchestration), OpenAI-compatible model API
- **Infra**: Docker / Docker Compose, Swagger (API docs), zap/logrus (logging)

## Architecture

The system has three backend service layers behind a Gin API gateway:

1. **Triage Service** — Calls Eino LLM orchestrator to extract structured symptoms, risk signals, candidate departments, and suggested priority from patient chief complaint. A **Rule Engine** then applies safety overrides (high-risk keyword matching, final priority adjudication). All decisions produce audit logs.
2. **Queue Service** — Manages waiting order and real-time call numbers. Pushes updates via WebSocket to all frontends.
3. **Basic Service** — Handles user/patient CRUD and foundational data.

Data flow: Patient input → LLM structured extraction → Rule engine safety check → Queue assignment → WebSocket notification.

## Key Design Constraints

- LLM handles natural language understanding; the rule engine makes final safety decisions — never skip the rule engine for high-risk symptoms
- All triage results must be auditable and traceable
- The project prioritizes a working, demonstrable end-to-end demo with clear AI integration

## Build / Test / Run Commands

### Backend
```bash
cd backend
go mod tidy          # install dependencies
go build ./...       # compile
go run main.go       # start server on :8080
go test ./... -v     # run tests (requires MySQL with triageflow_test database)
```

### Frontend
```bash
cd frontend
npm install          # install dependencies
npm run dev          # start dev server on :5173
npm test             # run tests (vitest)
```

### Environment Variables (backend)
- `DB_HOST` (default: `127.0.0.1`)
- `DB_PORT` (default: `3306`)
- `DB_USER` (default: `root`)
- `DB_PASS` (default: `1234`)
- `DB_NAME` (default: `triageflow`)
