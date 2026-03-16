# GameVote

A real-time party voting application for selecting games to play together. Built with Go backend and Vue.js frontend.

## Project Structure

```
gamevote-x/
├── gamevote-api-go/    # Go backend API
│   ├── cmd/api/        # Main application entry point
│   ├── internal/       # Internal application code
│   │   ├── handler/    # HTTP handlers
│   │   ├── service/    # Business logic
│   │   ├── storage/    # Database repositories
│   │   ├── models/     # Data models
│   │   ├── helpers/    # Utility functions
│   │   └── logger/     # Logging configuration
│   └── docs/           # Swagger documentation
└── gamevote-ui/        # Vue.js frontend
    ├── src/            # Source code
    └── public/         # Static assets
```

## Features

- **Party Management**: Create and manage voting parties with unique codes
- **Game Nomination**: Add Steam games as voting options
- **Real-time Voting**: Vote on games with live updates via Server-Sent Events (SSE)
- **Beer Tracking**: Track beer consumption per attendee
- **Steam Integration**: Search and add games from Steam's catalog
- **User Management**: Simple username-based authentication with session persistence

## Tech Stack

### Backend
- **Go** - Main programming language
- **Gin** - HTTP web framework
- **SurrealDB** - Database
- **Swagger** - API documentation
- **Server-Sent Events** - Real-time updates

### Frontend
- **Vue.js 3** - Frontend framework
- **TypeScript** - Type safety
- **Bun** - Package manager and build tool

## Prerequisites

- Go 1.21+
- Node.js 18+
- Bun (package manager)
- SurrealDB

## Setup & Installation

### Database Setup

1. Install and start SurrealDB:
```bash
# Install SurrealDB
curl --proto '=https' --tlsv1.2 -sSf https://install.surrealdb.com | sh

# Start SurrealDB
surreal start --log trace --user root --pass root memory
```

### Backend Setup

1. Navigate to the backend directory:
```bash
cd gamevote-api-go
```

2. Install dependencies:
```bash
go mod download
```

3. Create environment file (optional):
```bash
cp .env.example .env
```

4. Environment variables:
```env
SURREAL_WS=ws://localhost:8000/rpc
SURREAL_USER=root
SURREAL_PASS=root
SURREAL_NS=gamevote
SURREAL_DB=gamevote
SERVER_PORT=8080
```

5. Run the backend:
```bash
go run cmd/api/main.go
```

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd gamevote-ui
```

2. Install dependencies:
```bash
bun install
```

3. Start development server:
```bash
bun run dev
```

## Authentication

The application uses a simple username-based authentication system with the following behavior:

- **First Visit**: Users are prompted to enter a username on the login page
- **Session Persistence**: Username is stored in browser cookies for automatic login
- **Database Validation**: If a cookie exists but the user is not found in the database, the user is automatically redirected to the login page to re-authenticate
- **Auto-Registration**: New usernames are automatically registered in the database upon first login

## Usage

1. Open your browser and navigate to `http://localhost:5173`
2. Enter a username if prompted (first visit or invalid session)
3. Create a new party or join existing one with a party code
4. Add attendees to the party
5. Nominate games by searching Steam's catalog
6. Start voting phase when ready
7. Vote on games (thumbs up/down/neutral)
8. View results when all attendees have voted

## API Documentation

The API documentation is available via Swagger at:
`http://localhost:8080/swagger/index.html`

## Development

### Generate API Client

The frontend uses a generated TypeScript client. To regenerate:

```bash
cd gamevote-api-go
./scripts/generate-ts-client.sh
```

### Key Endpoints

- `GET /parties` - List all parties
- `POST /parties` - Create a new party
- `GET /parties/{code}` - Get party details
- `GET /parties/{code}/stream` - SSE stream for real-time updates
- `POST /parties/{code}/votes/{attendee}` - Submit votes
- `GET /games?q=search` - Search Steam games

## Architecture

### Party States
- **NOMINATION**: Adding games and attendees
- **VOTING**: Attendees vote on nominated games
- **RESULTS**: Display voting results

### Real-time Updates
Uses Server-Sent Events (SSE) to broadcast:
- Party state changes
- Online user list updates
- Outstanding voters updates

### Authentication Flow
- Username-based sessions with cookie persistence
- Automatic user registration on first login
- Session validation against database on each request
- Automatic redirect to login page if cookie exists but user not found in database

### Database Schema
- **parties**: Party information and state
- **polls**: Individual voting sessions
- **users**: User accounts
- **games**: Steam game cache
- **beers**: Beer consumption tracking
- **drink_types**: Drink presets

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

[Add your license information here]