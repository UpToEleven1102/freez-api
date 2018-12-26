package db

//Merchant
const schemaMerchant = `
	CREATE TABLE merchant(
		id VARCHAR(64) NOT NULL,
		online BOOLEAN NOT NULL DEFAULT 0,
		mobile BOOLEAN NOT NULL DEFAULT 0,
		phone_number text NOT NULL,
		email VARCHAR(64) NOT NULL UNIQUE,
		name text NOT NULL,
		password text NOT NULL,
		last_location POINT,
		image VARCHAR(255) NOT NULL DEFAULT 'https://freeze-app.s3.us-west-2.amazonaws.com/blank-profile-picture.jpg',
		PRIMARY KEY (id)
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
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(64) NOT NULL,
	price DECIMAL(10,2) NOT NULL,
	merchant_id VARCHAR(64) NOT NULL,
	image VARCHAR(255) NOT NULL DEFAULT '',
	PRIMARY KEY (id),
	FOREIGN KEY fk_product_merchant_id(merchant_id)
		REFERENCES merchant(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
)`

const schemaMerchantMOption = `
	CREATE TABLE merchant_m_option(
		id INT NOT NULL AUTO_INCREMENT,
		merchant_id VARCHAR(64) NOT NULL,
		add_convenience_fee BOOL NOT NULL DEFAULT false,
		PRIMARY KEY (id),
		FOREIGN KEY fk_option_merchant_id(merchant_id)
			REFERENCES merchant(id)
			ON UPDATE CASCADE
			ON DELETE CASCADE
	);
`

const triggerInsertMerchantMOption = `
	CREATE TRIGGER merchant_m_option AFTER INSERT ON merchant
	FOR EACH ROW INSERT INTO merchant_m_option (merchant_id) VALUES(NEW.id);
`


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
		email VARCHAR(64) NOT NULL UNIQUE,
		name text NOT NULL,
		password text NOT NULL,
		image VARCHAR(255) NOT NULL DEFAULT 'https://freeze-app.s3.us-west-2.amazonaws.com/blank-profile-picture.jpg',
		last_location POINT,
		PRIMARY KEY (id)
	);
`

const schemaMOption = `
	CREATE TABLE m_option(
		id INT AUTO_INCREMENT,
		user_id VARCHAR(64) NOT NULL,
		notif_fav_nearby BOOL NOT NULL DEFAULT TRUE,
		PRIMARY KEY (id),
		FOREIGN KEY fk_user (user_id)
			REFERENCES user(id)
			ON UPDATE CASCADE
			ON DELETE CASCADE
);
`

const triggerInsertUserMOption = `
	CREATE TRIGGER user_m_option AFTER INSERT ON user
	FOR EACH ROW INSERT INTO m_option (user_id) VALUES(NEW.id);
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
