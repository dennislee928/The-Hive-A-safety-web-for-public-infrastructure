# ERH Crowd-in-the-Loop Safety Decision PoC

This repository contains a **complete Proof of Concept (PoC)** implementation for a **crowd-in-the-loop safety decision system** across four high-density zones:

1. **Station Interior** (concourse / gates / platforms)
2. **Train Car** (inside carriage)
3. **Station Perimeter** (station surroundings, entrances/exits, transfer corridors)
4. **Other High-Density Areas** (events, plazas, festivals, dispersal flows)

The PoC implements two public-communication routes that coexist:

- **Route 1 (Baseline / No Download):** Standards-based public warning delivery (e.g., cell broadcast / location-based SMS) with a CAP-centered alert format.
- **Route 2 (App / Optional Download):** A public app that enables bidirectional interactions (structured crowd reports, personalized guidance, check-in), while adding abuse and privacy controls.

A core constraint is **ERH governance**: as system complexity increases, critical misjudgments ("ethical primes") must remain **measurably bounded** rather than exploding with scale.

## Architecture

The system consists of:

- **Backend API** (Go/Gin): Core safety decision system with ERH governance
- **Frontend Web** (Next.js): Client-facing interface and admin dashboard
- **Mobile App** (Flutter): Route 2 mobile application
- **Database** (PostgreSQL/TimescaleDB): Data persistence
- **Cache** (Redis): Rate limiting and caching

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development)
- Node.js 20+ (for frontend development)
- Flutter 3.22+ (for mobile app development)

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Taiwanese-Cultural-Heritage-based-Physical-Post-Quantom-Encryption-Method
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Access the services**
   - Frontend: http://localhost:3000
   - API: http://localhost:8080
   - API Health: http://localhost:8080/health

### Local Development

#### Backend

```bash
# Install dependencies
go mod download

# Set up environment variables
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/erh_safety?sslmode=disable
export REDIS_URL=redis://localhost:6379

# Run database migrations
# (using your preferred migration tool)

# Start the server
go run cmd/server/main.go
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Set environment variables
export NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# Start development server
npm run dev
```

#### Mobile App

```bash
cd mobile_app

# Install dependencies
flutter pub get

# Run on device/emulator
flutter run
```

## Project Structure

```
.
├── cmd/server/          # Backend server entry point
├── internal/            # Backend implementation
│   ├── handler/        # HTTP handlers
│   ├── service/        # Business logic
│   ├── model/          # Data models
│   ├── decision/       # Decision state machine
│   ├── erh/            # ERH governance
│   ├── cap/            # CAP message generation
│   ├── route1/         # Route 1 adapters
│   └── route2/         # Route 2 services
├── frontend/            # Next.js web application
│   ├── src/app/        # Pages (client & admin)
│   └── src/lib/        # API services
├── mobile_app/          # Flutter mobile application
├── database/            # Database migrations
├── tests/               # Robot Framework tests
└── docs/                # Documentation
```

## Features

### Backend
- ✅ Signal aggregation and trust scoring
- ✅ Decision state machine (D0-D6)
- ✅ ERH governance and ethical primes
- ✅ CAP message generation
- ✅ Route 1 adapters (Cell Broadcast, SMS, etc.)
- ✅ Route 2 API (guidance, assistance, feedback)
- ✅ Audit logging and evidence archiving
- ✅ High-impact action gating

### Frontend
- ✅ Client interface for CAP message viewing
- ✅ Admin dashboard with real-time monitoring
- ✅ Decision management interface
- ✅ ERH monitoring and metrics
- ✅ Audit log viewer
- ✅ CAP message management

### Mobile App
- ✅ Device registration and authentication
- ✅ Structured crowd reporting
- ✅ Personalized guidance
- ✅ Assistance requests
- ✅ Push notifications

## API Documentation

The API follows RESTful conventions and is documented through:
- OpenAPI/Swagger (if configured)
- Handler code with request/response structures
- See `internal/handler/` for endpoint implementations

Main API groups:
- `/api/v1/dashboard` - Dashboard data
- `/api/v1/operator` - Operator decision endpoints
- `/api/v1/cap` - CAP message endpoints
- `/api/v1/route2` - Route 2 App endpoints
- `/api/v1/erh` - ERH governance endpoints
- `/api/v1/audit` - Audit and evidence endpoints

## Testing

### Backend Tests

```bash
go test ./...
```

### Robot Framework Tests

```bash
cd tests
pip install -r requirements.txt
./run_tests.sh
```

## CI/CD

The project includes GitHub Actions workflows for:
- Go tests and linting
- Docker image building
- Flutter app building (Android & iOS)
- Robot Framework tests

See `.github/workflows/` for details.

## Documentation

- [`docs/`](./docs/) - Complete system documentation
- [`plan.md`](./plan.md) - PoC plan and milestones
- [`agent.md`](./agent.md) - Agent roles and constraints
- [`structure.md`](./structure.md) - Repository structure

## License

See [LICENSE.txt](./LICENSE.txt) for license information.

## Contributing

This is a PoC implementation. For production use, please refer to the documentation and ensure all security, privacy, and compliance requirements are met.

