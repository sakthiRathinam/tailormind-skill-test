#!/bin/bash
set -e

# This script initializes the PostgreSQL database with tables and seed data
# The database 'school_mgmt' is already created by docker-compose

echo "Starting database initialization for school_mgmt..."

# Run table creation first
echo "Creating tables..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/tables.sql

echo "Tables created successfully"

# Run seed data insertion
echo "Inserting seed data..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-initdb.d/seed-db.sql

echo "Seed data inserted successfully"
echo "Database initialization completed!" 