CREATE TABLE Persons
(
P_Id int NOT NULL,
LastName varchar(255) NOT NULL,
FirstName varchar(255),
Address varchar(255),
City varchar(255),
CHECK (P_Id>0)
)

ALTER TABLE Persons
ADD CHECK (P_Id>0)

ALTER TABLE Persons
DROP CONSTRAINT chk_Person