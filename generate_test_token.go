package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"ua/shared/auth"
)

func main() {
	// Use the same JWT secret as in docker-compose.yml
	jwtSecret := "your-super-secret-jwt-key-change-in-production"
	
	// Generate test player IDs
	player1ID := uuid.New()
	player2ID := uuid.New()
	
	// Generate JWT tokens for both test players
	token1, err := auth.GenerateTokenPair(player1ID, "test_player_1", jwtSecret)
	if err != nil {
		log.Fatalf("Failed to generate token for player 1: %v", err)
	}
	
	token2, err := auth.GenerateTokenPair(player2ID, "test_player_2", jwtSecret)
	if err != nil {
		log.Fatalf("Failed to generate token for player 2: %v", err)
	}
	
	fmt.Printf("Player 1 ID: %s\n", player1ID.String())
	fmt.Printf("Player 1 Access Token: %s\n\n", token1.AccessToken)
	
	fmt.Printf("Player 2 ID: %s\n", player2ID.String())
	fmt.Printf("Player 2 Access Token: %s\n\n", token2.AccessToken)
	
	fmt.Println("Export these for your integration tests:")
	fmt.Printf("export PLAYER1_ID='%s'\n", player1ID.String())
	fmt.Printf("export PLAYER1_TOKEN='%s'\n", token1.AccessToken)
	fmt.Printf("export PLAYER2_ID='%s'\n", player2ID.String())
	fmt.Printf("export PLAYER2_TOKEN='%s'\n", token2.AccessToken)
}