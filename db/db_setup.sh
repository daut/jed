DB_NAME="simpshop"

# create a setup function but print all errors and wrap everything in a transaction
setup_database() {
  echo "Setting up database..."
  psql -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;" || (echo "Failed to drop database" && exit 1)
  psql -d postgres -c "CREATE DATABASE $DB_NAME;" || (echo "Failed to create database" && exit 1)
  psql -d $DB_NAME -f ./schema.sql || (echo "Failed to create schema" && exit 1)
  psql -d $DB_NAME -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;" || (echo "Failed to create pgcrypto extension" && exit 1)
  psql -d $DB_NAME -f ./seeds/admins.sql || (echo "Failed to seed admins" && exit 1)
  psql -d $DB_NAME -f ./seeds/products.sql || (echo "Failed to seed products" && exit 1)
  echo "Database setup complete."
}

# cd into the directory of this script
echo $(dirname "$0")
cd "$(dirname "$0")"
setup_database
