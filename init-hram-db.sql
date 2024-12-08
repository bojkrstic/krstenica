-- Create the 'hrams' database
CREATE DATABASE hrams;

-- Connect to the 'hrams' database
\c hrams;

-- Create the 'hram' table
CREATE TABLE "hram" (
    hram_id SERIAL PRIMARY KEY,
    naziv_hrama VARCHAR(100) UNIQUE NOT NULL,
    status VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hram') THEN
        -- Dodaj kolonu 'status' ako ne postoji
        ALTER TABLE hram ADD COLUMN IF NOT EXISTS status VARCHAR(50);
    END IF;
END $$;


-- Display the created tables
\dt