package model

import "time"

// QueueEntry represents a patient's position in a department queue.
// It is created automatically when a task completes triage.
type QueueEntry struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	TaskID      uint       `json:"task_id" gorm:"uniqueIndex;not null"`
	PatientName string     `json:"patient_name" gorm:"not null"`
	Department  string     `json:"department" gorm:"index;not null"`
	Priority    string     `json:"priority" gorm:"not null"`
	QueueStatus string     `json:"queue_status" gorm:"index;default:waiting"`
	QueueNumber int        `json:"queue_number" gorm:"not null"`
	QueueOrder  int        `json:"queue_order" gorm:"index;not null"`
	CalledAt    *time.Time `json:"called_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Denormalized from Task for display convenience
	ChiefComplaint string `json:"chief_complaint" gorm:"not null"`
}

// PriorityWeight returns a numeric weight for priority-based ordering.
// Lower weight = higher priority (sorted first).
func PriorityWeight(priority string) int {
	switch priority {
	case "urgent":
		return 0
	case "high":
		return 1
	case "normal":
		return 2
	default:
		return 3
	}
}
