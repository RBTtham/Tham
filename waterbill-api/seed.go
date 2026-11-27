package main

import (
	"time"
	"waterbill-api/models"

	"gorm.io/gorm"
)

func seedData(db *gorm.DB) error {
	var count int64
	db.Model(&models.Customer{}).Count(&count)
	if count > 0 {
		return nil // มีข้อมูลในฐานข้อมูลแล้ว ไม่เพิ่มซ้ำ
	}

	// 1. เพิ่มผู้ใช้น้ำ 6 คน (ทดสอบการแบ่งหน้า page/limit และค้นหาชื่อได้)
	customers := []models.Customer{
		{FullName: "สมชาย ใจดี", Phone: "081-111-1111", Address: "12/3 หมู่ 1 ต.ในเมือง"},
		{FullName: "วิภาดา รักสงบ", Phone: "082-222-2222", Address: "45/8 หมู่ 2 ต.ในเมือง"},
		{FullName: "อนันต์ มั่นคง", Phone: "083-333-3333", Address: "99/1 หมู่ 4 ต.รอบเวียง"},
		{FullName: "กานดา สดใส", Phone: "084-444-4444", Address: "18/5 ซอยสุขใจ ต.ในเมือง"},
		{FullName: "สมศักดิ์ พัฒนา", Phone: "085-555-5555", Address: "204 หมู่ 7 ต.หนองบัว"},
		{FullName: "นารี ศรีสุข", Phone: "086-666-6666", Address: "77/2 ถ.มิตรภาพ ต.ในเมือง"},
	}
	if err := db.Create(&customers).Error; err != nil {
		return err
	}

	// 2. เพิ่มมิเตอร์น้ำ 8 ตัว (บางคนมี 2 ตัว)
	meters := []models.Meter{
		{CustomerID: customers[0].ID, MeterSerial: "MTR-2026-001", Location: "บ้านพักอาศัย (สมชาย)"},
		{CustomerID: customers[0].ID, MeterSerial: "MTR-2026-002", Location: "ร้านกาแฟหน้าบ้าน (สมชาย)"},
		{CustomerID: customers[1].ID, MeterSerial: "MTR-2026-003", Location: "บ้านพักอาศัย (วิภาดา)"},
		{CustomerID: customers[2].ID, MeterSerial: "MTR-2026-004", Location: "อาคารพาณิชย์ (อนันต์)"},
		{CustomerID: customers[3].ID, MeterSerial: "MTR-2026-005", Location: "บ้านเดี่ยว (กานดา)"},
		{CustomerID: customers[4].ID, MeterSerial: "MTR-2026-006", Location: "บ้านพักอาศัย (สมศักดิ์)"},
		{CustomerID: customers[4].ID, MeterSerial: "MTR-2026-007", Location: "สวนเกษตรหลังบ้าน (สมศักดิ์)"},
		{CustomerID: customers[5].ID, MeterSerial: "MTR-2026-008", Location: "ร้านอาหารตามสั่ง (นารี)"},
	}
	if err := db.Create(&meters).Error; err != nil {
		return err
	}

	// 3. เพิ่มรายการจดมิเตอร์ 8 รายการ (ครอบคลุมทั้งเรท 0-10, 11-20, 21+ หน่วย)
	readings := []models.Reading{
		{MeterID: meters[0].ID, BillingMonth: "2026-10", PrevReading: 100, CurrReading: 115, UnitsUsed: 15}, // ID 1 -> 85 บ.
		{MeterID: meters[1].ID, BillingMonth: "2026-10", PrevReading: 250, CurrReading: 275, UnitsUsed: 25}, // ID 2 -> 165 บ.
		{MeterID: meters[2].ID, BillingMonth: "2026-10", PrevReading: 50, CurrReading: 58, UnitsUsed: 8},    // ID 3 -> 40 บ.
		{MeterID: meters[3].ID, BillingMonth: "2026-09", PrevReading: 400, CurrReading: 420, UnitsUsed: 20}, // ID 4 -> 120 บ.
		{MeterID: meters[3].ID, BillingMonth: "2026-10", PrevReading: 420, CurrReading: 450, UnitsUsed: 30}, // ID 5 -> 210 บ.
		{MeterID: meters[4].ID, BillingMonth: "2026-10", PrevReading: 10, CurrReading: 22, UnitsUsed: 12},   // ID 6 -> 64 บ.
		{MeterID: meters[5].ID, BillingMonth: "2026-10", PrevReading: 300, CurrReading: 305, UnitsUsed: 5},  // ID 7 -> ยังไม่ออกบิล (เอาไว้ลอง POST /bills/generate)
		{MeterID: meters[7].ID, BillingMonth: "2026-10", PrevReading: 500, CurrReading: 540, UnitsUsed: 40}, // ID 8 -> ยังไม่ออกบิล (เอาไว้ลอง POST /bills/generate)
	}
	if err := db.Create(&readings).Error; err != nil {
		return err
	}

	// 4. เพิ่มบิลค่าน้ำ 6 ใบ (จ่ายแล้ว 3 ใบ, ค้างจ่าย 3 ใบ)
	now := time.Now()
	bills := []models.Bill{
		{ReadingID: readings[0].ID, Amount: 85, Status: "PAID", DueDate: now.AddDate(0, 0, -5), PaidDate: &now},
		{ReadingID: readings[1].ID, Amount: 165, Status: "PAID", DueDate: now.AddDate(0, 0, -5), PaidDate: &now},
		{ReadingID: readings[2].ID, Amount: 40, Status: "UNPAID", DueDate: now.AddDate(0, 0, 10), PaidDate: nil},
		{ReadingID: readings[3].ID, Amount: 120, Status: "PAID", DueDate: now.AddDate(0, -1, 0), PaidDate: &now},
		{ReadingID: readings[4].ID, Amount: 210, Status: "UNPAID", DueDate: now.AddDate(0, 0, 10), PaidDate: nil},
		{ReadingID: readings[5].ID, Amount: 64, Status: "UNPAID", DueDate: now.AddDate(0, 0, 10), PaidDate: nil},
	}
	return db.Create(&bills).Error
}
