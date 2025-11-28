-- add role column to users table
ALTER TABLE users
ADD COLUMN role VARCHAR(10) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user'));
