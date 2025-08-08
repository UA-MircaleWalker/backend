# Redis Data Structures for UA Card Battle Game

This document describes the Redis key patterns and data structures used across all microservices.

## Key Naming Conventions

All keys follow the pattern: `{service}:{entity}:{identifier}[:{subentity}]`

- **TTL Policy**: Most keys have appropriate expiration times to prevent memory bloat
- **Encoding**: All complex objects are stored as JSON strings
- **Versioning**: Keys include version suffixes for breaking changes

## 1. Authentication & Sessions

### JWT Blacklist
```
auth:blacklist:{jti}
Type: String
Value: "blacklisted"
TTL: Token expiration time
Usage: Store blacklisted JWT tokens
```

### User Sessions
```
auth:session:{user_id}
Type: Hash
Fields:
  - device_id: "web_chrome_123"
  - ip_address: "192.168.1.1"
  - last_activity: "2024-01-15T10:30:00Z"
  - login_at: "2024-01-15T09:00:00Z"
TTL: 24 hours
Usage: Track active user sessions
```

### Rate Limiting
```
ratelimit:{service}:{endpoint}:{user_id}
Type: String
Value: Request count
TTL: Rate limit window (e.g., 60 seconds)
Usage: API rate limiting per user
```

## 2. Card Service Caching

### Card Data Cache
```
cards:data:{card_id}
Type: String (JSON)
Value: Complete card object with effects
TTL: 6 hours
Usage: Cache frequently accessed card data
```

### Card Search Results
```
cards:search:{hash}
Type: String (JSON)
Value: Array of card IDs matching search criteria
TTL: 30 minutes
Usage: Cache search results to reduce DB load
```

### Popular Cards
```
cards:popular:daily
Type: ZSet
Score: Usage count
Members: card_id
TTL: 24 hours
Usage: Track and display popular cards
```

### Card Balance Adjustments
```
cards:balance:{card_id}
Type: Hash
Fields:
  - original_bp: "5"
  - adjusted_bp: "4"
  - adjustment_reason: "overpowered"
  - applied_at: "2024-01-15T12:00:00Z"
TTL: 7 days
Usage: Track recent balance changes
```

## 3. User Management

### User Profile Cache
```
users:profile:{user_id}
Type: String (JSON)
Value: User profile with stats
TTL: 1 hour
Usage: Cache user profiles for quick access
```

### User Collections
```
users:collection:{user_id}
Type: Hash
Fields: {card_id}: quantity
TTL: 2 hours
Usage: Cache user card collections
```

### Active Decks
```
users:deck:active:{user_id}
Type: String (JSON)
Value: Complete deck object with cards
TTL: 1 hour
Usage: Cache active deck for quick game setup
```

### User Achievements Progress
```
users:achievements:{user_id}
Type: Hash
Fields: {achievement_id}: progress_value
TTL: 30 minutes
Usage: Track achievement progress
```

## 4. Matchmaking System

### Matchmaking Queues
```
matchmaking:queue:{mode}
Type: ZSet
Score: Timestamp (for FIFO) + rank_adjustment
Members: user_id
TTL: None (managed by cleanup)
Usage: Store players waiting for matches

Modes: RANKED, CASUAL, FRIEND
```

### User Queue Status
```
matchmaking:user:{user_id}
Type: String (JSON)
Value: {
  "mode": "RANKED",
  "joined_at": "2024-01-15T10:30:00Z",
  "rank_range": 1500,
  "estimated_wait": 120
}
TTL: 5 minutes
Usage: Store user's current queue state
```

### Active Matches
```
matchmaking:active_matches
Type: ZSet
Score: Creation timestamp
Members: match_id
TTL: None (cleaned after completion)
Usage: Track ongoing matches
```

### Match History Cache
```
matchmaking:history:{user_id}
Type: List
Values: JSON objects of recent matches
TTL: 1 hour
Usage: Cache recent match history
```

### Queue Statistics
```
matchmaking:stats:current
Type: Hash
Fields:
  - ranked_queue_size: "15"
  - casual_queue_size: "8"
  - avg_wait_time: "45"
  - matches_created_today: "127"
TTL: 30 seconds
Usage: Real-time queue statistics
```

## 5. Game Battle System

### Game State Storage
```
game:{game_id}:state
Type: String (JSON)
Value: Complete game state including:
  - turn, phase, active_player
  - player hands, decks, fields
  - board state, action history
TTL: 24 hours
Usage: Store real-time game state
```

### Player Active Games
```
game:player:{user_id}:active
Type: Set
Members: game_id
TTL: 24 hours
Usage: Track which games a player is in
```

### Game Actions History
```
game:{game_id}:actions
Type: List
Values: JSON action objects in chronological order
TTL: 24 hours
Usage: Store game action sequence for replay
```

