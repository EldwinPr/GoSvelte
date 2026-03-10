package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"sort"
	"time"

	"gorm.io/gorm"
)

type WishlistService struct {
	BaseService
	Repo           *repositories.BaseRepository[models.Wishlist]
	AccountService *AccountService
}

func NewWishlistService(db *gorm.DB, accountService *AccountService) *WishlistService {
	return &WishlistService{
		BaseService:    BaseService{DB: db},
		Repo:           &repositories.BaseRepository[models.Wishlist]{DB: db},
		AccountService: accountService,
	}
}

// CreateItem adds a new item to the wishlist.
func (s *WishlistService) CreateItem(tenantID string, item *models.Wishlist) error {
	item.ID = repositories.GenerateULID()
	item.TenantID = tenantID
	item.Status = "Active"
	return s.Repo.Create(item)
}

// GetRankedWishlist retrieves items and calculates their 'Smart Score'.
func (s *WishlistService) GetRankedWishlist(tenantID string) ([]models.Wishlist, error) {
	// 1. Fetch all active items
	items, err := s.Repo.FindAll(s.Repo.FilterScope("tenant_id", "=", tenantID), s.Repo.FilterScope("status", "=", "Active"))
	if err != nil {
		return nil, err
	}

	// 2. Fetch current Total Net Worth
	summary, err := s.AccountService.GetDashboardSummary(tenantID)
	if err != nil {
		return nil, err
	}
	netWorth := summary["net_worth"].(float64)

	// 3. Calculate Scores for each item
	for i := range items {
		item := &items[i]

		// Financial Impact (Ability to pay)
		// 100% impact means price is >= Net Worth
		// 0% impact means price is negligible
		if netWorth > 0 {
			item.FinancialImpact = (item.TargetAmount / netWorth) * 100
		} else {
			item.FinancialImpact = 100 // Very high impact if no money
		}

		// Calculate Smart Score (Weighted Average)
		// We want high Need, high Want, high Productivity, but LOW Financial Impact
		
		// Convert Financial Impact to a 1-6 scale (inverted: low impact = 6, high = 1)
		financialScore := 6.0 - (item.FinancialImpact / 100 * 6)
		if financialScore < 1 {
			financialScore = 1
		}

		// Weighting: Need (40%), Want (20%), Productivity (20%), Financial Ability (20%)
		item.Score = (float64(item.Need) * 0.4) + 
					 (float64(item.Want) * 0.2) + 
					 (float64(item.Productivity) * 0.2) + 
					 (financialScore * 0.2)
	}

	// 4. Sort by Score descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].Score > items[j].Score
	})

	return items, nil
}

// UpdateProgress allows manual allocation of funds toward an item.
func (s *WishlistService) UpdateProgress(tenantID string, itemID string, amount float64) error {
	return s.DB.Model(&models.Wishlist{}).
		Where("id = ? AND tenant_id = ?", itemID, tenantID).
		Update("current_saved", amount).Error
}

// MarkPurchased records the purchase as a transaction and archives the item.
func (s *WishlistService) MarkPurchased(tenantID string, itemID string, accountID string) error {
	return s.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch the item
		var item models.Wishlist
		if err := tx.First(&item, "id = ? AND tenant_id = ?", itemID, tenantID).Error; err != nil {
			return errors.New("item not found")
		}

		// 2. Create Transaction (Expense)
		transaction := &models.Transaction{
			ID:          repositories.GenerateULID(),
			TenantID:    tenantID,
			AccountID:   accountID,
			Amount:      item.TargetAmount,
			Description: "Wishlist Purchase: " + item.Name,
			Date:        time.Now(),
			Type:        "Expense",
		}

		// Update Account Balance
		var account models.Account
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&account, "id = ? AND tenant_id = ?", accountID, tenantID).Error; err != nil {
			return err
		}
		account.Balance -= transaction.Amount
		if err := tx.Save(&account).Error; err != nil {
			return err
		}

		// Save Transaction
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 3. Update Wishlist Item status
		item.Status = "Purchased"
		return tx.Save(&item).Error
	})
}

// DeleteItem removes an item from the wishlist.
func (s *WishlistService) DeleteItem(tenantID string, id string) error {
	var item models.Wishlist
	if err := s.DB.Where("id = ? AND tenant_id = ?", id, tenantID).First(&item).Error; err != nil {
		return errors.New("item not found")
	}
	return s.Repo.Delete(&item)
}
