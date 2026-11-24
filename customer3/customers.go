package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// 1. อ่านไฟล์ customers.json
	data, err := os.ReadFile("customer.json")
	if err != nil {
		fmt.Println("อ่านไฟล์ไม่สำเร็จ:", err)
		return
	}

	// 2. แปลงข้อมูล JSON เข้า Slice
	var customers []Customer
	if err := json.Unmarshal(data, &customers); err != nil {
		fmt.Println("แปลงข้อมูล JSON ไม่สำเร็จ:", err)
		return
	}

	// 3. เรียงลำดับตามชื่อ (ก -> ฮ)
	sort.Slice(customers, func(i, j int) bool {
		return customers[i].Name < customers[j].Name
	})

	// 4. พิมพ์รายชื่อที่เรียงแล้วออกมา
	fmt.Println("=== รายชื่อลูกค้า (เรียงตามชื่อ) ===")
	for i, c := range customers {
		fmt.Printf("%d. %s (ID: %d | Email: %s)\n", i+1, c.Name, c.ID, c.Email)
	}
}
