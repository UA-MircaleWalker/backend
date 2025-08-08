# UA Card Battle Game - API Documentation

## Overview

This document provides comprehensive API documentation for the UA Card Battle Game microservices architecture. The system consists of 5 main microservices, each responsible for specific game functionality.

## Base URL

```
Development: http://localhost
Production: https://your-domain.com
```

## Authentication

All protected endpoints require a Bearer token in the Authorization header:

```http
Authorization: Bearer <jwt_token>
```

### Obtaining a Token

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "your_password"
}
```

Response:
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "username": "player123",
      "display_name": "Player 123"
    }
  }
}
```

## Service Endpoints

### 1. Card Service (Port 8001)

Manages card data, rules, and validation.

#### Get All Cards
```http
GET /api/v1/cards?page=1&limit=20&card_type=CHARACTER&rarity=RARE
```

Query Parameters:
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 20, max: 100)
- `card_type` (string): Filter by card type (CHARACTER, FIELD, EVENT, AP)
- `work_code` (string): Filter by work code (e.g., "UA2")
- `rarity` (string): Filter by rarity (COMMON, UNCOMMON, RARE, SUPER_RARE, SPECIAL)
- `characteristics` (string): Comma-separated characteristics
- `keywords` (string): Comma-separated keywords
- `min_bp`, `max_bp` (int): BP range filter
- `min_ap_cost`, `max_ap_cost` (int): AP cost range filter
- `search_name` (string): Search by card name

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": "card-uuid",
      "card_number": "UA25-001",
      "name": "Hero Character",
      "card_type": "CHARACTER",
      "work_code": "UA2",
      "bp": 5,
      "ap_cost": 2,
      "energy_cost": {"red": 1, "blue": 1},
      "rarity": "RARE",
      "characteristics": ["Hero", "Human"],
      "keywords": ["レイド"],
      "effect_text": "When played, draw a card.",
      "image_url": "https://example.com/card.jpg"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 156,
    "total_pages": 8
  }
}
```

#### Get Card by ID
```http
GET /api/v1/cards/{card_id}
```

#### Search Cards
```http
GET /api/v1/cards/search?q=hero&limit=10
```

#### Validate Deck Composition
```http
POST /api/v1/cards/validate-deck
Content-Type: application/json

[
  {"card_id": "uuid", "quantity": 3},
  {"card_id": "uuid", "quantity": 2}
]
```

Response:
```json
{
  "success": true,
  "data": {
    "is_valid": true,
    "errors": [],
    "warnings": ["Deck contains only 45 cards, consider adding more"],
    "card_count": 45,
    "work_breakdown": {"UA2": 30, "LLG": 15},
    "type_breakdown": {"CHARACTER": 25, "EVENT": 15, "FIELD": 5}
  }
}
```

#### Create Card (Admin)
```http
POST /api/v1/cards
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "card_number": "UA25-157",
  "name": "New Hero",
  "card_type": "CHARACTER",
  "work_code": "UA2",
  "bp": 6,
  "ap_cost": 3,
  "energy_cost": {"red": 2},
  "rarity": "SUPER_RARE",
  "characteristics": ["Hero", "Legendary"],
  "effect_text": "Destroy all enemy characters.",
  "keywords": ["ダメージ3"],
  "image_url": "https://example.com/new-hero.jpg"
}
```

### 2. User Service (Port 8002)

Handles authentication, user profiles, and deck management.

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "player123",
  "email": "player@example.com",
  "password": "secure_password",
  "display_name": "Player 123"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "player@example.com",
  "password": "secure_password"
}
```

#### Get User Profile
```http
GET /api/v1/users/profile
Authorization: Bearer <token>
```

Response:
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "user-uuid",
      "username": "player123",
      "display_name": "Player 123",
      "level": 15,
      "experience": 12500,
      "rank": 42,
      "rank_points": 1750
    },
    "stats": {
      "games_played": 89,
      "games_won": 54,
      "win_rate": 0.6067
    },
    "recent_matches": [...],
    "achievements": [...]
  }
}
```

#### Get User Decks
```http
GET /api/v1/users/{user_id}/decks
Authorization: Bearer <token>
```

#### Create Deck
```http
POST /api/v1/users/decks
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "My Tournament Deck",
  "cards": [
    {"card_id": "uuid", "quantity": 3},
    {"card_id": "uuid", "quantity": 2}
  ]
}
```

#### Get User Collection
```http
GET /api/v1/users/{user_id}/collection
Authorization: Bearer <token>
```

### 3. Matchmaking Service (Port 8003)

Manages player queues and match creation.

#### Join Queue
```http
POST /api/v1/matchmaking/queue
Authorization: Bearer <token>
Content-Type: application/json

