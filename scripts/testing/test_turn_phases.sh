#!/bin/bash

# Test Turn Phases - Union Arena Game Battle Service
# This script tests the complete turn phase flow with proper token switching

set -e

echo "=== Testing Turn Phase Flow ==="

# Configuration
BASE_URL="http://localhost:8004/api/v1"
GAME_BATTLE_URL="http://localhost:8004"

# Generate test tokens
echo "Generating test tokens..."
cd /c/Users/weilo/Desktop/ua
TOKEN_OUTPUT=$(go run scripts/testing/generate_test_token.go)

# Extract tokens and IDs from output
PLAYER1_ID=$(echo "$TOKEN_OUTPUT" | grep "Player 1 ID:" | cut -d' ' -f4)
PLAYER1_TOKEN=$(echo "$TOKEN_OUTPUT" | grep "Player 1 Access Token:" | cut -d' ' -f5)
PLAYER2_ID=$(echo "$TOKEN_OUTPUT" | grep "Player 2 ID:" | cut -d' ' -f4)
PLAYER2_TOKEN=$(echo "$TOKEN_OUTPUT" | grep "Player 2 Access Token:" | cut -d' ' -f5)

echo "Player 1 ID: $PLAYER1_ID"
echo "Player 1 Token: ${PLAYER1_TOKEN:0:20}..."
echo "Player 2 ID: $PLAYER2_ID" 
echo "Player 2 Token: ${PLAYER2_TOKEN:0:20}..."

# Function to make API calls with proper token
make_api_call() {
    local method=$1
    local endpoint=$2
    local token=$3
    local data=$4
    local description=$5
    
    echo -e "\n--- $description ---"
    
    if [ -z "$data" ]; then
        curl -X $method \
            -H "Authorization: Bearer $token" \
            -H "Content-Type: application/json" \
            "$BASE_URL$endpoint" | jq '.'
    else
        curl -X $method \
            -H "Authorization: Bearer $token" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint" | jq '.'
    fi
}

# Function to check game state
check_game_state() {
    local token=$1
    local player_name=$2
    echo -e "\n--- Checking game state ($player_name) ---"
    
    curl -X GET \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        "$BASE_URL/games/$GAME_ID" | jq '{
            game: {
                id: .data.game.id,
                status: .data.game.status,
                current_turn: .data.game.current_turn,
                phase: .data.game.phase,
                active_player: .data.game.active_player
            },
            game_state: {
                turn: .data.game_state.turn,
                phase: .data.game_state.phase,
                active_player: .data.game_state.active_player
            }
        }'
}

# Load test deck data
echo -e "\n=== Loading test deck data ==="
TEST_DECK=$(cat test_data/FULL_50_CARDS_DECK.json)

# 1. Create Game (Public endpoint - no token needed)
echo -e "\n=== Step 1: Creating Game ==="
CREATE_RESPONSE=$(curl -X POST \
    -H "Content-Type: application/json" \
    -d "{
        \"player1_id\": \"$PLAYER1_ID\",
        \"player2_id\": \"$PLAYER2_ID\",
        \"game_mode\": \"standard\",
        \"player1_deck\": $TEST_DECK,
        \"player2_deck\": $TEST_DECK
    }" \
    "$BASE_URL/games")

echo "$CREATE_RESPONSE" | jq '.'
GAME_ID=$(echo "$CREATE_RESPONSE" | jq -r '.data.game.id')
echo "Game ID: $GAME_ID"

# 2. Player 1 joins game
make_api_call "POST" "/games/$GAME_ID/join" "$PLAYER1_TOKEN" "" "Player 1 joins game"

# 3. Player 2 joins game (should start the game)
make_api_call "POST" "/games/$GAME_ID/join" "$PLAYER2_TOKEN" "" "Player 2 joins game (game should start)"

# 4. Both players perform mulligan
make_api_call "POST" "/games/$GAME_ID/mulligan" "$PLAYER1_TOKEN" '{"mulligan": false}' "Player 1 mulligan (keep hand)"
make_api_call "POST" "/games/$GAME_ID/mulligan" "$PLAYER2_TOKEN" '{"mulligan": false}' "Player 2 mulligan (keep hand)"

# 5. Check initial game state
check_game_state "$PLAYER1_TOKEN" "Player 1"

# Get the active player from game state to determine who goes first
ACTIVE_PLAYER_ID=$(curl -s -X GET \
    -H "Authorization: Bearer $PLAYER1_TOKEN" \
    "$BASE_URL/games/$GAME_ID" | jq -r '.data.game_state.active_player')

echo -e "\n=== Active Player: $ACTIVE_PLAYER_ID ==="

