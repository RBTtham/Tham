CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    address TEXT
);

CREATE TABLE IF NOT EXISTS meters (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    meter_serial TEXT NOT NULL UNIQUE,
    location TEXT
);

CREATE TABLE IF NOT EXISTS readings (
    id SERIAL PRIMARY KEY,
    meter_id INTEGER NOT NULL REFERENCES meters(id) ON DELETE CASCADE,
    billing_month TEXT NOT NULL,
    prev_reading INTEGER NOT NULL,
    curr_reading INTEGER NOT NULL,
    units_used INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS bills (
    id SERIAL PRIMARY KEY,
    reading_id INTEGER NOT NULL UNIQUE REFERENCES readings(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'UNPAID',
    due_date TIMESTAMPTZ,
    paid_date TIMESTAMPTZ
);