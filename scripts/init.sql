CREATE TABLE IF NOT EXISTS series (
  id UUID PRIMARY KEY NOT NULL,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  episodes INTEGER NOT NULL,
  begin_year SMALLINT NOT NULL,
  end_year SMALLINT DEFAULT 0 NOT NULL,
  creator VARCHAR(150)
);

CREATE OR REPLACE FUNCTION make_tsvector(title TEXT, description TEXT)
  RETURNS tsvector AS $$
BEGIN
  RETURN (setweight(to_tsvector('english', title), 'A') ||
    setweight(to_tsvector('english', description), 'B'));
END
$$ LANGUAGE 'plpgsql' IMMUTABLE;

CREATE INDEX IF NOT EXISTS idx_fts_series ON series
  USING gin(make_tsvector(title, description));

CREATE TABLE IF NOT EXISTS  reviews (
  id UUID NOT NULL,
  series_id UUID NOT NULL,
  author_id UUID NOT NULL,
  text TEXT NOT NULL,
  PRIMARY KEY(id),
  UNIQUE(author_id, series_id)
);
