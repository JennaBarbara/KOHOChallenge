
##About

This project is written in go version 1.15.5. It connects to a MySql database using gorm.

It is designed to

##Before You Start

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

##Running the Program

Running the unit tests
```
go tests
```

Running the program

```
go run main.go
```

##Assumptions

The following assumptions were made when designing this program:
- going off given times for rate limits
- assuming all the requests are coming from the same time zone
- assuming all load amounts are using the same currency
- assuming day is count by calendar date rather than 24h cycle
- assuming week is counted by calendar week starting with Monday
- assuming that load requests with same Id and Customer_id as a previous request should be ignored, even if that request was not accepted
- assuming load requests should be maintained between uses of the program (hence the db)
