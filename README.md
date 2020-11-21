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
- assuming day is count by calendar date rather than 24h cycle
- assuming week is counted by calendar week, following the ISO week date system

TODO -
add mutex to db
add error handling when id is taken
move conversion of load request format to separate function/file
add checks forloadlimiting
add tests

add checks that
