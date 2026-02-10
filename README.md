# 🎯 Real-Time Quiz API

A high-performance real-time quiz application built with **Go (Golang)**, featuring WebSocket-based live communication, real-time score updates, and dynamic leaderboard functionality.

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [API Documentation](#-api-documentation)
- [WebSocket Events](#-websocket-events)
- [Database Schema](#-database-schema)
- [Architecture](#-architecture)
- [Testing](#-testing)

---

## 🎮 Overview

This project implements a real-time quiz system that allows users to:

- Join quiz sessions using a unique quiz ID
- Answer questions and receive instant score updates
- Compete with other participants simultaneously
- View a live leaderboard that updates in real-time

## ✨ Features

### Core Requirements

| Feature                | Status | Description                                                  |
| ---------------------- | ------ | ------------------------------------------------------------ |
| **User Participation** | ✅     | Users can join quiz sessions via unique quiz ID              |
| **Multi-user Support** | ✅     | Multiple users can join the same quiz session simultaneously |
| **Real-Time Scores**   | ✅     | Scores update instantly as users submit answers              |
| **Live Leaderboard**   | ✅     | Leaderboard reflects current standings in real-time          |
| **Idempotency**        | ✅     | Prevents duplicate answer submissions                        |

### Additional Features

- 🔐 **Authentication Middleware** - Secure user sessions
- 📊 **Accurate Scoring System** - Consistent and fair score calculation
- 🚀 **High Performance** - Built with Go for maximum concurrency
- 🔄 **WebSocket Communication** - Low-latency real-time updates
- 📝 **RESTful API** - Standard HTTP endpoints for quiz management

## 🛠 Tech Stack

| Component            | Technology              |
| -------------------- | ----------------------- |
| **Language**         | Go 1.23+                |
| **Web Framework**    | Gin                     |
| **WebSocket**        | Gorilla WebSocket       |
| **Database**         | PostgreSQL 16           |
| **Cache/Pub-Sub**    | Redis 7                 |
| **Containerization** | Docker & Docker Compose |

## 📁 Project Structure

```
realtime-quiz/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── config/
│   ├── local.yaml               # Local development config
│   └── production.yaml          # Production config
├── internal/
│   ├── bootstrap/               # Application initialization
│   │   ├── config.go            # Configuration loader
│   │   ├── logger.go            # Logger setup
│   │   ├── postgres.go          # PostgreSQL connection
│   │   ├── redis.go             # Redis connection
│   │   ├── router.go            # HTTP router setup
│   │   └── server.go            # Server initialization
│   ├── domain/                  # Business domain objects
│   │   ├── errors.go            # Custom error definitions
│   │   ├── leaderboard.go       # Leaderboard domain logic
│   │   ├── quiz.go              # Quiz domain logic
│   │   └── session.go           # Session domain logic
│   ├── handler/                 # HTTP handlers
│   │   ├── quiz.handler.go      # Quiz endpoints
│   │   └── user.handler.go      # User endpoints
│   ├── middleware/              # HTTP middlewares
│   │   ├── auth.middleware.go   # Authentication
│   │   ├── cors.middleware.go   # CORS handling
│   │   ├── logger.middleware.go # Request logging
│   │   └── recovery.middleware.go # Panic recovery
│   ├── models/                  # Database models
│   │   ├── question.go          # Question model
│   │   ├── quiz.go              # Quiz model
│   │   └── result.go            # Result model
│   ├── realtime/                # WebSocket & real-time logic
│   │   ├── broadcaster.go       # Message broadcasting
│   │   ├── client.go            # WebSocket client
│   │   ├── handler.go           # WebSocket handler
│   │   ├── hub.go               # Connection hub
│   │   ├── message.go           # Message types
│   │   ├── player.go            # Player state
│   │   └── session.go           # Game session
│   ├── repository/              # Data access layer
│   │   ├── interfaces.go        # Repository interfaces
│   │   ├── question.repo.go     # Question repository
│   │   ├── quiz.repo.go         # Quiz repository
│   │   └── answer.repo.go       # Answer repository
│   └── service/                 # Business logic layer
│       ├── interfaces.go        # Service interfaces
│       ├── leaderboard.service.go # Leaderboard logic
│       └── quiz.service.go      # Quiz logic
├── migrations/                  # Database migrations
├── pkg/
│   └── response/                # HTTP response helpers
│       └── response.go
├── tests/                       # Integration tests
│   ├── api/                     # API Integration tests
│   └── manual/                  # Manual test scripts
├── web/                         # Static web files (optional)
├── docker-compose.yaml          # Docker services config
├── go.mod                       # Go module definition
└── Makefile                     # Build automation
```

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or higher
- Docker & Docker Compose
- Make (optional)

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd realtime-quiz
   ```

2. **Start infrastructure services**

   ```bash
   docker-compose up -d
   ```

3. **Install dependencies**

   ```bash
   go mod download
   ```

4. **Run database migrations**

   ```bash
   # Using go-migrate or similar tool
   migrate -path ./migrations -database "postgres://quiz:quiz123@localhost:5433/realtime_quiz?sslmode=disable" up
   ```

5. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

The server will be available at `http://localhost:8080`

### Environment Configuration

| Variable      | Default       | Description       |
| ------------- | ------------- | ----------------- |
| `SERVER_PORT` | 8080          | HTTP server port  |
| `DB_HOST`     | localhost     | PostgreSQL host   |
| `DB_PORT`     | 5433          | PostgreSQL port   |
| `DB_USER`     | quiz          | Database user     |
| `DB_PASSWORD` | quiz123       | Database password |
| `DB_NAME`     | realtime_quiz | Database name     |
| `REDIS_HOST`  | localhost     | Redis host        |
| `REDIS_PORT`  | 6379          | Redis port        |

## 📖 API Documentation

### REST Endpoints

#### Quiz Management

| Method | Endpoint                     | Description         |
| ------ | ---------------------------- | ------------------- |
| `POST` | `/api/v1/quiz`               | Create a new quiz   |
| `GET`  | `/api/v1/quiz/:id`           | Get quiz details    |
| `GET`  | `/api/v1/quiz/:id/questions` | Get quiz questions  |
| `POST` | `/api/v1/quiz/:id/join`      | Join a quiz session |

#### User

| Method | Endpoint                | Description         |
| ------ | ----------------------- | ------------------- |
| `POST` | `/api/v1/user/register` | Register a new user |
| `POST` | `/api/v1/user/login`    | User login          |

#### Leaderboard

| Method | Endpoint                       | Description             |
| ------ | ------------------------------ | ----------------------- |
| `GET`  | `/api/v1/quiz/:id/leaderboard` | Get current leaderboard |

### WebSocket Endpoint

```
ws://localhost:8080/ws/quiz/:quiz_id?token=<jwt_token>
```

_Note: If the `Authorization` header cannot be set (e.g., in standard JS WebSocket), pass the token via the `token` query parameter._

## 🔌 WebSocket Events

### Client → Server

| Event           | Payload                                           | Description        |
| --------------- | ------------------------------------------------- | ------------------ |
| `join_session`  | `{ "user_id": "string", "username": "string" }`   | Join quiz session  |
| `submit_answer` | `{ "question_id": "string", "answer": "string" }` | Submit answer      |
| `leave_session` | `{}`                                              | Leave quiz session |

### Server → Client

| Event                | Payload                                                  | Description          |
| -------------------- | -------------------------------------------------------- | -------------------- |
| `session_joined`     | `{ "session_id": "string", "participants": [...] }`      | Confirmation of join |
| `new_question`       | `{ "question": {...}, "time_limit": 30 }`                | Next question        |
| `score_update`       | `{ "user_id": "string", "score": 100, "correct": true }` | Score update         |
| `leaderboard_update` | `{ "rankings": [...] }`                                  | Updated leaderboard  |
| `quiz_ended`         | `{ "final_rankings": [...], "winner": {...} }`           | Quiz completion      |

## 💾 Database Schema

### Quiz

| Column        | Type      | Description           |
| ------------- | --------- | --------------------- |
| `id`          | UUID      | Primary key           |
| `title`       | VARCHAR   | Quiz title            |
| `description` | TEXT      | Quiz description      |
| `created_at`  | TIMESTAMP | Creation timestamp    |
| `updated_at`  | TIMESTAMP | Last update timestamp |

### Question

| Column           | Type    | Description               |
| ---------------- | ------- | ------------------------- |
| `id`             | UUID    | Primary key               |
| `quiz_id`        | UUID    | Foreign key to Quiz       |
| `content`        | TEXT    | Question text             |
| `options`        | JSONB   | Answer options            |
| `correct_answer` | VARCHAR | Correct answer            |
| `points`         | INTEGER | Points for correct answer |
| `time_limit`     | INTEGER | Time limit in seconds     |

### Result

| Column         | Type      | Description          |
| -------------- | --------- | -------------------- |
| `id`           | UUID      | Primary key          |
| `quiz_id`      | UUID      | Foreign key to Quiz  |
| `user_id`      | UUID      | User identifier      |
| `score`        | INTEGER   | Total score          |
| `completed_at` | TIMESTAMP | Completion timestamp |

## 🏗 Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                            │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │  REST API   │    │  WebSocket  │    │   Web Frontend      │  │
│  └──────┬──────┘    └──────┬──────┘    └─────────────────────┘  │
└─────────┼──────────────────┼────────────────────────────────────┘
          │                  │
          ▼                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Handler Layer                            │
│  ┌─────────────┐    ┌─────────────┐                             │
│  │ Quiz Handler│    │  WS Handler │                             │
│  └──────┬──────┘    └──────┬──────┘                             │
└─────────┼──────────────────┼────────────────────────────────────┘
          │                  │
          ▼                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Service Layer                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │Quiz Service │    │ Leaderboard │    │   Session Manager   │  │
│  └──────┬──────┘    └──────┬──────┘    └──────────┬──────────┘  │
└─────────┼──────────────────┼─────────────────────┼──────────────┘
          │                  │                     │
          ▼                  ▼                     ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Repository Layer                           │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────┐  │
│  │ Quiz Repo   │    │Question Repo│    │    Result Repo      │  │
│  └──────┬──────┘    └──────┬──────┘    └──────────┬──────────┘  │
└─────────┼──────────────────┼─────────────────────┼──────────────┘
          │                  │                     │
          ▼                  ▼                     ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Data Layer                                 │
│  ┌─────────────────────────┐    ┌─────────────────────────────┐ │
│  │       PostgreSQL        │    │           Redis             │ │
│  │   (Persistent Data)     │    │   (Cache & Pub/Sub)         │ │
│  └─────────────────────────┘    └─────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### Key Design Decisions

1. **WebSocket Hub Pattern** - Centralized connection management for broadcasting
2. **Clean Architecture** - Separation of concerns between layers
3. **Redis Pub/Sub** - Enables horizontal scaling for real-time updates
4. **Repository Pattern** - Abstracts data access for testability

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/service/...

# Run integration tests
go test ./tests/api/...
```

---

## 📝 License

This project is created for assessment purposes only and is not intended for commercial use.

## 👤 Author

Created as a Back-End Golang Technical Assessment submission.
