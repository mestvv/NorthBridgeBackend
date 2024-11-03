package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/mestvv/NorthBridgeBackend/internal/config"
	"github.com/mestvv/NorthBridgeBackend/internal/repository"
	"github.com/mestvv/NorthBridgeBackend/pkg/auth"
	"github.com/mestvv/NorthBridgeBackend/pkg/email"
	"github.com/mestvv/NorthBridgeBackend/pkg/hash"
	"github.com/mestvv/NorthBridgeBackend/pkg/otp"
)

type Services struct {
	Users Users
}

type Deps struct {
	Logger       *slog.Logger
	Config       *config.Config
	Hasher       hash.PasswordHasher
	TokenManager auth.TokenManager
	OtpGenerator otp.Generator
	EmailSender  email.Sender
	Repos        *repository.Repositories
}

func NewServices(deps Deps) *Services {
	emailService := newEmailsService(deps.EmailSender, deps.Config.Email)

	return &Services{
		Users: newUserService(deps.Repos.Users,
			deps.Repos.RefreshSession,
			deps.Hasher,
			deps.TokenManager,
			deps.OtpGenerator,
			emailService,
			deps.Config.Auth,
		),
	}
}

type Emails interface {
	SendUserVerificationEmail(VerificationEmailInput) error
}

type Users interface {
	Register(ctx context.Context, input *UserRegisterInput) error
	Auth(ctx context.Context, input *UserAuthInput) (*Tokens, error)
	createSession(ctx context.Context, userID *uuid.UUID, userAgent *string, userIP *string) (*Tokens, error)
}
