package services

import (
	"errors"
	"gosvelte/app/models"
	"gosvelte/app/repositories"
	"time"

	"gorm.io/gorm"
)

type RecurringPaymentService struct {
	BaseService
	Repo               *repositories.BaseRepository[models.RecurringPayment]
	InstanceRepo       *repositories.BaseRepository[models.RecurringInstance]
	TransactionService *TransactionService
}

func NewRecurringPaymentService(db *gorm.DB, txService *TransactionService) *RecurringPaymentService {
	return &RecurringPaymentService{
		BaseService: BaseService{DB: db},
		Repo: &repositories.BaseRepository[models.RecurringPayment]{DB: db},
		InstanceRepo: &repositories.BaseRepository[models.RecurringInstance]{DB: db},
		TransactionService: txService,
	}
}

// CreateSchedule sets up a new recurring payment schedule.
func (s *RecurringPaymentService) CreateSchedule(tenantID string, schedule *models.RecurringPayment) error {
	schedule.ID = repositories.GenerateULID()
	schedule.TenantID = tenantID
	
	if schedule.NextDate.IsZero() {
		schedule.NextDate = schedule.StartDate
	}

	return s.Repo.Create(schedule)
}

// GetSchedules retrieves all schedules for a tenant.
func (s *RecurringPaymentService) GetSchedules(tenantID string) ([]models.RecurringPayment, error) {
	return s.Repo.FindAll(s.Repo.FilterScope("tenant_id", "=", tenantID))
}

// ToggleSchedule activates or deactivates a schedule.
func (s *RecurringPaymentService) ToggleSchedule(tenantID string, id string, active bool) error {
	return s.DB.Model(&models.RecurringPayment{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Update("is_active", active).Error
}

// ProcessPending generates 'pending' instances for all due schedules.
// This does NOT create transactions or update account balances.
func (s *RecurringPaymentService) ProcessPending(tenantID string) (int, error) {
	var schedules []models.RecurringPayment
	// Find active schedules where NextDate is today or in the past
	err := s.DB.Where("tenant_id = ? AND is_active = ? AND next_date <= ?", tenantID, true, time.Now()).
		Find(&schedules).Error
	if err != nil {
		return 0, err
	}

	generatedCount := 0
	for _, schedule := range schedules {
		err := s.Transaction(func(tx *gorm.DB) error {
			// 1. Create a 'pending' instance
			instance := &models.RecurringInstance{
				ID:                 repositories.GenerateULID(),
				RecurringPaymentID: schedule.ID,
				DueDate:            schedule.NextDate,
				Status:             "pending",
			}
			if err := tx.Create(instance).Error; err != nil {
				return err
			}

			// 2. Advance the NextDate on the schedule
			nextDate := s.calculateNextDate(schedule.NextDate, schedule.Frequency)
			return tx.Model(&schedule).Update("next_date", nextDate).Error
		})

		if err == nil {
			generatedCount++
		}
	}

	return generatedCount, nil
}

// GetPendingInstances retrieves all instances that require manual confirmation.
func (s *RecurringPaymentService) GetPendingInstances(tenantID string) ([]models.RecurringInstance, error) {
	var instances []models.RecurringInstance
	err := s.DB.Preload("RecurringPayment").
		Joins("JOIN recurring_payments ON recurring_payments.id = recurring_instances.recurring_payment_id").
		Where("recurring_payments.tenant_id = ? AND recurring_instances.status = ?", tenantID, "pending").
		Find(&instances).Error
	return instances, err
}

// ConfirmInstance manually processes a pending instance, creating a real transaction and updating balance.
func (s *RecurringPaymentService) ConfirmInstance(tenantID string, instanceID string, actualAmount float64) error {
	return s.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch the instance and its schedule
		var instance models.RecurringInstance
		if err := tx.Preload("RecurringPayment").First(&instance, "id = ?", instanceID).Error; err != nil {
			return errors.New("instance not found")
		}

		// Verify ownership
		if instance.RecurringPayment.TenantID != tenantID {
			return errors.New("unauthorized")
		}

		if instance.Status != "pending" {
			return errors.New("instance is already processed or skipped")
		}

		// 2. Create the real Transaction (using TransactionService logic)
		transaction := &models.Transaction{
			ID:                 repositories.GenerateULID(),
			TenantID:           tenantID,
			AccountID:          *instance.RecurringPayment.AccountID,
			CategoryID:         instance.RecurringPayment.CategoryID,
			RecurringPaymentID: &instance.RecurringPaymentID,
			Amount:             actualAmount,
			Description:        instance.RecurringPayment.Description,
			Date:               time.Now(), // Date of actual confirmation
			Type:               "Expense",
		}

		// Lock and Update Account
		var account models.Account
		if err := tx.Clauses(gorm.Expr("FOR UPDATE")).First(&account, "id = ? AND tenant_id = ?", transaction.AccountID, tenantID).Error; err != nil {
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

		// 3. Update Instance status and link the transaction
		instance.Status = "processed"
		instance.TransactionID = &transaction.ID
		return tx.Save(&instance).Error
	})
}

// SkipInstance marks a pending instance as 'skipped'.
func (s *RecurringPaymentService) SkipInstance(tenantID string, instanceID string) error {
	var instance models.RecurringInstance
	if err := s.DB.Preload("RecurringPayment").First(&instance, "id = ?", instanceID).Error; err != nil {
		return err
	}

	if instance.RecurringPayment.TenantID != tenantID {
		return errors.New("unauthorized")
	}

	return s.DB.Model(&instance).Update("status", "skipped").Error
}

func (s *RecurringPaymentService) calculateNextDate(current time.Time, frequency string) time.Time {
	switch frequency {
	case "Daily":
		return current.AddDate(0, 0, 1)
	case "Weekly":
		return current.AddDate(0, 0, 7)
	case "Bi-Weekly":
		return current.AddDate(0, 0, 14)
	case "Monthly":
		return current.AddDate(0, 1, 0)
	case "Yearly":
		return current.AddDate(1, 0, 0)
	default:
		return current.AddDate(0, 1, 0)
	}
}
