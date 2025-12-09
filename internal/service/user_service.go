package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	apperror "github.com/hafiztri123/kki-be/internal/app_error"
	"github.com/hafiztri123/kki-be/internal/constants"
	"github.com/hafiztri123/kki-be/internal/dto"
	"github.com/hafiztri123/kki-be/internal/models"
	"github.com/hafiztri123/kki-be/internal/repository"
	"github.com/hafiztri123/kki-be/internal/utils"
	"golang.org/x/crypto/bcrypt"
)


type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}


func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	return s.userRepo.InsertUser(
		ctx,
		&models.User{
			ID: uuid.New(),
			Username: req.Username,
			Password: string(hashedPassword),
			Email: req.Email,
			Role: req.Role,
			Name: req.Name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{},

		},
	)
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	fetchedUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(fetchedUser.ID, fetchedUser.Role)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
	}, nil
}

func (s *UserService) CreateCashier(ctx context.Context, req *dto.CreateCashierRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.InsertUser(
		ctx,
		&models.User{
			ID:        uuid.New(),
			Username:  req.Username,
			Password:  string(hashedPassword),
			Email:     req.Email,
			Role:      constants.RoleCashier,
			Name:      req.Name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{},
		},
	)
}

func (s *UserService) GetCashiers(ctx context.Context, limit, offset int) ([]dto.UserResponse, int64, error) {
	users, totalCount, err := s.userRepo.GetUsersByRole(ctx, "cashier", limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		})
	}

	return responses, totalCount, nil
}

func (s *UserService) GetCashierByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.Role != "cashier" {
		return nil, apperror.ErrNotFound
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserService) UpdateCashier(ctx context.Context, id string, req *dto.UpdateCashierRequest) error {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if user.Role != "cashier" {
		return apperror.ErrNotFound
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Username = req.Username
	user.Email = req.Email
	user.Name = req.Name
	user.Password = string(newHashedPassword)

	user.UpdatedAt = time.Now()

	return s.userRepo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteCashier(ctx context.Context, id string) error {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if user.Role != "cashier" {
		return apperror.ErrNotFound
	}

	return s.userRepo.DeleteUser(ctx, id)
}


