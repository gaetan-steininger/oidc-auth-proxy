package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate state: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func extractUserGroupsClaims(claims map[string]interface{}, user_claim string, groups_claim string) (string, []string, error) {
	user_claim_interface_value, exists := claims[user_claim]
	if !exists {
		return "", []string{}, fmt.Errorf("User claim %s not found", user_claim)
	}

	switch user_claim_value := user_claim_interface_value.(type) {
	case string:
		if groups_claim == "" {
			return user_claim_value, []string{}, nil
		}

		groups_claim_interface_value, exists := claims[groups_claim]
		if !exists {
			return "", []string{}, fmt.Errorf("Groups claim %s not found", groups_claim)
		}

		switch v := groups_claim_interface_value.(type) {
		case []interface{}:
			groups_claim_value := make([]string, len(v))
			for i, j := range v {
				if str, ok := j.(string); ok {
					groups_claim_value[i] = str
				} else {
					return user_claim_value, []string{}, fmt.Errorf("Invalid type in groups claim : %T", j)
				}
			}
			return user_claim_value, groups_claim_value, nil
		default:
			return user_claim_value, []string{}, fmt.Errorf("Invalid type for groups claim : %T", v)
		}
	default:
		return "", []string{}, fmt.Errorf("Invalid type for user claim")
	}
}
