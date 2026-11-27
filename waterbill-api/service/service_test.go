package service

import (
	"errors"
	"testing"
	"waterbill-api/models"
	"waterbill-api/repository"
)

// MockRepository จำลองฐานข้อมูลสำหรับรัน Test
type MockRepository struct {
	bills map[uint]*models.Bill
}

func (m *MockRepository) GetBillByID(id uint) (*models.Bill, error) {
	if b, ok := m.bills[id]; ok {
		return b, nil
	}
	return nil, repository.ErrNotFound
}
func (m *MockRepository) UpdateBill(b *models.Bill) error {
	m.bills[b.ID] = b
	return nil
}

// เมธอดอื่นที่ไม่ได้ใช้ในเคสนี้
func (m *MockRepository) CreateCustomer(c *models.Customer) error { return nil }
func (m *MockRepository) GetCustomers(name string, page, limit int) ([]models.Customer, int64, error) {
	return nil, 0, nil
}
func (m *MockRepository) GetCustomerByID(id uint) (*models.Customer, error) { return nil, nil }
func (m *MockRepository) CreateMeter(mt *models.Meter) error                { return nil }
func (m *MockRepository) GetMeterByID(id uint) (*models.Meter, error)       { return nil, nil }
func (m *MockRepository) GetLatestReadingByMeterID(meterID uint) (*models.Reading, error) {
	return nil, nil
}
func (m *MockRepository) CreateReading(r *models.Reading) error                   { return nil }
func (m *MockRepository) GetReadingByID(id uint) (*models.Reading, error)         { return nil, nil }
func (m *MockRepository) CreateBill(b *models.Bill) error                         { return nil }
func (m *MockRepository) GetBillByReadingID(readingID uint) (*models.Bill, error) { return nil, nil }

func TestCalculateWaterBill(t *testing.T) {
	tests := []struct {
		name      string
		units     int
		expected  int
		expectErr bool
	}{
		{"ใช้น้ำ 0 หน่วย", 0, 0, false},
		{"ใช้น้ำ 10 หน่วย", 10, 50, false},
		{"ใช้น้ำ 15 หน่วย", 15, 85, false},
		{"ใช้น้ำ 25 หน่วย", 25, 165, false},
		{"ติดลบต้อง error", -5, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateWaterBill(tt.units)
			if (err != nil) != tt.expectErr {
				t.Fatalf("error = %v, expectErr = %v", err, tt.expectErr)
			}
			if got != tt.expected {
				t.Errorf("ได้ %d; ต้องการ %d", got, tt.expected)
			}
		})
	}
}

func TestPayBill(t *testing.T) {
	mockRepo := &MockRepository{
		bills: map[uint]*models.Bill{
			1: {ID: 1, Amount: 165, Status: "UNPAID"},
			2: {ID: 2, Amount: 85, Status: "PAID"},
		},
	}
	svc := NewWaterBillService(mockRepo)

	t.Run("จ่ายบิลสำเร็จ สถานะต้องเปลี่ยนเป็น PAID", func(t *testing.T) {
		bill, err := svc.PayBill(1)
		if err != nil {
			t.Fatalf("ไม่ควรเกิด error: %v", err)
		}
		if bill.Status != "PAID" || bill.PaidDate == nil {
			t.Errorf("สถานะไม่เปลี่ยนเป็น PAID หรือไม่มีวันที่จ่าย")
		}
	})

	t.Run("จ่ายบิลที่จ่ายไปแล้ว ต้องคืน ErrAlreadyPaid (409)", func(t *testing.T) {
		_, err := svc.PayBill(2)
		if !errors.Is(err, ErrAlreadyPaid) {
			t.Errorf("ต้องการ ErrAlreadyPaid แต่ได้ %v", err)
		}
	})
}
