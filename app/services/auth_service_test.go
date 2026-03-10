package services

import (
	"gosvelte/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_Register(t *testing.T) {
	db := setupTestDB()
	service := NewAuthService(db)

	t.Run("Successful Registration", func(t *testing.T) {
		user, err := service.Register("Test Tenant", "John Doe", "john@example.com", "password123")
		
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.NotEmpty(t, user.TenantID)

		// Verify database records
		var count int64
		db.Model(&models.Tenant{}).Where("id = ?", user.TenantID).Count(&count)
		assert.Equal(t, int64(1), count)

		db.Model(&models.UserSettings{}).Where("user_id = ?", user.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Duplicate Email Registration", func(t *testing.T) {
		db := setupTestDB()
		service := NewAuthService(db)
		service.Register("T1", "U1", "john@example.com", "p")
		
		_, err := service.Register("T2", "User 2", "john@example.com", "pass")
		assert.Error(t, err)
	})
}

func TestAuthService_Login(t *testing.T) {
	db := setupTestDB()
	service := NewAuthService(db)

	// Pre-register a user
	service.Register("Tenant", "Jane Doe", "jane@example.com", "secret")

	t.Run("Successful Login", func(t *testing.T) {
		user, err := service.Login("jane@example.com", "secret")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "Jane Doe", user.Name)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		user, err := service.Login("jane@example.com", "wrong")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("Non-existent User", func(t *testing.T) {
		user, err := service.Login("nobody@example.com", "secret")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestAuthService_Sessions(t *testing.T) {
	t.Run("Create and Validate Session", func(t *testing.T) {
		db := setupTestDB()
		service := NewAuthService(db)
		user, _ := service.Register("T", "U", "u@e.com", "p")

		session, err := service.CreateSession(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, session)

		// Validate
		validSession, err := service.ValidateSession(session.ID)
		assert.NoError(t, err)
		assert.NotNil(t, validSession)
		assert.Equal(t, user.ID, validSession.UserID)
	})

	t.Run("Logout revokes session", func(t *testing.T) {
		db := setupTestDB()
		service := NewAuthService(db)
		user, _ := service.Register("T", "U", "u@e.com", "p")
		session, _ := service.CreateSession(user.ID)
		
		err := service.Logout(session.ID)
		assert.NoError(t, err)

		// Should now be invalid
		_, err = service.ValidateSession(session.ID)
		assert.Error(t, err)
		assert.Equal(t, "session revoked", err.Error())
	})
}
