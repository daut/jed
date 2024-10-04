CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255),
  description TEXT,
  price DECIMAL(10, 2)
);

CREATE TABLE admins (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255),
  password VARCHAR(255)
);
