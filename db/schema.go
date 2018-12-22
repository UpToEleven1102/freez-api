package db

//Merchant
const schemaMerchant = `
	CREATE TABLE merchant(
		id VARCHAR(64) NOT NULL,
		online BOOLEAN NOT NULL DEFAULT 0,
		mobile BOOLEAN NOT NULL DEFAULT 0,
		phone_number text NOT NULL,
		email VARCHAR(64) NOT NULL,
		name text NOT NULL,
		password text NOT NULL,
		image VARCHAR(255) NOT NULL DEFAULT './assets/profile.jpg',
		PRIMARY KEY (email)
	);
`

const schemaLocation = `
	CREATE TABLE location(
		id INT NOT NULL AUTO_INCREMENT,
		ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP UNIQUE,
		merchant_id VARCHAR(64) NOT NULL,
		location POINT NOT NULL,
		SPATIAL INDEX(location),
		PRIMARY KEY(id)
	);
`

const schemaProduct = `CREATE TABLE product(
	
)`

const schemaMerchantNotification = `CREATE TABLE merchant_notification(
	id INT NOT NULL AUTO_INCREMENT,
	ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP UNIQUE,
	merchant_id VARCHAR(64) NOT NULL,
	activity_type INT
	source_id VARCHAR(64) NOT NULL
);`

//User
const schemaUser = `
	CREATE TABLE user(
		id VARCHAR(64) NOT NULL,
		phone_number text NOT NULL,
		email VARCHAR(64) NOT NULL,
		name text NOT NULL,
		password text NOT NULL,
		image VARCHAR(255) NOT NULL DEFAULT './assets/profile.jpg',
		last_location POINT,
		PRIMARY KEY (email)
	);
`

const schemaFavorites = `
	CREATE TABLE favorite(
		user_id VARCHAR(64) NOT NULL,
		merchant_id VARCHAR(64) NOT NULL,
		PRIMARY KEY(user_id, merchant_id)
	);
`

const schemaUserNotification = `CREATE TABLE user_notification(
	id INT NOT NULL AUTO_INCREMENT,
	ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP UNIQUE,
	user_id VARCHAR(64) NOT NULL,
	activity_type INT
	source_id VARCHAR(64) NOT NULL
);`


const schemaRequest = `CREATE TABLE request(
		id INT NOT NULL AUTO_INCREMENT,
		user_id VARCHAR(64) NOT NULL,
		merchant_id VARCHAR(64) NOT NULL,
		location POINT NOT NULL,
		SPATIAL INDEX(location),
		comment VARCHAR(200) NOT NULL DEFAULT '',
		accepted TINYINT(1) DEFAULT -1 NOT NULL,
		PRIMARY KEY (id)
	);
`

