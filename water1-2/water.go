package main

import (
	"errors"
	"fmt"
)

func WaterBill(units int) (int, error) {
	if units < 0 {
		return 0, errors.New("(คืน error")
	}

	total := 0

	if units <= 10 {
		total = units * 5
	} else if units <= 20 {
		total = (10 * 5) + ((units - 10) * 7)
	} else {
		total = (10 * 5) + (10 * 7) + ((units - 20) * 9)
	}

	return total, nil
}

func main() {
	var units int
	fmt.Print("กรอกจำนวนน้ำ: ")
	_, err := fmt.Scan(&units)
	if err != nil {
		fmt.Println("Error: กรุณากรอกเป็นตัวเลขจำนวนเต็ม")
		return
	}

	price, err := WaterBill(units)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("ใช้น้ำ %d หน่วย = %d บาท\n", units, price)
}
