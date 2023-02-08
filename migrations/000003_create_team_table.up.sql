CREATE TABLE team (
	id_team UUID PRIMARY KEY NOT NULL,
	team_photo VARCHAR(80) NOT NULL,
	id_stadium UUID REFERENCES stadium(id_stadium),
	club VARCHAR(25) NOT NULL
);
