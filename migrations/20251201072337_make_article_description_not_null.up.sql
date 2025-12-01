UPDATE articles
SET description = 'description cannot be empty'
WHERE description IS NULL;

ALTER TABLE articles
ALTER COLUMN description SET NOT NULL;
