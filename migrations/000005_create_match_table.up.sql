CREATE TABLE match (
	id_match UUID PRIMARY KEY,
	id_home UUID NOT NULL REFERENCES team(id_team),
	id_away UUID  NOT NULL REFERENCES team(id_team),
	id_stadium UUID  NOT NULL REFERENCES stadium(id_stadium),
	date_time DATE NOT NULL
);
