package models

import "time"

type Customer struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	FullName string  `gorm:"not null" json:"full_name"`
	Phone    string  `gorm:"not null" json:"phone"`
	Address  string  `json:"address"`
	Meters   []Meter `gorm:"foreignKey:CustomerID" json:"meters,omitempty"`
}

type Meter struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	CustomerID  uint   `gorm:"not null" json:"customer_id"`
	MeterSerial string `gorm:"uniqueIndex;not null" json:"meter_serial"`
	Location    string `json:"location"`
}

type Reading struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	MeterID      uint   `gorm:"not null" json:"meter_id"`
	BillingMonth string `gorm:"not null" json:"billing_month"`
	PrevReading  int    `gorm:"not null" json:"prev_reading"`
	CurrReading  int    `gorm:"not null" json:"curr_reading"`
	UnitsUsed    int    `gorm:"not null" json:"units_used"`
}

type Bill struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ReadingID uint       `gorm:"uniqueIndex;not null" json:"reading_id"`
	Amount    int        `gorm:"not null" json:"amount"`
	Status    string     `gorm:"not null;default:'UNPAID'" json:"status"` // UNPAID หรือ PAID
	DueDate   time.Time  `json:"due_date"`
	PaidDate  *time.Time `json:"paid_date"`
}

// --- DTO สำหรับรับค่า Request ---

type CreateCustomerRequest struct {
	FullName string `json:"full_name" example:"สมชาย ใจดี"`
	Phone    string `json:"phone" example:"081-1111111"`
	Address  string `json:"address" example:"12/3 ต.ในเมือง"`
}

type CreateMeterRequest struct {
	CustomerID  uint   `json:"customer_id" example:"1"`
	MeterSerial string `json:"meter_serial" example:"MTR-2026-001"`
	Location    string `json:"location" example:"บ้านพักอาศัย"`
}

type CreateReadingRequest struct {
	MeterID      uint   `json:"meter_id" example:"1"`
	BillingMonth string `json:"billing_month" example:"2026-11"`
	CurrReading  int    `json:"curr_reading" example:"115"`
}

type GenerateBillRequest struct {
	ReadingID uint `json:"reading_id" example:"1"`
}
