package db

const schema = `
	CREATE TABLE merchant(
		id VARCHAR(64) NOT NULL,
		phone_number text NOT NULL,
		email VARCHAR(64) NOT NULL,
		name text NOT NULL,
		password text NOT NULL,
		image VARCHAR(50) NOT NULL DEFAULT './assets/profile.jpg',
		PRIMARY KEY (email)
	);
`
