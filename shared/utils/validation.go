package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"ua/shared/models"
)

var (
	emailRegex    = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(strings.ToLower(email))
}

func IsValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

func IsValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 128
}

func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSecureToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

func ValidatePageAndLimit(page, limit int) (int, int, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		return 0, 0, fmt.Errorf("limit cannot exceed 100")
	}

	return page, limit, nil
}

// DeckValidationError represents deck validation errors
type DeckValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e DeckValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateDeck validates a deck according to Union Arena rules
func ValidateDeck(cards []models.Card) []DeckValidationError {
	var errors []DeckValidationError

	// Check deck size (50 cards total)
	if len(cards) != 50 {
		errors = append(errors, DeckValidationError{
			Field:   "deck_size",
			Message: fmt.Sprintf("Deck must contain exactly 50 cards, found %d", len(cards)),
		})
	}

	// Count AP cards
	apCount := 0
	colorCounts := make(map[string]int)
	cardCopies := make(map[string]int)

	for _, card := range cards {
		// Count AP cards
		if card.CardType == models.CardTypeAP {
			apCount++
		}

		// Count colors
		if card.Color != "" {
			colorCounts[card.Color]++
		}

		// Count card copies
		cardCopies[card.CardNumber]++
	}

	// Check AP card count (must be exactly 3)
	if apCount != 3 {
		errors = append(errors, DeckValidationError{
			Field:   "ap_cards",
			Message: fmt.Sprintf("Deck must contain exactly 3 AP cards, found %d", apCount),
		})
	}

	// Check single color requirement
	var deckColor string
	colorCount := 0
	for color, count := range colorCounts {
		if count > 0 {
			colorCount++
			if deckColor == "" {
				deckColor = color
			}
		}
	}

	if colorCount > 1 {
		errors = append(errors, DeckValidationError{
			Field:   "deck_color",
			Message: "Deck must contain cards of only one color",
		})
	}

	// Check card copy limits (usually max 4 copies per card)
	for cardNumber, copies := range cardCopies {
		if copies > 4 {
			errors = append(errors, DeckValidationError{
				Field:   "card_copies",
				Message: fmt.Sprintf("Card %s has %d copies, maximum allowed is 4", cardNumber, copies),
			})
		}
	}

	return errors
}

// ValidateColorEffect checks if a color effect is valid
func ValidateColorEffect(color string) bool {
	colorEffects := models.GetColorEffects()
	_, exists := colorEffects[color]
	return exists
}

// ValidateTriggerEffect checks if a trigger effect is valid
func ValidateTriggerEffect(triggerEffect string) bool {
	validEffects := []string{
		models.TriggerEffectDrawCard,
		models.TriggerEffectColor,
		models.TriggerEffectActiveBP3000,
		models.TriggerEffectAddToHand,
		models.TriggerEffectRushOrAddToHand,
		models.TriggerEffectSpecial,
		models.TriggerEffectFinal,
		models.TriggerEffectNil,
	}

	for _, effect := range validEffects {
		if effect == triggerEffect {
			return true
		}
	}
	return false
}
