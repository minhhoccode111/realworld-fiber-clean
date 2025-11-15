ALTER TABLE article_tags
  DROP CONSTRAINT article_tags_tag_id_fkey;

ALTER TABLE article_tags
  ADD CONSTRAINT article_tags_tag_id_fkey
    FOREIGN KEY (tag_id)
    REFERENCES tags(id)
    ON DELETE NO ACTION;
