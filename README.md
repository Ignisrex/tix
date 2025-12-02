# Tix - Event Ticket Booking System

A distributed microservices-based ticket booking platform built as a take-home assignment. The system enables users to search for events, reserve tickets, and complete purchases with real-time availability tracking.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Scaling Considerations](#scaling-considerations)
- [Future Improvements](#future-improvements)
- [Tradeoffs & Design Decisions](#tradeoffs--design-decisions)

## Overview

Tix is a ticket booking system that demonstrates a microservices architecture with the following core capabilities:

- **Event Search**: Full-text search across events using Elasticsearch
- **Ticket Reservation**: Atomic reservation system using Redis for distributed locking
- **Ticket Purchase**: Transactional purchase flow with payment processing
- **Real-time Availability**: Seat selection with live ticket status updates
- **Reservation Timer**: Client-side countdown for active reservations (180 seconds TTL)

The system is designed with scalability and high availability in mind, using a service-oriented architecture that can be horizontally scaled.

## Architecture

The system consists of three main microservices:

### Core Service (API Gateway)
- **Port**: 8080
- **Role**: Main API gateway handling event CRUD operations, venue management, and routing to specialized services
- **Responsibilities**:
  - Event creation, retrieval, and management
  - Venue management
  - Ticket listing for events
  - Orchestration of search and booking services
  - Elasticsearch indexing for new events

### Booking Service
- **Port**: 8081
- **Role**: Handles all ticket reservation and purchase operations
- **Responsibilities**:
  - Atomic ticket reservation using Redis distributed locks
  - Ticket purchase with payment processing
  - Reservation TTL management (180 seconds)
  - Transaction management for purchases

### Search Service
- **Port**: 8082
- **Role**: Dedicated search service for event queries
- **Responsibilities**:
  - Full-text search across events using Elasticsearch
  - Query optimization and result ranking

### Data Stores

- **PostgreSQL**: Primary database for events, venues, tickets, and purchases
- **Redis**: Distributed locking for ticket reservations
- **Elasticsearch**: Full-text search index for events

### Frontend

- **Next.js 16**: React-based UI with TypeScript
- **Features**:
  - Event search interface
  - Seat selection with visual representation
  - Reservation timer with countdown
  - Checkout flow
  - Purchase confirmation

## Features

### Implemented Features

✅ **Event Management**
- Create events with ticket allocations (VIP, GA, Front Row)
- List all events
- Get event details with available tickets
- Search events by query (title, description, venue)

✅ **Ticket Reservation**
- Atomic multi-ticket reservation
- Distributed locking via Redis
- 180-second reservation TTL
- Automatic rollback on partial failures

✅ **Ticket Purchase**
- Transactional purchase flow
- Payment processing (mock Stripe integration)
- Purchase history tracking
- Automatic reservation release on purchase

✅ **Search**
- Full-text search across event titles, descriptions, and venues
- Elasticsearch-powered indexing
- Pagination support

✅ **UI/UX**
- Modern, responsive design
- Real-time reservation timer
- Seat selection interface
- Smooth checkout flow

### Partially Implemented

⚠️ **Health Checks**: Endpoints are stubbed but not fully implemented

## Tech Stack

### Backend
- **Go 1.21+**: All services written in Go
- **Chi Router**: HTTP routing
- **SQLC**: Type-safe SQL code generation
- **Goose**: Database migrations
- **PostgreSQL 16**: Primary database
- **Redis 7**: Distributed locking
- **Elasticsearch 8.15**: Full-text search

### Frontend
- **Next.js 16**: React framework
- **TypeScript**: Type safety
- **Tailwind CSS**: Styling
- **Radix UI**: Component primitives

### Infrastructure
- **Docker & Docker Compose**: Containerization and orchestration
- **Multi-network architecture**: Public and internal networks for service isolation

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for UI development)

### Running the Application

1. **Clone the repository**
   ```bash
   git clone https://github.com/Ignisrex/tix.git
   cd tix
   ```

2. **Start all services with Docker Compose**
   ```bash
   docker-compose up --build
   ```

   This will start:
   - PostgreSQL database (port 5532)
   - Redis (internal network)
   - Elasticsearch (port 9200)
   - Core service (port 8080)
   - Booking service (port 8081)
   - Search service (port 8082)
   - Database migrations (automatic)

3. **Start the UI (in a separate terminal)**
   ```bash
   cd ui
   npm install
   npm run dev
   ```

   The UI will be available at `http://localhost:3000`

   Note: this isnt put into container for the ease of development; it should be added to the container in the future.

4. **Seed Data** (optional)
   
   The core service automatically seeds sample data on startup if `SEED_ON_START=true` is set (default in docker-compose.yaml).

### Environment Variables

#### Core Service
```
PORT=8080
DB_USER=postgres
DB_PASSWORD=password
DB_HOST=db
DB_PORT=5432
DB_NAME=tix_db
ES_HOST=elasticsearch
ES_PORT=9200
SEARCH_SERVICE_URL=http://search:8082
BOOKING_SERVICE_URL=http://booking:8081
SEED_ON_START=true
```

#### Booking Service
```
PORT=8081
DB_USER=postgres
DB_PASSWORD=password
DB_HOST=db
DB_PORT=5432
DB_NAME=tix_db
REDIS_HOST=ticket-lock
REDIS_PORT=6379
```

#### Search Service
```
PORT=8082
ES_HOST=elasticsearch
ES_PORT=9200
DB_USER=postgres
DB_PASSWORD=password
DB_HOST=db
DB_PORT=5432
DB_NAME=tix_db
```

## API Documentation

### Core Service (Port 8080)

#### Events

**GET `/api/v1/events`**
- List all events

**GET `/api/v1/events/:id`**
- Get event details

**POST `/api/v1/events`**
- Create a new event
- Body:
  ```json
  {
    "title": "Event Title",
    "description": "Event Description",
    "start_date": "2024-12-31T20:00:00Z",
    "venue_id": "uuid",
    "ticket_allocation": {
      "vip": 10,
      "ga": 100,
      "front_row": 20
    }
  }
  ```

**GET `/api/v1/events/:id/tickets`**
- Get all tickets for an event

**GET `/api/v1/events/search?q=query&limit=10&offset=0`**
- Search events (delegates to search service)

#### Venues

**GET `/api/v1/venues`**
- List all venues

**POST `/api/v1/venues`**
- Create a new venue
- Body:
  ```json
  {
    "name": "Venue Name",
    "location": "City, State"
  }
  ```

#### Booking

**POST `/api/v1/booking/reserve`**
- Reserve tickets
- Body:
  ```json
  {
    "ticket_ids": ["uuid1", "uuid2"]
  }
  ```
- Returns: Reservation confirmation with ticket IDs and TTL

**POST `/api/v1/booking/purchase`**
- Purchase reserved tickets
- Body:
  ```json
  {
    "ticket_ids": ["uuid1", "uuid2"]
  }
  ```
- Returns: Purchase confirmation with total amount

## Scaling Considerations

### Service Scaling

**Search & Booking Services**
- Scale horizontally by adding load balancers in front
- Implement service discovery for dynamic instance management
- Stateless design allows for easy horizontal scaling

**Core Service**
- Options for scaling:
  1. **High Availability**: Run multiple replicas with failover capabilities
  2. **Split Architecture**: Extract CRUD operations into separate microservice, make core solely an API gateway
- Recommended: Use Kubernetes for orchestration, auto-scaling, and self-healing

### Database Scaling

**PostgreSQL**
- Primary with read replicas for read-heavy workloads
- Event creation can tolerate eventual consistency
- **Ticket Status Challenge**: For sold tickets, introduce Redis cache layer:
  - Cache sold ticket IDs in Redis
  - Enrich ticket reads with cache data
  - Trade-off: Additional computational overhead for improved consistency

**Redis**
- Scale horizontally (Redis Cluster)
- Used for distributed locking and reservation management

**Elasticsearch**
- Scale horizontally
- **Production Recommendation**: Implement CDC (Change Data Capture) from PostgreSQL → Elasticsearch for real-time index updates

## Future Improvements

### High Priority

1. **Geospatial Search**
   - Integrate PostGIS or geospatial database
   - Enhance search with location-based queries
   - Implement quadtree or use Redis geospatial features for efficient location queries
   - Display suggested events based on user location

2. **Scalping Detection & Protection**
   - Implement rate limiting per user/IP
   - Detect and prevent automated ticket purchasing
   - Add CAPTCHA for high-demand events
   - Monitor for suspicious patterns

3. **High Traffic Event Handling**
   - Mark events as "high traffic" at creation
   - Implement virtual queue/line system
   - Gradual user release based on:
     - System load
     - Ticket availability
     - Time-based release windows
   - Dynamic throttling to prevent system overload

4. **Health Check Endpoints**
   - Implement `/healthz` endpoints for all services
   - Include database connectivity checks
   - Elasticsearch connectivity checks
   - Redis connectivity checks
   - Critical for Kubernetes liveness/readiness probes

5. **Observability**
   - Integrate Prometheus for metrics
   - Set up Grafana dashboards
   - Structured logging
   - Distributed tracing

### Medium Priority

6. **Enhanced Seat Map**
   - Store venue-specific seating layout JSON
   - Visual representation of actual venue layout
   - Clear seat-to-ticket associations
   - Interactive seat selection

7. **Global Utils Package**
   - Extract shared types and utilities into common package
   - Separate Go module at root level
   - Imported by all services
   - Reduces code duplication

## Known issues and edges:
1. **Reservation timer only reflects the most recent reservation time**
   - When a user reserves some tickets, waits, and then reserves additional tickets, the UI currently shows a single countdown based on the **latest** reservation.
   - Earlier tickets still expire at their original TTL; they are **not** extended, but the single timer can give the impression that all tickets share the new, later expiry.
   - This can be confusing to the user.

2. **Planned improvements for reservation UX**
   - **Group tickets by TTL**: have the booking service return per-ticket expiration times (for example `expires_at` / `ttl_seconds` per ticket), then group tickets in the UI by shared expiry and display them as a list, e.g.:
     - “2 tickets reserved until 10:05:12”
     - “1 ticket reserved until 10:03:47”
   - **Make the reservation timer dismissible**: similar to the cart drawer, allow users to hide the reservation timer/banner once they understand the timing, while still enforcing expiry server-side.

## Tradeoffs & Design Decisions

### Atomic Operations vs. Partial Success

**Current Implementation: All-or-Nothing**

The system currently uses atomic operations for ticket reservations and purchases. If any ticket in a multi-ticket request fails, the entire operation is rolled back.

**Trade-off Analysis:**

**All-or-Nothing (Current Approach)**
- ✅ **Pros**:
  - Simpler implementation
  - Guarantees consistency
  - No partial state to handle
  - Easier to reason about
- ❌ **Cons**:
  - Poor user experience if only one ticket is unavailable
  - Users must retry with different selections
  - Can lead to frustration during high-demand events

**Partial Success (Alternative Approach)**
- ✅ **Pros**:
  - Better user experience - users get what's available
  - Reduces retry attempts
  - More resilient to partial failures
  - Users can proceed with available tickets
- ❌ **Cons**:
  - More complex implementation
  - Need to handle partial state
  - Requires careful rollback logic
  - More edge cases to test

**Recommendation for Production:**

For a production system, **partial success** would provide a significantly better user experience. The implementation would:

1. Attempt to reserve each ticket individually
2. Return a response indicating which tickets were successfully reserved
3. Allow users to proceed with available tickets or cancel
4. Provide clear feedback about unavailable tickets

This approach is especially important for high-demand events where ticket availability changes rapidly. The additional complexity is justified by the improved user experience and reduced support burden.

### Eventual Consistency for Event Creation

**Decision**: Event creation uses eventual consistency for Elasticsearch indexing.

- Event is created in PostgreSQL immediately
- Elasticsearch indexing happens asynchronously
- Trade-off: Search may not immediately reflect new events, but system remains available

### Reservation TTL

**Decision**: 180-second (3-minute) reservation window.

- Balances user checkout time with ticket availability
- Prevents long-held reservations from blocking sales
- Client-side timer provides user feedback
- Future: Server-side validation and webhook notifications for accuracy (partially implemented)

## Project Structure

```
tix/
├── core/              # API Gateway & CRUD service
│   ├── cmd/          # Application entry points
│   ├── internal/     # Internal packages
│   ├── service/      # Business logic
│   ├── mappers/      # Data transformations
│   └── sql/          # SQL queries
├── booking/          # Booking service
│   ├── cmd/
│   ├── internal/
│   └── service/
├── search/           # Search service
│   ├── cmd/
│   ├── internal/
│   └── service/
├── ui/               # Next.js frontend
│   ├── src/
│   │   ├── app/      # Next.js app router
│   │   ├── components/
│   │   ├── lib/
│   │   └── types/
├── db/               # Database migrations
│   └── migrations/
└── docker-compose.yaml
```

## License

This project was created as a take-home assignment.

