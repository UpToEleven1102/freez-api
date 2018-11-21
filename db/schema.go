package db

const schema = `
	CREATE TABLE merchant(
		id INT NOT NULL AUTO_INCREMENT,
		phone_number text NOT NULL,
		email text NOT NULL,
		name text NOT NULL,
		password text NOT NULL,
		image VARCHAR(50) NOT NULL DEFAULT './assets/profile.jpg',
		PRIMARY KEY (id)
	);
`
