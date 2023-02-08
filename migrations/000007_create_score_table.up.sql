CREATE TABLE score (
	id_score UUID PRIMARY KEY,
	id_match UUID  NOT NULL REFERENCES match(id_match),
	id_team UUID  NOT NULL REFERENCES team(id_team),
	id_player UUID  NOT NULL REFERENCES player(id_player),
	goals SMALLINT NOT NULL CHECK(goals >= 0),
	assists SMALLINT NOT NULL CHECK(assists >= 0)
);
