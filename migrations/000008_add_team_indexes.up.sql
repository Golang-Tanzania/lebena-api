CREATE INDEX IF NOT EXISTS team_club_idx ON team USING GIN (to_tsvector('simple', club));