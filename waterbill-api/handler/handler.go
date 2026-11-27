package handler

import (
	"errors"
	"waterbill-api/models"
	"waterbill-api/service"

	"github.com/gofiber/fiber/v2"
)

type WaterBillHandler struct {
	svc *service.WaterBillService
}

func NewWaterBillHandler(svc *service.WaterBillService) *WaterBillHandler {
	return &WaterBillHandler{svc: svc}
}

// CreateCustomer godoc
// @Summary เพิ่มผู้ใช้น้ำ
// @Tags customers
// @Accept json
// @Produce json
// @Param request body models.CreateCustomerRequest true "ข้อมูลผู้ใช้น้ำ"
// @Success 201 {object} models.Customer
// @Failure 422 {object} map[string]string
// @Router /customers [post]
func (h *WaterBillHandler) CreateCustomer(c *fiber.Ctx) error {
	var req models.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid json"})
	}

	customer, err := h.svc.CreateCustomer(req)
	if errors.Is(err, service.ErrValidation) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "ชื่อและเบอร์โทรห้ามว่าง"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(customer)
}

// ListCustomers godoc
// @Summary ดูรายชื่อผู้ใช้น้ำ (ค้นหาด้วยชื่อและแบ่งหน้าได้)
// @Tags customers
// @Produce json
// @Param name query string false "ค้นหาตามชื่อ"
// @Param page query int false "หน้าที่ (default 1)"
// @Param limit query int false "จำนวนต่อหน้า (default 10)"
// @Success 200 {object} map[string]interface{}
// @Router /customers [get]
func (h *WaterBillHandler) ListCustomers(c *fiber.Ctx) error {
	name := c.Query("name")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	customers, total, err := h.svc.ListCustomers(name, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  customers,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GetCustomer godoc
// @Summary ดูข้อมูลผู้ใช้น้ำ 1 คน
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} models.Customer
// @Failure 404 {object} map[string]string
// @Router /customers/{id} [get]
func (h *WaterBillHandler) GetCustomer(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	customer, err := h.svc.GetCustomer(uint(id))
	if errors.Is(err, service.ErrCustomerMissing) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ไม่พบผู้ใช้น้ำ"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(customer)
}

// AddMeter godoc
// @Summary เพิ่มมิเตอร์ให้ผู้ใช้น้ำ
// @Tags meters
// @Accept json
// @Produce json
// @Param request body models.CreateMeterRequest true "ข้อมูลมิเตอร์"
// @Success 201 {object} models.Meter
// @Failure 404 {object} map[string]string
// @Router /meters [post]
func (h *WaterBillHandler) AddMeter(c *fiber.Ctx) error {
	var req models.CreateMeterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}

	meter, err := h.svc.AddMeter(req)
	if errors.Is(err, service.ErrCustomerMissing) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ไม่มีลูกค้าคนนี้"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(meter)
}

// RecordReading godoc
// @Summary บันทึกเลขที่จดได้จากมิเตอร์
// @Tags readings
// @Accept json
// @Produce json
// @Param request body models.CreateReadingRequest true "ข้อมูลเลขมิเตอร์"
// @Success 201 {object} models.Reading
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /readings [post]
func (h *WaterBillHandler) RecordReading(c *fiber.Ctx) error {
	var req models.CreateReadingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}

	reading, err := h.svc.RecordReading(req)
	if errors.Is(err, service.ErrResourceMissing) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ไม่พบมิเตอร์"})
	}
	if errors.Is(err, service.ErrValidation) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "เลขมิเตอร์น้อยกว่าเดือนก่อน"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(reading)
}

// GenerateBill godoc
// @Summary สร้างบิลค่าน้ำจากเลขที่จด
// @Tags bills
// @Accept json
// @Produce json
// @Param request body models.GenerateBillRequest true "รหัสการจดมิเตอร์"
// @Success 201 {object} models.Bill
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /bills/generate [post]
func (h *WaterBillHandler) GenerateBill(c *fiber.Ctx) error {
	var req models.GenerateBillRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}

	bill, err := h.svc.GenerateBill(req.ReadingID)
	if errors.Is(err, service.ErrResourceMissing) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ไม่พบรายการจดมิเตอร์"})
	}
	if errors.Is(err, service.ErrBillExists) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "ออกบิลสำหรับรายการนี้ไปแล้ว"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(bill)
}

// PayBill godoc
// @Summary จ่ายบิลค่าน้ำ
// @Tags bills
// @Produce json
// @Param id path int true "Bill ID"
// @Success 200 {object} models.Bill
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /bills/{id}/pay [post]
func (h *WaterBillHandler) PayBill(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	bill, err := h.svc.PayBill(uint(id))
	if errors.Is(err, service.ErrResourceMissing) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ไม่พบบิล"})
	}
	if errors.Is(err, service.ErrAlreadyPaid) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "บิลนี้จ่ายไปแล้ว"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(bill)
}
