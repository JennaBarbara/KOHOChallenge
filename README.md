
## About

This project is written in go version 1.15.5. It connects to a MySql database using gorm.

## Before You Start

As this program requires a MySQL database, you will need a MySql table with the following structure in a MySql Database running locally to run the program:

```
CREATE TABLE loaded_funds (
    Id int NOT NULL,
    Customer_id int NOT NULL,
    Load_amount DECIMAL(10,2) NOT NULL,
    Time datetime NOT NULL,
    Accepted BOOLEAN NOT NULL
);
```

Please specify the details of your local MySql database (such as port, name, username, and password) in config/config.go.

You can also specify the intended input and output files for your program in config/config.go.

## Running the Program

Running the unit tests:
```
go tests
```

Running the program:
```
go run main.go
```

## Assumptions

The following assumptions were made when designing this program:
- Load requests should be maintained between uses of the program, hence the use of a database
- It is appropriate to go off the times given in the input for velocity rate limits
- All the inputs are coming from the same time zone
- All load amounts are using the same currency
- Daily velocity limits are counted by calendar date rather than 24h cycle
- Weekly velocity limits are counted by calendar week, with each week starting with Monday, rather than 7 day cycle
- Load requests with the same Id and Customer_id as a previous request should be ignored, even if that request was not accepted
