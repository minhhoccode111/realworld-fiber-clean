ALTER TABLE comments
ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();

-- update existing comments
UPDATE comments
SET updated_at = created_at;
