-- 1. ตาราง customers
CREATE TABLE public.customers (
    customer_id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    address VARCHAR(255)
);

-- 2. ตาราง meters (ยังไม่ใส่ Foreign Key เพราะข้อมูลเดิมมี customer_id 9-11 ที่ถูกลบไปแล้ว)
CREATE TABLE public.meters (
    meter_id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    meter_serial VARCHAR(50) UNIQUE NOT NULL,
    location VARCHAR(100)
);

-- 3. ตาราง readings (ยังไม่ใส่ Foreign Key เพราะข้อมูลเดิมมี meter_id 9-13 ที่ถูกลบไปแล้ว)
CREATE TABLE public.readings (
    reading_id SERIAL PRIMARY KEY,
    meter_id INT NOT NULL,
    billing_month VARCHAR(7) NOT NULL,
    prev_reading INT NOT NULL,
    curr_reading INT NOT NULL,
    units_used INT NOT NULL
);

-- 4. ตาราง bills 
CREATE TABLE public.bills (
    bill_id SERIAL PRIMARY KEY,
    reading_id INT UNIQUE NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(10) NOT NULL,
    due_date TIMESTAMPTZ NOT NULL,
    paid_date TIMESTAMPTZ,
    FOREIGN KEY (reading_id) REFERENCES public.readings(reading_id)
);