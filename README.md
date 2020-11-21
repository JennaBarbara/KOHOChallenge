Mysql Database structure:

 CREATE DATABASE loadfundsapp;

CREATE TABLE loaded_funds (
    Id int NOT NULL,
    Customer_id int NOT NULL,
    Load_amount DECIMAL(10,2) NOT NULL,
    Time datetime NOT NULL
);


Assumptions
- going off given times for rate limits
- assuming all the requests are coming from the same time zone
- assuming all load amounts are using the same currency
- assuming day is count by calendar date rather than 24h cycle
- assuming week is counted by calendar week starting with Monday
- assuming all load requests are unique
TODO -
add tests
add unique database primary key
add delete utility function for tests
