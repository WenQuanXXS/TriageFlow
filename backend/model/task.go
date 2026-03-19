package model

import "time"

type Task struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	PatientName    string    `json:"patient_name" gorm:"not null"`
	ChiefComplaint string    `json:"chief_complaint" gorm:"not null"`
	Status         string    `json:"status" gorm:"default:pending"`
	Priority       string    `json:"priority" gorm:"default:normal"`
	Department     string    `json:"department"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
