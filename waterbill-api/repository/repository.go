package repository

import (
	"errors"
	"waterbill-api/models"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("record not found")

type Repository interface {
	CreateCustomer(c *models.Customer) error
	GetCustomers(name string, page, limit int) ([]models.Customer, int64, error)
	GetCustomerByID(id uint) (*models.Customer, error)

	CreateMeter(m *models.Meter) error
	GetMeterByID(id uint) (*models.Meter, error)

	GetLatestReadingByMeterID(meterID uint) (*models.Reading, error)
	CreateReading(r *models.Reading) error
	GetReadingByID(id uint) (*models.Reading, error)

	CreateBill(b *models.Bill) error
	GetBillByReadingID(readingID uint) (*models.Bill, error)
	GetBillByID(id uint) (*models.Bill, error)
	UpdateBill(b *models.Bill) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateCustomer(c *models.Customer) error {
	return r.db.Create(c).Error
}

func (r *GormRepository) GetCustomers(name string, page, limit int) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	query := r.db.Model(&models.Customer{})
	if name != "" {
		query = query.Where("full_name ILIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Order("id ASC").Find(&customers).Error
	return customers, total, err
}

func (r *GormRepository) GetCustomerByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Preload("Meters").First(&customer, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &customer, err
}

func (r *GormRepository) CreateMeter(m *models.Meter) error {
	return r.db.Create(m).Error
}

func (r *GormRepository) GetMeterByID(id uint) (*models.Meter, error) {
	var meter models.Meter
	err := r.db.First(&meter, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &meter, err
}

func (r *GormRepository) GetLatestReadingByMeterID(meterID uint) (*models.Reading, error) {
	var reading models.Reading
	err := r.db.Where("meter_id = ?", meterID).Order("id DESC").First(&reading).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &reading, err
}

func (r *GormRepository) CreateReading(rd *models.Reading) error {
	return r.db.Create(rd).Error
}

func (r *GormRepository) GetReadingByID(id uint) (*models.Reading, error) {
	var reading models.Reading
	err := r.db.First(&reading, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &reading, err
}

func (r *GormRepository) CreateBill(b *models.Bill) error {
	return r.db.Create(b).Error
}

func (r *GormRepository) GetBillByReadingID(readingID uint) (*models.Bill, error) {
	var bill models.Bill
	err := r.db.Where("reading_id = ?", readingID).First(&bill).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &bill, err
}

func (r *GormRepository) GetBillByID(id uint) (*models.Bill, error) {
	var bill models.Bill
	err := r.db.First(&bill, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &bill, err
}

func (r *GormRepository) UpdateBill(b *models.Bill) error {
	return r.db.Save(b).Error
}