# Determine which token to use based on active player
if [ "$ACTIVE_PLAYER_ID" = "$PLAYER1_ID" ]; then
    CURRENT_TOKEN="$PLAYER1_TOKEN"
    CURRENT_PLAYER="Player 1"
    OTHER_TOKEN="$PLAYER2_TOKEN"
    OTHER_PLAYER="Player 2"
else
    CURRENT_TOKEN="$PLAYER2_TOKEN"
    CURRENT_PLAYER="Player 2"
    OTHER_TOKEN="$PLAYER1_TOKEN"
    OTHER_PLAYER="Player 1"
fi

echo "Current turn: $CURRENT_PLAYER"

# 6. Test Start Phase - Draw Card (should be automatic, but test manual draw)
echo -e "\n=== Testing Start Phase - Draw Card ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "DRAW_CARD"}' "$CURRENT_PLAYER attempts to draw card"

# 7. Test wrong player trying to act
echo -e "\n=== Testing Turn Validation (Wrong Player) ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$OTHER_TOKEN" '{"action_type": "DRAW_CARD"}' "$OTHER_PLAYER attempts to act (should fail - not their turn)"

# 8. Test Start Phase - Extra Draw
echo -e "\n=== Testing Start Phase - Extra Draw ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "EXTRA_DRAW"}' "$CURRENT_PLAYER attempts extra draw (costs 1 AP)"

# 9. Test invalid phase actions
echo -e "\n=== Testing Invalid Phase Actions ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "PLAY_CHARACTER"}' "$CURRENT_PLAYER attempts to play card in Start Phase (should fail)"

# 10. Advance to Move Phase
echo -e "\n=== Testing Phase Advancement to Move Phase ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "END_PHASE"}' "$CURRENT_PLAYER ends Start Phase"

check_game_state "$CURRENT_TOKEN" "$CURRENT_PLAYER"

# 11. Test Move Phase
echo -e "\n=== Testing Move Phase ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "MOVE_CHARACTER"}' "$CURRENT_PLAYER attempts to move character"

# 12. Advance to Main Phase
echo -e "\n=== Testing Phase Advancement to Main Phase ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "END_PHASE"}' "$CURRENT_PLAYER ends Move Phase"

check_game_state "$CURRENT_TOKEN" "$CURRENT_PLAYER"

# 13. Test Main Phase - Play Card (would need specific card data)
echo -e "\n=== Testing Main Phase - Play Card ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "PLAY_CHARACTER"}' "$CURRENT_PLAYER attempts to play character card"

# 14. Advance to Attack Phase
echo -e "\n=== Testing Phase Advancement to Attack Phase ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "END_PHASE"}' "$CURRENT_PLAYER ends Main Phase"

check_game_state "$CURRENT_TOKEN" "$CURRENT_PLAYER"

# 15. Test Attack Phase
echo -e "\n=== Testing Attack Phase ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "ATTACK"}' "$CURRENT_PLAYER attempts to attack"

# 16. Advance to End Phase
echo -e "\n=== Testing Phase Advancement to End Phase ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "END_PHASE"}' "$CURRENT_PLAYER ends Attack Phase"

check_game_state "$CURRENT_TOKEN" "$CURRENT_PLAYER"

# 17. End Turn (should switch to other player)
echo -e "\n=== Testing Turn End (Should Switch Players) ==="
make_api_call "POST" "/games/$GAME_ID/actions" "$CURRENT_TOKEN" '{"action_type": "END_TURN"}' "$CURRENT_PLAYER ends turn"

check_game_state "$OTHER_TOKEN" "$OTHER_PLAYER"

# 18. Verify turn switched
echo -e "\n=== Verifying Turn Switch ==="
NEW_ACTIVE_PLAYER=$(curl -s -X GET \
    -H "Authorization: Bearer $OTHER_TOKEN" \
    "$BASE_URL/games/$GAME_ID" | jq -r '.data.game_state.active_player')

if [ "$NEW_ACTIVE_PLAYER" != "$ACTIVE_PLAYER_ID" ]; then
    echo "✅ Turn successfully switched to other player!"
    echo "New active player: $NEW_ACTIVE_PLAYER"
else
    echo "❌ Turn did not switch properly"
fi

# 19. Test Redis game info endpoint
echo -e "\n=== Testing Redis Game Info Endpoint ==="
curl -X GET "$GAME_BATTLE_URL/api/v1/game-info/$GAME_ID" | jq '.'

echo -e "\n=== Turn Phase Flow Test Complete ==="
echo "Game ID for further testing: $GAME_ID"
echo "Player 1 ID: $PLAYER1_ID (Token: ${PLAYER1_TOKEN:0:20}...)"
echo "Player 2 ID: $PLAYER2_ID (Token: ${PLAYER2_TOKEN:0:20}...)"