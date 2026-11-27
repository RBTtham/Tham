package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	_ "waterbill-api/docs"
	"waterbill-api/handler"
	"waterbill-api/repository"
	"waterbill-api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Water Bill API
// @version 1.0
// @description API ระบบบริหารจัดการผู้ใช้น้ำ มิเตอร์ และการออกบิล
// @host localhost:3000
// @BasePath /
func main() {
	_ = godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// 1. รัน Migration จากโฟลเดอร์ migrations/*.up.sql
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatal("โหลดไฟล์ Migration ไม่สำเร็จ:", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("รัน Migration ไม่สำเร็จ:", err)
	}

	// 2. เชื่อมต่อ GORM
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("เชื่อมต่อฐานข้อมูลไม่สำเร็จ:", err)
	}

	repo := repository.NewGormRepository(db)
	svc := service.NewWaterBillService(repo)
	h := handler.NewWaterBillHandler(svc)

	app := fiber.New()

	// หน้า Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// เส้นทาง API ทั้ง 7 เส้นทาง
	app.Post("/customers", h.CreateCustomer)
	app.Get("/customers", h.ListCustomers)
	app.Get("/customers/:id", h.GetCustomer)
	app.Post("/meters", h.AddMeter)
	app.Post("/readings", h.RecordReading)
	app.Post("/bills/generate", h.GenerateBill)
	app.Post("/bills/:id/pay", h.PayBill)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3001"
	}

	log.Fatal(app.Listen(":" + port))
}
