package service

import (
	"errors"
	"time"

	"github.com/triageflow/backend/model"
	"gorm.io/gorm"
)

// Standard queue errors.
var (
	ErrAlreadyEnqueued    = errors.New("task already has an active queue entry")
	ErrQueueEntryNotFound = errors.New("queue entry not found")
	ErrInvalidTransition  = errors.New("invalid queue status transition")
)

// QueueService encapsulates all queue domain logic.
type QueueService struct{}

func NewQueueService() *QueueService {
	return &QueueService{}
}

// ComputeQueueOrder produces a single sortable integer that encodes
// both priority and creation time.
//
// Formula: PriorityWeight * 1_000_000_000 + unix_timestamp
//
// This guarantees:
//   - urgent (weight=0) < high (weight=1) < normal (weight=2)
//   - within the same priority band, earlier timestamps sort first
//   - ties (same second) are broken by auto-increment ID at query time
func ComputeQueueOrder(priority string, createdAt time.Time) int {
	return model.PriorityWeight(priority)*1_000_000_000 + int(createdAt.Unix())
}

// Enqueue creates a QueueEntry for a triaged task.
// It is idempotent: calling it twice for the same task is a no-op.
// The entire operation runs in a transaction to prevent race conditions
// on queue_number generation.
func (s *QueueService) Enqueue(db *gorm.DB, task *model.Task) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Check for existing entry (unique constraint on task_id, but check first to be idempotent)
		var count int64
		tx.Model(&model.QueueEntry{}).Where("task_id = ?", task.ID).Count(&count)
		if count > 0 {
			return nil
		}

		// Generate queue number: next sequential in department for today.
		// Inside a transaction so concurrent inserts don't collide.
		var maxNum int
		tx.Model(&model.QueueEntry{}).
			Where("department = ? AND DATE(created_at) = CURDATE()", task.FinalDepartment).
			Select("COALESCE(MAX(queue_number), 0)").
			Scan(&maxNum)

		entry := model.QueueEntry{
			TaskID:         task.ID,
			PatientName:    task.PatientName,
			Department:     task.FinalDepartment,
			Priority:       task.FinalPriority,
			QueueStatus:    "waiting",
			QueueNumber:    maxNum + 1,
			QueueOrder:     ComputeQueueOrder(task.FinalPriority, task.CreatedAt),
			ChiefComplaint: task.ChiefComplaint,
		}

		return tx.Create(&entry).Error
	})
}

// ListByDepartment returns queue entries for a department, ordered by
// queue_order ASC, id ASC (priority first, then arrival time).
func (s *QueueService) ListByDepartment(db *gorm.DB, department string, status string) ([]model.QueueEntry, error) {
	var entries []model.QueueEntry
	query := db

	if department != "" {
		query = query.Where("department = ?", department)
	}
	if status != "" {
		query = query.Where("queue_status = ?", status)
	}

	err := query.Order("queue_order ASC, id ASC").Find(&entries).Error
	if entries == nil {
		entries = []model.QueueEntry{}
	}
	return entries, err
}

// PositionInfo contains everything a patient needs to see about their queue state.
type PositionInfo struct {
	TaskID      uint               `json:"task_id"`
	QueueNumber int                `json:"queue_number"`
	Department  string             `json:"department"`
	QueueStatus string             `json:"queue_status"`
	Ahead       int                `json:"ahead"`
	NowServing  []model.QueueEntry `json:"now_serving"`
	WaitingList []model.QueueEntry `json:"waiting_list"`
}

// GetPosition returns the patient's queue position plus the department board
// (who is being served and who is waiting), all in one query-set.
func (s *QueueService) GetPosition(db *gorm.DB, taskID uint) (*PositionInfo, error) {
	var entry model.QueueEntry
	if err := db.Where("task_id = ?", taskID).First(&entry).Error; err != nil {
		return nil, ErrQueueEntryNotFound
	}

	// Count waiting patients ahead in the same department
	var ahead int64
	db.Model(&model.QueueEntry{}).
		Where("department = ? AND queue_status = ? AND (queue_order < ? OR (queue_order = ? AND id < ?))",
			entry.Department, "waiting", entry.QueueOrder, entry.QueueOrder, entry.ID).
		Count(&ahead)

	// Fetch currently called patients in this department (now serving)
	var nowServing []model.QueueEntry
	db.Where("department = ? AND queue_status = ?", entry.Department, "called").
		Order("called_at DESC").
		Limit(3).
		Find(&nowServing)
	if nowServing == nil {
		nowServing = []model.QueueEntry{}
	}

	// Fetch next waiting patients in this department
	var waitingList []model.QueueEntry
	db.Where("department = ? AND queue_status = ?", entry.Department, "waiting").
		Order("queue_order ASC, id ASC").
		Limit(8).
		Find(&waitingList)
	if waitingList == nil {
		waitingList = []model.QueueEntry{}
	}

	return &PositionInfo{
		TaskID:      entry.TaskID,
		QueueNumber: entry.QueueNumber,
		Department:  entry.Department,
		QueueStatus: entry.QueueStatus,
		Ahead:       int(ahead),
		NowServing:  nowServing,
		WaitingList: waitingList,
	}, nil
}

// CallPatient transitions waiting → called and syncs task status.
// Returns the updated entry or an error.
func (s *QueueService) CallPatient(db *gorm.DB, taskID uint) (*model.QueueEntry, error) {
	var result *model.QueueEntry

	err := db.Transaction(func(tx *gorm.DB) error {
		var entry model.QueueEntry
		if err := tx.Where("task_id = ?", taskID).First(&entry).Error; err != nil {
			return ErrQueueEntryNotFound
		}

		if entry.QueueStatus != "waiting" {
			return ErrInvalidTransition
		}

		now := time.Now()
		entry.QueueStatus = "called"
		entry.CalledAt = &now
		if err := tx.Save(&entry).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Task{}).Where("id = ?", entry.TaskID).
			Update("status", "in_progress").Error; err != nil {
			return err
		}

		result = &entry
		return nil
	})

	return result, err
}

// CompletePatient transitions called → completed and syncs task status.
// Returns the updated entry or an error.
func (s *QueueService) CompletePatient(db *gorm.DB, taskID uint) (*model.QueueEntry, error) {
	var result *model.QueueEntry

	err := db.Transaction(func(tx *gorm.DB) error {
		var entry model.QueueEntry
		if err := tx.Where("task_id = ?", taskID).First(&entry).Error; err != nil {
			return ErrQueueEntryNotFound
		}

		if entry.QueueStatus != "called" {
			return ErrInvalidTransition
		}

		now := time.Now()
		entry.QueueStatus = "completed"
		entry.CompletedAt = &now
		if err := tx.Save(&entry).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Task{}).Where("id = ?", entry.TaskID).
			Update("status", "completed").Error; err != nil {
			return err
		}

		result = &entry
		return nil
	})

	return result, err
}
