USE Food_Manager
CREATE TABLE Sender(
			sid VARCHAR(10) NOT NULL UNIQUE,
			sname VARCHAR(30) NOT NULL,
			sphone VARCHAR(15) NOT NULL,
			aid VARCHAR(10) NOT NULL,
			PRIMARY KEY(sid),
			FOREIGN KEY (aid) REFERENCES Account ON DELETE NO ACTION);