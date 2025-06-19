CREATE TABLE
	privileges (
		id CHAR(36) NOT NULL PRIMARY KEY,
		name VARCHAR(64) NOT NULL UNIQUE,
		order_index INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
	);

INSERT INTO
	privileges (id, name, order_index)
VALUES
	(
		'706773b3-ea2c-40dd-a8e0-77482910c28b',
		'owner',
		100
	),
	(
		'4f52c9cd-514c-4fe0-ae70-5a5eb342d749',
		'admin',
		80
	),
	(
		'a95ca130-7687-4c0f-9c1e-aec1790be777',
		'moderator',
		60
	),
	(
		'b43583ec-22a4-4bbf-9f72-a137ef68e97d',
		'tester',
		40
	);