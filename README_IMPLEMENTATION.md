# ERH Safety System - Implementation Guide

This document describes the implementation of the ERH Crowd-in-the-Loop Safety Decision System based on the technical implementation plan.

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── model/                   # Data models (GORM)
│   ├── dto/                     # Data Transfer Objects (request/response)
│   ├── vo/                      # Value Objects (response)
│   ├── handler/                 # HTTP handlers (Gin)
│   ├── service/                 # Business logic services
│   ├── aggregation/             # Signal aggregation engine
│   ├── middleware/              # HTTP middleware (rate limiting, auth)
│   ├── database/                # Database initialization
│   └── redis/                   # Redis client
├── database/
│   └── migrations/              # Database migration files
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Phase 1 Implementation Status

### ✅ Completed

1. **Project Initialization**
   - Go module setup (`go.mod`)
   - Project structure
   - Configuration management
   - Docker setup

2. **Database Schema**
   - Initial migration file (`001_initial_schema.up.sql`)
   - All core tables defined (signals, aggregated_summaries, decision_states, etc.)
   - Indexes and constraints

3. **Data Models**
   - `Signal` model with JSONB support
   - `AggregatedSummary` model
   - Custom types for PostgreSQL arrays and JSONB

4. **Signal Reception Layer**
   - Four signal handlers implemented:
     - `CrowdHandler` - Route 2 App reports
     - `StaffHandler` - Staff reports
     - `InfrastructureHandler` - Infrastructure signals
     - `EmergencyHandler` - Emergency calls
   - Service layer for signal creation
   - DTOs for request validation

5. **Rate Limiting**
   - Redis-based rate limiter
   - Middleware for Gin
   - Configurable limits per action type

6. **Signal Aggregation Engine**
   - Time window aggregation
   - Weighted aggregation by zone
   - Effective signal filtering (quality + trust score)
   - Outlier detection (Z-score method)

### 🚧 In Progress / TODO

1. **Trust Scoring Engine**
   - Implement trust score calculation
   - Device integrity checks
   - Historical accuracy tracking

2. **Testing**
   - Unit tests for services
   - Integration tests for handlers
   - Aggregation engine tests

3. **Additional Features**
   - Authentication/Authorization middleware
   - Device ID extraction from tokens
   - Staff ID extraction from tokens
   - Error handling improvements
   - Logging improvements

## Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL 15+ (via Docker)
- Redis 7+ (via Docker)

### Running the Application

1. **Start dependencies:**
   ```bash
   docker-compose up -d postgres redis
   ```

2. **Run migrations:**
   ```bash
   # Connect to PostgreSQL and run migrations manually, or use a migration tool
   psql -h localhost -U postgres -d erh_safety -f database/migrations/001_initial_schema.up.sql
   ```

3. **Run the application:**
   ```bash
   go run cmd/server/main.go
   ```

   Or use Docker Compose:
   ```bash
   docker-compose up api
   ```

### API Endpoints

- `GET /health` - Health check
- `POST /api/v1/reports` - Submit crowd report (Rate limited: 3/hour)
- `POST /api/v1/staff/reports` - Submit staff report
- `POST /api/v1/infrastructure/signals` - Submit infrastructure signal
- `POST /api/v1/emergency/calls` - Submit emergency call

## Configuration

Configuration is loaded from environment variables. See `internal/config/config.go` for all available options.

Default values:
- Server port: `8080`
- Database: `localhost:5432/erh_safety`
- Redis: `localhost:6379`

## Next Steps

See the technical implementation plan for Phase 2 and beyond:
- Trust scoring system
- Decision state machine
- ERH complexity calculation
- High-impact action gates
- CAP message engine
- Route 1/Route 2 adapters

