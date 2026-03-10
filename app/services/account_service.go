package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"

	"gorm.io/gorm"
)

type AccountService struct {
	BaseService
	Repo *repositories.AccountRepository
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{
		BaseService: BaseService{DB: db},
		Repo: &repositories.AccountRepository{
			BaseRepository: repositories.BaseRepository[models.Account]{DB: db},
		},
	}
}

// CreateAccount initializes a new account for a tenant.
func (s *AccountService) CreateAccount(tenantID string, account *models.Account) error {
	account.ID = repositories.GenerateULID()
	account.TenantID = tenantID
	return s.Repo.Create(account)
}

// GetAccounts retrieves all accounts for a tenant.
func (s *AccountService) GetAccounts(tenantID string) ([]models.Account, error) {
	return s.Repo.FindAll(s.Repo.FilterScope("tenant_id", "=", tenantID))
}

// GetAccount retrieves a specific account by ID, ensuring it belongs to the tenant.
func (s *AccountService) GetAccount(tenantID string, id string) (*models.Account, error) {
	var account models.Account
	err := s.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &account, nil
}

// UpdateAccount updates an existing account after verifying ownership.
func (s *AccountService) UpdateAccount(tenantID string, account *models.Account) error {
	// Verify ownership first
	var existing models.Account
	if err := s.DB.Where("id = ? AND tenant_id = ?", account.ID, tenantID).First(&existing).Error; err != nil {
		return errors.New("unauthorized or account not found")
	}

	// Only update editable fields (prevent TenantID manipulation)
	existing.Name = account.Name
	existing.Type = account.Type
	existing.Currency = account.Currency
	// Balance is typically updated via transactions, but we allow it here for manual corrections if needed.
	// In a more strict system, we'd exclude Balance from manual updates.
	existing.Balance = account.Balance

	return s.Repo.Update(&existing)
}

// DeleteAccount soft-deletes an account after verifying ownership.
func (s *AccountService) DeleteAccount(tenantID string, id string) error {
	var account models.Account
	if err := s.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&account).Error; err != nil {
		return errors.New("unauthorized or account not found")
	}
	return s.Repo.Delete(&account)
}

// GetDashboardSummary calculates the net worth and account breakdown for a tenant.
func (s *AccountService) GetDashboardSummary(tenantID string) (map[string]interface{}, error) {
	accounts, err := s.GetAccounts(tenantID)
	if err != nil {
		return nil, err
	}

	var netWorth float64
	breakdown := make(map[string]float64)

	for _, acc := range accounts {
		netWorth += acc.Balance
		breakdown[acc.Type] += acc.Balance
	}

	return map[string]interface{}{
		"net_worth": netWorth,
		"breakdown": breakdown,
		"accounts":  accounts,
	}, nil
}
