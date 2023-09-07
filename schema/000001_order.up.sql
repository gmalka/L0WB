CREATE TABLE IF NOT EXISTS Orders (
    order_uid VARCHAR(250) PRIMARY KEY,
    track_number VARCHAR(250),
    entry VARCHAR(100),
    locale VARCHAR(100),
    internal_signature VARCHAR(100),
    customer_id VARCHAR(100),
    delivery_service VARCHAR(100),
    shardkey VARCHAR(100),
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS Deliveries (
    name VARCHAR(100) PRIMARY KEY,
    order_id VARCHAR(250) UNIQUE REFERENCES Orders(order_uid),
    phone VARCHAR(30),
    zip VARCHAR(100),
    city VARCHAR(100),
    address VARCHAR(100),
    region VARCHAR(100),
    email VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS Payments (
    transaction VARCHAR(100) PRIMARY KEY,
    order_id VARCHAR(250) UNIQUE REFERENCES Orders(order_uid),
    request_id VARCHAR(100),
    currency VARCHAR(100),
    provider VARCHAR(100),
    amount INTEGER,
    payment_dt INTEGER,
    bank VARCHAR(100),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);

CREATE TABLE IF NOT EXISTS Items (
    chrt_id INTEGER PRIMARY KEY,
    order_id VARCHAR(250) REFERENCES Orders(order_uid),
    track_number VARCHAR(100),
    price INTEGER,
    rid VARCHAR(100),
    name VARCHAR(100),
    sale INTEGER,
    size VARCHAR(100),
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR(100),
    status INTEGER
);