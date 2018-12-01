package db

const schemaMerchant = `
	CREATE TABLE merchant(
		id VARCHAR(64) NOT NULL,
		online BOOLEAN NOT NULL DEFAULT 0,
		phone_number text NOT NULL,
		email VARCHAR(64) NOT NULL,
		name text NOT NULL,
		password text NOT NULL,
		image VARCHAR(50) NOT NULL DEFAULT './assets/profile.jpg',
		PRIMARY KEY (email)
	);
`

const schemaLocation = `
	CREATE TABLE location(
		id INT NOT NULL AUTO_INCREMENT,
		ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		merchant_id VARCHAR(64) NOT NULL,
		location POINT NOT NULL,
		SPATIAL INDEX(location),
		PRIMARY KEY(id)
	);
`

const schemaUser = `
	CREATE TABLE user(
		id VARCHAR(64) NOT NULL,
		phone_number text NOT NULL,
		email VARCHAR(64) NOT NULL,
		name text NOT NULL,
		password text NOT NULL,
		image VARCHAR(50) NOT NULL DEFAULT './assets/profile.jpg',
		PRIMARY KEY (email)
	);
`

const schemaRequest = `CREATE TABLE request(
		id INT NOT NULL AUTO_INCREMENT,
		user_id VARCHAR(64) NOT NULL,
		merchant_id VARCHAR(64) NOT NULL,
		location POINT NOT NULL,
		SPATIAL INDEX(location),
		PRIMARY KEY (id)
	);
`
