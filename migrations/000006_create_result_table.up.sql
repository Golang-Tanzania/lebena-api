CREATE TABLE result (
	id_match UUID PRIMARY KEY REFERENCES match(id_match),
	home SMALLINT NOT NULL CHECK(home >= 0),
	away SMALLINT NOT NULL CHECK(away >= 0)
);
