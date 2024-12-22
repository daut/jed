-- This file contains the SQL schema for the database.

CREATE TABLE IF NOT EXISTS products (
  id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  price DECIMAL(10, 2),
  inventory_count INTEGER DEFAULT 1 NOT NULL,

  CONSTRAINT positive_price CHECK (price > 0),
  CONSTRAINT positive_inventory_count CHECK (inventory_count > 0)
);

CREATE TABLE IF NOT EXISTS admins (
  id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens (
  hash bytea PRIMARY KEY,
  admin_id INTEGER NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
  expires_at TIMESTAMP(0) WITH TIME ZONE NOT NULL
);
