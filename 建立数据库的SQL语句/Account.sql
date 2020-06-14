USE Food_Manager
CREATE TABLE Account(
			aid VARCHAR(10) NOT NULL UNIQUE,
			balance FLOAT(5) NOT NULL,
			level INT NOT NULL,
			PRIMARY KEY(aid));