{
  "mode": "RANKED",
  "rank_range": 1500
}
```

Response:
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Successfully joined queue",
    "position": 3,
    "estimated_wait": 45,
    "joined_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Leave Queue
```http
DELETE /api/v1/matchmaking/queue/{user_id}
Authorization: Bearer <token>
```

#### Get Queue Status
```http
GET /api/v1/matchmaking/status/{user_id}
Authorization: Bearer <token>
```

#### Get Queue Statistics
```http
GET /api/v1/matchmaking/stats
```

Response:
```json
{
  "success": true,
  "data": {
    "ranked_queue": 15,
    "casual_queue": 8,
    "active_matches": 23
  }
}
```

#### Accept/Decline Match
```http
POST /api/v1/matchmaking/accept
Authorization: Bearer <token>
Content-Type: application/json

{
  "match_id": "match-uuid"
}
```

```http
POST /api/v1/matchmaking/decline
Authorization: Bearer <token>
Content-Type: application/json

{
  "match_id": "match-uuid"
}
```

### 4. Game Battle Service (Port 8004)

Real-time game logic and state management.

#### Create Game
```http
POST /api/v1/games
Authorization: Bearer <token>
Content-Type: application/json

{
  "player1_id": "uuid",
  "player2_id": "uuid",
  "mode": "RANKED"
}
```

#### Get Game State
```http
GET /api/v1/games/{game_id}
Authorization: Bearer <token>
```

Response:
```json
{
  "success": true,
  "data": {
    "id": "game-uuid",
    "status": "IN_PROGRESS",
    "current_turn": 3,
    "phase": "MAIN",
    "active_player": "player1-uuid",
    "game_state": {
      "turn": 3,
      "phase": "MAIN",
      "active_player": "player1-uuid",
      "players": {
        "player1-uuid": {
          "ap": 5,
          "max_ap": 5,
          "energy": {"red": 3, "blue": 2},
          "hand": [...],
          "characters": [...],
          "fields": [...]
        },
        "player2-uuid": {...}
      },
      "board": {...}
    }
  }
}
```

#### Perform Game Action
```http
POST /api/v1/games/{game_id}/actions
Authorization: Bearer <token>
Content-Type: application/json

