-- update existing users with null bio value
UPDATE users SET bio = '' WHERE bio IS NULL;

ALTER TABLE users
ALTER COLUMN bio SET DEFAULT '',
ALTER COLUMN bio SET NOT NULL;
