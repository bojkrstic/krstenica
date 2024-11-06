-- Create the 'eparhije' database
CREATE DATABASE eparhije;

-- Connect to the 'eparhije' database
\c eparhije;

-- Create the 'eparhija' table
CREATE TABLE "eparhija" (
    eparhija_id SERIAL PRIMARY KEY,
    naziv_eparhije VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Display the created tables
\dt