{
  "action_type": "PLAY_CARD",
  "action_data": {
    "card_id": "card-uuid",
    "position": {"zone": 1, "slot": 2},
    "target_id": "target-uuid"
  }
}
```

Response:
```json
{
  "success": true,
  "data": {
    "success": true,
    "game_state": {...},
    "effects": [
      {
        "type": "CARD_PLAYED",
        "source": "card-uuid",
        "description": "Hero Character was played",
        "applied": true
      }
    ],
    "events_triggered": [
      {
        "type": "CARD_PLAYED",
        "source": "card-uuid",
        "timestamp": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

#### Get Game Actions History
```http
GET /api/v1/games/{game_id}/actions?from_index=0
Authorization: Bearer <token>
```

#### WebSocket Connection
```javascript
const ws = new WebSocket('ws://localhost/ws');

// Join a game room
ws.send(JSON.stringify({
  type: 'JOIN_GAME',
  game_id: 'game-uuid'
}));

// Listen for game updates
ws.onmessage = function(event) {
  const message = JSON.parse(event.data);
  console.log('Game update:', message);
};
```

### 5. Game Result Service (Port 8005)

Statistics, leaderboards, and analytics.

#### Record Game Result
```http
POST /api/v1/results
Authorization: Bearer <token>
Content-Type: application/json

{
  "game_id": "game-uuid",
  "player1_id": "player1-uuid",
  "player2_id": "player2-uuid",
  "winner": "player1-uuid",
  "game_duration": 720,
  "total_turns": 12,
  "end_reason": "NORMAL_WIN",
  "game_mode": "RANKED"
}
```

#### Get Player Statistics
```http
GET /api/v1/results/{user_id}/stats
```

Response:
```json
{
  "success": true,
  "data": {
    "stats": {
      "user_id": "user-uuid",
      "games_played": 156,
      "games_won": 94,
      "games_lost": 59,
      "games_drawn": 3,
      "win_rate": 0.6026,
      "current_streak": 5,
      "best_streak": 12,
      "total_game_time": 45600,
      "avg_game_time": 292,
      "rank_points": 1847,
      "current_rank": 42
    },
    "achievements": [...],
    "rank_history": [...]
  }
}
```

#### Get Leaderboard
```http
GET /api/v1/leaderboard?page=1&limit=50&time_frame=week&mode=ranked
```

Response:
```json
{
  "success": true,
  "data": {
    "entries": [
      {
        "rank": 1,
        "user": {
          "id": "top-player-uuid",
          "username": "champion",
          "display_name": "Champion Player"
        },
        "stats": {
          "rank_points": 2456,
          "games_played": 203,
          "win_rate": 0.7537
        },
        "rank_change": 2
      }
    ],
    "total": 1250,
    "page": 1,
    "limit": 50,
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Get Match History
```http
GET /api/v1/results/{user_id}/history?page=1&limit=20
Authorization: Bearer <token>
```

#### Get Game Analytics
```http
GET /api/v1/analytics/overview?start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer <token>
```

Response:
```json
{
  "success": true,
  "data": {
    "overview": {
      "total_games": 15678,
      "total_players": 3456,
      "active_players": 1234,
      "avg_game_length": 425,
      "popular_cards": [
        {
          "card_id": "card-uuid",
          "card_name": "Popular Hero",
          "usage_count": 8934,
          "win_rate": 0.6234
        }
      ]
    },
    "trend_data": [...],
    "top_performers": [...]
  }
}
```

#### Compare Players
```http
GET /api/v1/results/compare?player1=uuid1&player2=uuid2
```

## Error Handling

All endpoints return consistent error responses:

```json
{
  "success": false,
  "error": "Error type or message",
  "details": "Additional error information (optional)"
}
```

### HTTP Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict (e.g., user already in queue)
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

## Rate Limits

- Authentication endpoints: 5 requests/second
- General API endpoints: 10 requests/second
- Game actions: 20 requests/second

Rate limit headers are included in responses:
```http
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 7
X-RateLimit-Reset: 1642248000
```

## WebSocket Events

Real-time game updates are delivered via WebSocket connections.

### Client → Server Messages

```javascript
// Join game room
{
  "type": "JOIN_GAME",
  "game_id": "game-uuid"
}

// Leave game room
{
  "type": "LEAVE_GAME",
  "game_id": "game-uuid"
}

// Heartbeat
{
  "type": "PING"
}
```

### Server → Client Messages

```javascript
// Game state update
{
  "type": "GAME_STATE_UPDATE",
  "game_id": "game-uuid",
  "game_state": {...},
  "timestamp": "2024-01-15T10:30:00Z"
}

// Game action performed
{
  "type": "GAME_ACTION",
  "game_id": "game-uuid",
  "player_id": "player-uuid",
  "action": {...},
  "timestamp": "2024-01-15T10:30:00Z"
}

// Match found
{
  "type": "MATCH_FOUND",
  "match_id": "match-uuid",
  "opponent": {...},
  "game_mode": "RANKED"
}

// Game ended
{
  "type": "GAME_ENDED",
  "game_id": "game-uuid",
  "winner": "player-uuid",
  "reason": "NORMAL_WIN"
}
```

## SDK Examples

### JavaScript/TypeScript
```typescript
import { UAGameClient } from 'ua-game-sdk';

const client = new UAGameClient({
  baseUrl: 'http://localhost',
  token: 'your-jwt-token'
});

// Get cards
const cards = await client.cards.list({
  page: 1,
  limit: 20,
  card_type: 'CHARACTER'
});

// Join matchmaking
await client.matchmaking.joinQueue({
  mode: 'RANKED'
});

// Connect to WebSocket
const gameClient = client.games.connect();
gameClient.on('gameStateUpdate', (data) => {
  console.log('Game updated:', data);
});
```

### Python
```python
from ua_game_client import UAGameClient

client = UAGameClient(
    base_url='http://localhost',
    token='your-jwt-token'
)

# Get player stats
stats = client.results.get_player_stats('player-uuid')
print(f"Win rate: {stats['win_rate']:.2%}")

# Get leaderboard
leaderboard = client.results.get_leaderboard(
    page=1,
    limit=10,
    time_frame='week'
)
```

## Monitoring and Health Checks

Each service provides health check endpoints:

```http
GET /{service}/health
```

Response:
```json
{
  "status": "healthy",
  "service": "card-service",
  "version": "1.0.0",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Development Setup

See the main README.md file for development setup instructions including Docker Compose configuration and local development guidelines.

## Support

For API support and questions:
- GitHub Issues: https://github.com/your-org/ua-card-battle/issues
- Documentation: https://docs.your-domain.com
- Discord: https://discord.gg/your-server