package auth

import (
	"context"
	"errors"

	serviceauth "github.com/dockloud/service-auth/core"
)

var ErrValidationFailed = errors.New("validation failed")

func IdentifyCallerService(ctx context.Context, expected serviceauth.Principal) error {
	actual, ok := serviceauth.PrincipalFromContext(ctx)
	if !ok || actual != expected {
		return ErrValidationFailed
	}
	return nil
}
