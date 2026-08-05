#!/bin/sh
set -e

create_database_if_not_exists() {
    db="$1"
    exists=$(psql -U "$POSTGRES_USER" -tAc "SELECT 1 FROM pg_database WHERE datname = '$db'")

    if [ "$exists" != "1" ]; then
        echo "Creating database: $db"
        psql -U "$POSTGRES_USER" -c "CREATE DATABASE $db;"
    else
        echo "Database '$db' already exists, skipping creation."
    fi
}

create_database_if_not_exists "user_db"
create_database_if_not_exists "catalog_db"
create_database_if_not_exists "cart_db"
create_database_if_not_exists "order_db"
create_database_if_not_exists "payment_db"
create_database_if_not_exists "notification_db"