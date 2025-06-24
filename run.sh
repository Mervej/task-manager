#!/bin/bash

# Create data directory if it doesn't exist
mkdir -p data

# Check if SQLite3 command-line tool is available
if command -v sqlite3 &> /dev/null; then
    echo "Creating database schema..."
    # Create the database schema using the migration file
    sqlite3 data/tasks.db < migrations/001_create_tasks_table.up.sql
    
    if [ $? -eq 0 ]; then
        echo "Database schema created successfully"
    else
        echo "Error creating database schema"
        exit 1
    fi
else
    echo "SQLite3 command-line tool not found. The database will be initialized at runtime."
fi

# Run the server
echo "Starting the server..."
go run cmd/server/main.go
