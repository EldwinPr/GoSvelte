package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	BaseService
	Repo        *repositories.UserRepository
	SessionRepo *repositories.BaseRepository[models.Session]
	TenantRepo  *repositories.BaseRepository[models.Tenant]
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		BaseService: BaseService{DB: db},
		Repo: &repositories.UserRepository{
			BaseRepository: repositories.BaseRepository[models.User]{DB: db},
		},
		SessionRepo: &repositories.BaseRepository[models.Session]{DB: db},
		TenantRepo:  &repositories.BaseRepository[models.Tenant]{DB: db},
	}
}

// Register creates a new Tenant and the first User for that tenant.
func (s *AuthService) Register(tenantName, name, email, password string) (*models.User, error) {
	var user models.User

	err := s.Transaction(func(tx *gorm.DB) error {
		// 1. Create Tenant
		tenant := &models.Tenant{
			ID:   repositories.GenerateULID(),
			Name: tenantName,
		}
		if err := tx.Create(tenant).Error; err != nil {
			return err
		}

		// 2. Hash Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// 3. Create User
		user = models.User{
			ID:       repositories.GenerateULID(),
			TenantID: tenant.ID,
			Name:     name,
			Email:    email,
			Password: string(hashedPassword),
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 4. Create Default User Settings
		settings := &models.UserSettings{
			ID:     repositories.GenerateULID(),
			UserID: user.ID,
		}
		if err := tx.Create(settings).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Login authenticates users.
func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// CreateSession generates a new session for a user.
func (s *AuthService) CreateSession(userID string) (*models.Session, error) {
	session := &models.Session{
		ID:        repositories.GenerateULID(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour), // 7 days
	}

	if err := s.SessionRepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

// ValidateSession checks if a session ID is valid and not expired.
func (s *AuthService) ValidateSession(sessionID string) (*models.Session, error) {
	var session models.Session
	err := s.DB.Preload("User").First(&session, "id = ? AND expires_at > ?", sessionID, time.Now()).Error
	if err != nil {
		return nil, errors.New("session invalid or expired")
	}

	if session.RevokedAt != nil {
		return nil, errors.New("session revoked")
	}

	return &session, nil
}

// Logout revokes a session.
func (s *AuthService) Logout(sessionID string) error {
	now := time.Now()
	return s.DB.Model(&models.Session{}).Where("id = ?", sessionID).Update("revoked_at", &now).Error
}

// GetContext retrieves user profile and TenantID.
func (s *AuthService) GetContext(userID string) (*models.User, error) {
	var user models.User
	err := s.DB.Preload("Tenant").First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
