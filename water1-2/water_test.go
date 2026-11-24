package main

import "testing"

func TestWaterBill(t *testing.T) {
	tests := []struct {
		name      string
		units     int
		expected  int
		expectErr bool
	}{
		{
			name:      "ใช้น้ำ 0 หน่วย",
			units:     0,
			expected:  0,
			expectErr: false,
		},
		{
			name:      "ใช้น้ำ 10 หน่วย",
			units:     10,
			expected:  50,
			expectErr: false,
		},
		{
			name:      "ใช้น้ำ 15 หน่วย",
			units:     15,
			expected:  85,
			expectErr: false,
		},
		{
			name:      "ใช้น้ำ 20 หน่วย",
			units:     20,
			expected:  120,
			expectErr: false,
		},
		{
			name:      "ใช้น้ำ 25 หน่วย",
			units:     25,
			expected:  165,
			expectErr: false,
		},
		{
			name:      "ใช้น้ำ 50 หน่วย",
			units:     50,
			expected:  390,
			expectErr: false,
		},
		{
			name:      "ค่าติดลบ (-1 หน่วย) ต้องคืน error",
			units:     -1,
			expected:  0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WaterBill(tt.units)

			if (err != nil) != tt.expectErr {
				t.Errorf("WaterBill(%d) error = %v, expectErr = %v", tt.units, err, tt.expectErr)
				return
			}

			if got != tt.expected {
				t.Errorf("WaterBill(%d) = %d บาท; ต้องการ %d บาท", tt.units, got, tt.expected)
			}
		})
	}
}
