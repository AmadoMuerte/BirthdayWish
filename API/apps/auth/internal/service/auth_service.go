package service

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/auth/internal/models"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/auth/internal/storage"
	authProto "github.com/AmadoMuerte/BirthdayWish/API/proto/auth"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage   *storage.Storage
	log       *slog.Logger
	secretKey string
	authProto.UnimplementedAuthServiceServer
}

func NewAuthService(storage *storage.Storage, log *slog.Logger, secretKey string) *AuthService {
	return &AuthService{
		storage:   storage,
		log:       log,
		secretKey: secretKey,
	}
}

func (s *AuthService) validateCredentials(validateType string, data string) error {
	var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()-_+=]{8,20}$`)
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	switch validateType {
	case "username":
		if !usernameRegex.MatchString(data) {
			return errors.New("invalid username")
		}
	case "password":
		if !passwordRegex.MatchString(data) {
			return errors.New("invalid password")
		}
	case "email":
		if !emailRegex.MatchString(data) {
			return errors.New("invalid email")
		}
	}
	return nil
}

func (s *AuthService) SignUp(ctx context.Context, req *authProto.SignUpRequest) (*authProto.SignUpResponse, error) {
	if err := s.validateCredentials("username", req.Username); err != nil {
		return nil, err
	}
	if err := s.validateCredentials("password", req.Password); err != nil {
		return nil, err
	}
	if err := s.validateCredentials("email", req.Email); err != nil {
		return nil, err
	}

	exists, err := s.storage.UserExists(ctx, req.Username, req.Email)
	if err != nil {
		s.log.Error("database error checking user existence", "error", err)
		return nil, errors.New("internal server error")
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, errors.New("internal server error")
	}

	user := models.User{
		Name:     req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := s.storage.CreateUser(ctx, &user); err != nil {
		s.log.Error("failed to create user", "error", err)
		return nil, errors.New("internal server error")
	}

	return &authProto.SignUpResponse{
		Message: "User created successfully",
		UserId:  strconv.FormatInt(user.ID, 10),
	}, nil
}

func (s *AuthService) SignIn(ctx context.Context, req *authProto.SignInRequest) (*authProto.SignInResponse, error) {
	if err := s.validateCredentials("username", req.Username); err != nil {
		return nil, errors.New("authentication failed")
	}
	if err := s.validateCredentials("password", req.Password); err != nil {
		return nil, errors.New("authentication failed")
	}

	user, err := s.storage.GetUserByUsername(ctx, req.Username)
	if err != nil {
		s.log.Error("error getting user by username", "error", err)
		return nil, errors.New("authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("authentication failed")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		s.log.Error("failed to generate token", "error", err)
		return nil, errors.New("internal server error")
	}

	expTime := time.Now().Add(24 * time.Hour).Unix()

	return &authProto.SignInResponse{
		Token:  tokenString,
		UserId: strconv.FormatInt(user.ID, 10),
		Email:  user.Email,
		Name:   user.Name,
		Exp:    expTime,
	}, nil
}
