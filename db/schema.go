package db

//Merchant
const schemaMerchant = `
	CREATE TABLE merchant(
		id VARCHAR(64) NOT NULL,
		facebook_id varchar(64) UNIQUE,
		category VARCHAR(64) NOT NULL DEFAULT 'ice_cream_truck',
		online BOOLEAN NOT NULL DEFAULT 0,
		mobile BOOLEAN NOT NULL DEFAULT 0,
		phone_number text NOT NULL,
		email VARCHAR(64) NOT NULL,
		name text NOT NULL,
		description text NULL,
		password text NOT NULL,
		last_location POINT,
		image VARCHAR(255) NOT NULL DEFAULT 'https://freeze-app.s3.us-west-2.amazonaws.com/blank-profile-picture.jpg',
		stripe_id VARCHAR(64) NOT NULL UNIQUE,
		PRIMARY KEY (id),
		FOREIGN KEY merchant_category (category)
			REFERENCES merchant_category(category)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
`

const schemaMerchantCategory = `
	CREATE TABLE merchant_category(
		category VARCHAR(64) NOT NULL,
		PRIMARY KEY (category)
	);
`


const schemaLocation = `
	CREATE TABLE location(
		id INT NOT NULL AUTO_INCREMENT,
		ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP UNIQUE,
		merchant_id VARCHAR(64) NOT NULL,
		location POINT NOT NULL,
		SPATIAL INDEX(location),
		PRIMARY KEY(id),
		FOREIGN KEY fk_merchant_location (merchant_id)
			REFERENCES merchant(id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
`

const schemaProduct = `CREATE TABLE product(
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(64) NOT NULL,
	description TEXT NULL,
	price DECIMAL(10,2) NOT NULL,
	merchant_id VARCHAR(64) NOT NULL,
	image VARCHAR(255) NOT NULL DEFAULT 'https://www.houstonfoodbank.org/wp-content/uploads/2018/01/homepage_boxfood-276x300.png',
	PRIMARY KEY (id),
	FOREIGN KEY fk_product_merchant_id(merchant_id)
		REFERENCES merchant(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
)`

const schemaOrder = `CREATE TABLE m_order(
	id INT NOT NULL AUTO_INCREMENT,
	user_id VARCHAR(64) NOT NULL,
	merchant_id VARCHAR(64) NOT NULL,
	stripe_id VARCHAR(64) NOT NULL,
	refund BOOL NOT NULL DEFAULT FALSE,
	amount FLOAT NOT NULL,
	date DATETIME NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id)
	);
	`

const schemaOrderProduct = `CREATE TABLE m_order_product(
	id INT NOT NULL AUTO_INCREMENT,
	order_id INT NOT NULL,
	product_id INT NOT NULL,
	quantity INT NOT NULL,
	price FLOAT NOT NULL,
	PRIMARY KEY(id),
	FOREIGN KEY fk_order_id(order_id)
		REFERENCES m_order(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	FOREIGN KEY fk_product_id(product_id)
		REFERENCES product(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
	)
	`

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

const schemaActivityType = `
	CREATE TABLE activity_type(
	id INT NOT NULL AUTO_INCREMENT,
	type VARCHAR(64) NOT NULL,
	PRIMARY KEY (id))
`

const schemaMerchantNotification = `CREATE TABLE merchant_notification(
	id INT NOT NULL AUTO_INCREMENT,
	ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP UNIQUE,
	merchant_id VARCHAR(64) NOT NULL,
	activity_type INT,
	source_id INT NOT NULL,
	unread BOOL NOT NULL DEFAULT true,
	message VARCHAR(225) NOT NULL DEFAULT '',
	PRIMARY KEY (id),
	FOREIGN KEY fk_merchant_notification_type(activity_type)
		REFERENCES activity_type(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
);`

const schemaUserNotification = `CREATE TABLE user_notification(
	id INT NOT NULL AUTO_INCREMENT,
	ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP UNIQUE,
	user_id VARCHAR(64) NOT NULL,
	activity_type INT,
	source_id int NOT NULL,
	merchant_id VARCHAR(64) DEFAULT '' NULL,
	unread BOOL NOT NULL DEFAULT true,
	message VARCHAR(225) NOT NULL DEFAULT '',
	PRIMARY KEY (id),
	FOREIGN KEY fk_user_notification_type(activity_type)
		REFERENCES activity_type(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
);`

//User
const schemaUser = `
	CREATE TABLE user(
		id VARCHAR(64) NOT NULL,
		phone_number text NOT NULL,
		email VARCHAR(64) NOT NULL UNIQUE,
		name text NOT NULL,
		password text NOT NULL,
		facebook_id varchar(64) UNIQUE,
		freez_point INT NOT NULL DEFAULT 0,
		image VARCHAR(255) NOT NULL DEFAULT 'https://freeze-app.s3.us-west-2.amazonaws.com/blank-profile-picture.jpg',
		last_location POINT,
		PRIMARY KEY (id)
	);
`

const schemaMOption = `
	CREATE TABLE m_option(
		id INT AUTO_INCREMENT,
		user_id VARCHAR(64) NOT NULL,
		notif_fav_nearby BOOLEAN NOT NULL DEFAULT TRUE,
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
		PRIMARY KEY(user_id, merchant_id),
		FOREIGN KEY fk_user_fav (user_id)
			REFERENCES user(id)
			ON UPDATE CASCADE
			ON DELETE CASCADE,
		FOREIGN KEY fk_user_fav_merchant (merchant_id)
			REFERENCES merchant(id)
			ON UPDATE CASCADE
	);
`
const schemaRequest = `CREATE TABLE request(
		id INT NOT NULL AUTO_INCREMENT,
		user_id VARCHAR(64) NOT NULL,
		merchant_id VARCHAR(64) NOT NULL,
		location POINT NOT NULL,
		SPATIAL INDEX(location),
		comment VARCHAR(200) NOT NULL DEFAULT '',
		active BOOL NOT NULL DEFAULT TRUE,
		accepted TINYINT(1) DEFAULT -1 NOT NULL,
		PRIMARY KEY (id)
	);
`
