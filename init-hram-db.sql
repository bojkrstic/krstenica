-- Create the 'hrams' database
CREATE DATABASE hrams;

-- Connect to the 'hrams' database
\c hrams;

-- Create the 'hram' table
CREATE TABLE "hram" (
    hram_id SERIAL PRIMARY KEY,
    naziv_hrama VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Display the created tables
\dt