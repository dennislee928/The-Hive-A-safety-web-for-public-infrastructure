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
│   │   ├── signal.go
│   │   ├── aggregated_summary.go
│   │   └── device_trust.go      # Phase 2: Device trust models
│   ├── dto/                     # Data Transfer Objects (request/response)
│   ├── vo/                      # Value Objects (response)
│   │   ├── signal_response.go
│   │   ├── decision_response.go # Phase 2: Decision responses
│   │   └── erh_response.go      # Phase 2: ERH responses
│   ├── handler/                 # HTTP handlers (Gin)
│   │   ├── crowd_handler.go
│   │   ├── staff_handler.go
│   │   ├── infrastructure_handler.go
│   │   ├── emergency_handler.go
│   │   ├── operator_handler.go  # Phase 2: Operator endpoints
│   │   └── dashboard_handler.go # Phase 2: Dashboard API
│   ├── service/                 # Business logic services
│   │   └── signal_service.go
│   ├── aggregation/             # Signal aggregation engine
│   ├── trust/                   # Phase 2: Trust scoring engine
│   │   └── scorer.go
│   ├── decision/                # Phase 2: Decision engine
│   │   ├── state_machine.go
│   │   ├── evaluator.go
│   │   ├── service.go
│   │   └── errors.go
│   ├── erh/                     # Phase 2: ERH governance
│   │   ├── complexity.go
│   │   ├── ethical_prime.go
│   │   └── breakpoint_detector.go
│   ├── middleware/              # HTTP middleware (rate limiting, auth)
│   ├── database/                # Database initialization
│   └── redis/                   # Redis client
├── database/
│   └── migrations/              # Database migration files
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## Implementation Status

### ✅ Phase 1: Completed

1. **Project Initialization**
   - Go module setup
   - Project structure
   - Configuration management
   - Docker setup

2. **Database Schema**
   - Initial migration file
   - All core tables defined

3. **Signal Reception Layer**
   - Four signal handlers implemented
   - Service layer for signal creation
   - DTOs for request validation

4. **Rate Limiting**
   - Redis-based rate limiter
   - Middleware for Gin

5. **Signal Aggregation Engine**
   - Time window aggregation
   - Weighted aggregation by zone
   - Effective signal filtering

### ✅ Phase 2: Completed

1. **Trust Scoring Engine**
   - Complete trust score calculation
   - Historical accuracy tracking
   - Frequency scoring
   - Device integrity checks (framework)
   - Cross-source corroboration

2. **Decision State Machine**
   - State definitions (D0-D6)
   - State transition logic
   - Decision service layer

3. **Decision Evaluator**
   - Decision evaluation logic
   - Corroboration checking
   - Target state determination

4. **ERH Complexity Calculator**
   - x_s, x_d, x_c calculation
   - x_total calculation
   - Complexity level determination

5. **Ethical Prime Calculator**
   - FN-prime, FP-prime, Bias-prime, Integrity-prime
   - Framework for calculation

6. **Operator & Dashboard APIs**
   - Operator endpoints for decision management
   - Dashboard API for monitoring

## Getting Started

### Prerequisites

- Go 1.23+
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

#### Signal Endpoints
- `POST /api/v1/reports` - Submit crowd report (Rate limited: 3/hour)
- `POST /api/v1/staff/reports` - Submit staff report
- `POST /api/v1/infrastructure/signals` - Submit infrastructure signal
- `POST /api/v1/emergency/calls` - Submit emergency call

#### Operator Endpoints (Phase 2)
- `POST /api/v1/operator/decisions/:zone_id/d0` - Create D0 Pre-Alert
- `POST /api/v1/operator/decisions/:decision_id/transition` - Transition decision state
- `GET /api/v1/operator/zones/:zone_id/state` - Get latest decision state

#### Dashboard Endpoints (Phase 2)
- `GET /api/v1/dashboard/zones/:zone_id` - Get dashboard data (state, complexity, ethical primes)

## Configuration

Configuration is loaded from environment variables. See `internal/config/config.go` for all available options.

Default values:
- Server port: `8080`
- Database: `localhost:5432/erh_safety`
- Redis: `localhost:6379`

## Testing

Run all tests:
```bash
go test ./...
```

Run tests for specific packages:
```bash
go test ./internal/trust/...
go test ./internal/decision/...
go test ./internal/erh/...
```

## Next Steps

See the technical implementation plan for Phase 3 and beyond:
- High-impact action gates (dual control, keepalive, TTL)
- CAP message engine
- Route 1/Route 2 adapters
- Audit & sealing system
