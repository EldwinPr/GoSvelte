package app

import (
	"context"
	"errors"
	"gosvelte/app/models"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type contextKey string

const UserContextKey = "user"

type Middleware struct {
	DB *gorm.DB
}

func NewMiddleware(db *gorm.DB) *Middleware {
	return &Middleware{DB: db}
}

// Auth is a simple middleware to simulate authentication.
// In production, this would verify a JWT or session cookie.
func (m *Middleware) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized: No Session", http.StatusUnauthorized)
			return
		}

		userID := cookie.Value
		var user models.User
		if err := m.DB.First(&user, "id = ?", userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If user not found, clear the cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    "",
					Path:     "/",
					Expires:  time.Unix(0, 0),
					HttpOnly: true,
				})
			}
			http.Error(w, "Unauthorized: Session invalid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// RequireClearance ensures the user has at least the required level.
// Level 20 (Developer) is granted access to everything.
func (m *Middleware) RequireClearance(minLevel int, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(*models.User)
		if !ok {
			http.Error(w, "Internal Server Error: User context missing", http.StatusInternalServerError)
			return
		}

		// Developer (20) always has access.
		if user.Clearance < minLevel && user.Clearance != 20 {
			http.Error(w, "Forbidden: Insufficient Clearance Level "+string(rune(minLevel)), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}
