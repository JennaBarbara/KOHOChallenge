Mysql Database structure:

 CREATE DATABASE loadfundsapp;

CREATE TABLE LoadedFunds (
    Id int NOT NULL,
    Customer_id int NOT NULL,
    Load_amount DECIMAL(10,2) NOT NULL,
    Time datetime NOT NULL,
    PRIMARY KEY (Id)
);


Assumptions
- going off given times for rate limits
- assuming all the requests are coming from the same time zone
- assuming all load amounts are using the same currency
