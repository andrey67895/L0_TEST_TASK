-- Таблицы
CREATE TABLE orders (
                        id VARCHAR PRIMARY KEY,
                        track_number VARCHAR,
                        entry VARCHAR,
                        locale VARCHAR,
                        internal_signature VARCHAR,
                        customer_id VARCHAR,
                        delivery_service VARCHAR,
                        shardkey VARCHAR,
                        sm_id BIGINT,
                        date_created TIMESTAMP,
                        oof_shard TEXT
);

CREATE TABLE delivery (
                          id BIGSERIAL PRIMARY KEY,
                          order_id VARCHAR REFERENCES orders(id) ON DELETE CASCADE,
                          name VARCHAR,
                          phone VARCHAR,
                          zip VARCHAR,
                          city VARCHAR,
                          address VARCHAR,
                          region VARCHAR,
                          email VARCHAR
);

CREATE TABLE payment (
                         id BIGSERIAL PRIMARY KEY,
                         order_id VARCHAR REFERENCES orders(id) ON DELETE CASCADE,
                         transaction VARCHAR,
                         request_id VARCHAR,
                         currency VARCHAR,
                         provider VARCHAR,
                         amount NUMERIC,
                         payment_dt BIGINT,
                         bank VARCHAR,
                         delivery_cost NUMERIC,
                         goods_total NUMERIC,
                         custom_fee NUMERIC
);

CREATE TABLE items (
                       id BIGSERIAL PRIMARY KEY,
                       order_id VARCHAR REFERENCES orders(id) ON DELETE CASCADE,
                       chrt_id BIGINT,
                       track_number VARCHAR,
                       price NUMERIC,
                       rid VARCHAR,
                       name VARCHAR,
                       sale NUMERIC,
                       size VARCHAR,
                       total_price NUMERIC,
                       nm_id BIGINT,
                       brand VARCHAR,
                       status INT
);

-- Дополнительные индексы
CREATE INDEX idx_delivery_order_id ON delivery(order_id);
CREATE INDEX idx_payment_order_id ON payment(order_id);
CREATE INDEX idx_items_order_id ON items(order_id);