package service

import (
	"errors"
	"strings"
	"time"
	"waterbill-api/models"
	"waterbill-api/repository"
)

var (
	ErrValidation      = errors.New("validation error")
	ErrCustomerMissing = errors.New("customer not found")
	ErrResourceMissing = errors.New("resource not found")
	ErrAlreadyPaid     = errors.New("bill is already paid")
	ErrBillExists      = errors.New("bill already generated for this reading")
)

type WaterBillService struct {
	repo repository.Repository
}

func NewWaterBillService(repo repository.Repository) *WaterBillService {
	return &WaterBillService{repo: repo}
}

// CalculateWaterBill ฟังก์ชันคิดค่าน้ำจากบทที่ 3 (1-10: 5บ., 11-20: 7บ., 21+: 9บ.)
func CalculateWaterBill(units int) (int, error) {
	if units < 0 {
		return 0, errors.New("จำนวนหน่วยน้ำต้องไม่ติดลบ")
	}

	if units <= 10 {
		return units * 5, nil
	} else if units <= 20 {
		return (10 * 5) + ((units - 10) * 7), nil
	}
	return (10 * 5) + (10 * 7) + ((units - 20) * 9), nil
}
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func (s *WaterBillService) CreateCustomer(req models.CreateCustomerRequest) (*models.Customer, error) {
	name := strings.TrimSpace(req.FullName)
	phone := strings.ReplaceAll(strings.TrimSpace(req.Phone), "-", "")

	if name == "" || phone == "" || !isNumeric(phone) {
		return nil, ErrValidation // ตอบ 422 ทันทีถ้าชื่อว่างหรือเบอร์ไม่ใช่ตัวเลขล้วน
	}

	customer := &models.Customer{
		FullName: strings.TrimSpace(req.FullName),
		Phone:    strings.TrimSpace(req.Phone),
		Address:  strings.TrimSpace(req.Address),
	}
	if err := s.repo.CreateCustomer(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *WaterBillService) ListCustomers(name string, page, limit int) ([]models.Customer, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.GetCustomers(name, page, limit)
}

func (s *WaterBillService) GetCustomer(id uint) (*models.Customer, error) {
	customer, err := s.repo.GetCustomerByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrCustomerMissing
	}
	return customer, err
}

func (s *WaterBillService) AddMeter(req models.CreateMeterRequest) (*models.Meter, error) {
	if _, err := s.repo.GetCustomerByID(req.CustomerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCustomerMissing
		}
		return nil, err
	}

	meter := &models.Meter{
		CustomerID:  req.CustomerID,
		MeterSerial: req.MeterSerial,
		Location:    req.Location,
	}
	if err := s.repo.CreateMeter(meter); err != nil {
		return nil, err
	}
	return meter, nil
}

func (s *WaterBillService) RecordReading(req models.CreateReadingRequest) (*models.Reading, error) {
	if _, err := s.repo.GetMeterByID(req.MeterID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrResourceMissing
		}
		return nil, err
	}

	prevReading := 0
	latest, err := s.repo.GetLatestReadingByMeterID(req.MeterID)
	if err == nil {
		prevReading = latest.CurrReading
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// กรณีพิเศษ: เลขที่จดได้น้อยกว่าเดือนก่อน ตอบ 422
	if req.CurrReading < prevReading {
		return nil, ErrValidation
	}

	reading := &models.Reading{
		MeterID:      req.MeterID,
		BillingMonth: req.BillingMonth,
		PrevReading:  prevReading,
		CurrReading:  req.CurrReading,
		UnitsUsed:    req.CurrReading - prevReading,
	}
	if err := s.repo.CreateReading(reading); err != nil {
		return nil, err
	}
	return reading, nil
}

func (s *WaterBillService) GenerateBill(readingID uint) (*models.Bill, error) {
	reading, err := s.repo.GetReadingByID(readingID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrResourceMissing
		}
		return nil, err
	}

	if _, err := s.repo.GetBillByReadingID(readingID); err == nil {
		return nil, ErrBillExists
	}

	amount, err := CalculateWaterBill(reading.UnitsUsed)
	if err != nil {
		return nil, ErrValidation
	}

	bill := &models.Bill{
		ReadingID: reading.ID,
		Amount:    amount,
		Status:    "UNPAID",
		DueDate:   time.Now().AddDate(0, 0, 15),
	}
	if err := s.repo.CreateBill(bill); err != nil {
		return nil, err
	}
	return bill, nil
}

func (s *WaterBillService) PayBill(billID uint) (*models.Bill, error) {
	bill, err := s.repo.GetBillByID(billID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrResourceMissing
		}
		return nil, err
	}

	// กรณีพิเศษ: จ่ายไปแล้ว ตอบ 409
	if bill.Status == "PAID" {
		return nil, ErrAlreadyPaid
	}

	now := time.Now()
	bill.Status = "PAID"
	bill.PaidDate = &now

	if err := s.repo.UpdateBill(bill); err != nil {
		return nil, err
	}
	return bill, nil
}
