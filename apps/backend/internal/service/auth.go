package service

import (
	"context"
	"fmt"

	"github.com/sriniously/tasker/internal/server"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
)

type AuthService struct {
	server *server.Server
}

func NewAuthService(s *server.Server) *AuthService {
	clerk.SetKey(s.Config.Auth.SecretKey)
	return &AuthService{
		server: s,
	}
}

func (s *AuthService) GetUserEmail(ctx context.Context, userID string) (string, error) {
	user, err := clerkUser.Get(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user from Clerk: %w", err)
	}

	if len(user.EmailAddresses) == 0 {
		return "", fmt.Errorf("user %s has no email addresses", userID)
	}

	for _, email := range user.EmailAddresses {
		if user.PrimaryEmailAddressID != nil && email.ID == *user.PrimaryEmailAddressID {
			return email.EmailAddress, nil
		}
	}

	return user.EmailAddresses[0].EmailAddress, nil
}
