CREATE TABLE products (
 id uuid PRIMARY KEY, name text NOT NULL, spec text NOT NULL, category text NOT NULL,
 price bigint NOT NULL CHECK(price>0), old bigint NOT NULL DEFAULT 0 CHECK(old>=0),
 stock integer NOT NULL CHECK(stock>=0), type text NOT NULL, color text NOT NULL, badge text NOT NULL DEFAULT ''
);
CREATE INDEX products_category_idx ON products(category);
CREATE TABLE orders (id uuid PRIMARY KEY,body jsonb NOT NULL,total bigint NOT NULL CHECK(total>0),status text NOT NULL CHECK(status IN ('Новый','Подтверждён','Собирается','Готов к выдаче','Завершён')),created_at timestamptz NOT NULL DEFAULT now());
CREATE INDEX orders_created_idx ON orders(created_at DESC);
