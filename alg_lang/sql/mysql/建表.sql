CREATE INDEX index_name
ON table_name (column_name);

CREATE UNIQUE INDEX index_name
ON table_name (column_name);

CREATE TABLE table_name
(
    id int NOT NULL AUTO_INCREMENT ,   -- PRIMARY KEY约束
    name varchar(255) NOT NULL,
    address varchar(255) DEFAULT 'gz',
    city varchar(255),
	PRIMARY KEY (`id`)
)ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
## check
CREATE TABLE Persons
(
    P_Id int NOT NULL,
    LastName varchar(255) NOT NULL,
    FirstName varchar(255),
    Address varchar(255),
    City varchar(255),
    CHECK (P_Id>0)
);

ALTER TABLE Persons
    ADD CHECK (P_Id>0);

ALTER TABLE Persons
    DROP CONSTRAINT chk_Person;

## 创建索引

CREATE INDEX index_name
ON table_name (column_name);

CREATE UNIQUE INDEX index_name
ON table_name (column_name);
