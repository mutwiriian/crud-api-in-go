# CRUD API IN Go
To run the code locally run
```
git clone https://github.com/mutwiriian/crud-api-in-go.git
```

## Examples 
Try the following commands on a Unix command line to interact with the API

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
