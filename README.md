# CRUD API in Go
To run the code locally:

First,ensure you have a Postgres server running in your system and execute the following:
```
create database customers;
```
Then run 
```
git clone https://github.com/mutwiriian/crud-api-in-go.git
cd crud-api-in-go
```
Create the environment variables
```
export DBHOST=host       #127.0.0.1
export DBPORT=port       #5432 
export DBUSER=user       #postgres
export DBPASS=password   #postgres
export DBNAME=name       #customers
```
Finally run the program
```
go run main.go
```

## Examples 
Try the following at the terminal to interact with the API

## CREATE Customer
```
curl -X POST http://localhost:8000/customers/create 
-d '{"name":"imma","email":"imma@yahoo.com","phone_number":"0743267419","address":"Nairobi"}'
```

## READ all Customers
```
curl -X GET http://localhost:8000/customers/get_all
```

## READ a Customer
```
curl -X GET http://localhost:8000/customers/get_email?email=imma@yahoo.com
```

## UPDATE a Customer
```
curl -X PATCH http://localhost:8000/customers/update?email=imma@yahoo.com -d '{"name":"munene"}'
```

## DELETE a Customer
```
curl -X DELETE http://localhost:8000/customers/delete_email?email=imma@yahoo.com
```