### Game Locks (for concurrency)
```
game:{game_id}:lock
Type: String
Value: lock_token
TTL: 5 seconds
Usage: Prevent concurrent game state modifications
```

### Turn Timeouts
```
game:{game_id}:timeout
Type: String
Value: "turn_timeout"
TTL: Turn time limit (default 90 seconds)
Usage: Handle turn timeouts automatically
```

## 6. Real-time WebSocket Management

### Connected Clients
```
websocket:clients:{service}
Type: Hash
Fields: {connection_id}: user_id
TTL: None (managed by connection lifecycle)
Usage: Map WebSocket connections to users
```

### Game Room Subscriptions
```
websocket:room:{game_id}
Type: Set
Members: connection_id
TTL: None (managed by game lifecycle)
Usage: Track which connections are subscribed to game updates
```

### User Presence
```
websocket:presence:{user_id}
Type: Hash
Fields:
  - status: "online"
  - last_seen: "2024-01-15T10:30:00Z"
  - current_game: "game_id_123"
TTL: 5 minutes (refreshed by heartbeat)
Usage: Track user online status
```

## 7. Game Results & Statistics

### Leaderboard Cache
```
results:leaderboard:{timeframe}
Type: ZSet
Score: rank_points
Members: user_id
TTL: 5 minutes

Timeframes: all, week, month
Usage: Cache leaderboard rankings
```

### Player Statistics Cache
```
results:stats:{user_id}
Type: String (JSON)
Value: Complete player statistics
TTL: 10 minutes
Usage: Cache calculated player stats
```

### Daily Statistics
```
results:daily:{date}
Type: Hash
Fields:
  - total_games: "156"
  - unique_players: "89"
  - avg_duration: "720"
  - new_registrations: "12"
TTL: 7 days
Usage: Daily analytics aggregation
```

### Card Usage Analytics
```
results:card_usage:{card_id}:{date}
Type: Hash
Fields:
  - games_played: "15"
  - games_won: "9"
  - avg_turn_played: "3.2"
  - win_rate: "0.6000"
TTL: 30 days
Usage: Track card performance metrics
```

## 8. Distributed Locking

### Service Locks
```
lock:{resource_type}:{resource_id}
Type: String
Value: lock_owner_id
TTL: Lock timeout (usually 30 seconds)
Usage: Distributed locking for critical operations
```

### Maintenance Mode
```
system:maintenance
Type: String
Value: "maintenance_mode_active"
TTL: None (manually managed)
Usage: Global maintenance mode flag
```

## 9. Event Streaming (Redis Streams)

### Game Events Stream
```
stream:game_events
Type: Stream
Entries: {
  "event_type": "CARD_PLAYED",
  "game_id": "uuid",
  "player_id": "uuid", 
  "timestamp": "2024-01-15T10:30:00Z",
  "data": "{...}"
}
Usage: Event sourcing for game actions
```

### Achievement Events
```
stream:achievements
Type: Stream
Entries: Achievement unlock events
Usage: Process achievement unlocks asynchronously
```

### Analytics Events
```
stream:analytics
Type: Stream
Entries: User behavior and game metrics
Usage: Real-time analytics processing
```

## 10. Configuration & Feature Flags

### Service Configuration
```
config:{service_name}
Type: Hash
Fields: Configuration key-value pairs
TTL: None (updated via admin API)
Usage: Runtime service configuration
```

### Feature Flags
```
features:{feature_name}
Type: Hash
Fields:
  - enabled: "true"
  - rollout_percentage: "50"
  - target_users: "premium_users"
TTL: None (managed via admin interface)
Usage: Feature rollout control
```

## Memory Usage Optimization

### Key Expiration Policies
- **Game states**: Auto-expire after 24 hours
- **Cache data**: Short TTL with refresh on access
- **Temporary data**: Aggressive expiration (1-5 minutes)
- **Analytics**: Longer retention (7-30 days)

### Memory Management Commands
```bash
# Monitor memory usage
MEMORY USAGE {key}

# Get key statistics
INFO keyspace

# Clean expired keys
SCAN 0 MATCH "pattern*" COUNT 1000
```

### Data Structure Selection
- **Hash**: User profiles, game states
- **ZSet**: Leaderboards, queues, timestamps
- **Set**: Simple collections, active games
- **List**: Action histories, event logs
- **String**: Simple values, JSON objects
- **Stream**: Event sourcing, real-time processing

## Backup and Monitoring

### Health Check Keys
```
health:{service_name}
Type: Hash
Fields:
  - status: "healthy"
  - last_updated: "2024-01-15T10:30:00Z"
  - version: "1.0.0"
TTL: 30 seconds
Usage: Service health monitoring
```

### Metrics Collection
```
metrics:{service}:{metric_name}:{timestamp}
Type: String or Hash
Usage: Time-series metrics for monitoring
```

This Redis schema provides high-performance caching, real-time data management, and scalable event processing for the UA Card Battle Game microservices architecture.