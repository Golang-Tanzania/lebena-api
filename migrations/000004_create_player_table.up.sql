CREATE TABLE player (
	id_player UUID PRIMARY KEY,
	id_team UUID NOT NULL REFERENCES team(id_team),
	first_name VARCHAR(15) NOT NULL,
	last_name VARCHAR(15) NOT NULL,
	kit SMALLINT NOT NULL CHECK(kit between 0 and 99),
	position VARCHAR(15) NOT NULL CHECK(
		position = 'Goalkeeper' OR position = 'Defender' OR position = 'Mid fielder' OR  position = 'Striker'),
	region VARCHAR(32),
    player_photo VARCHAR(80),
	previous_club VARCHAR(50)
);