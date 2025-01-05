# CRUD API IN GO
To run the code locally run
```
git clone https://github.com/mutwiriian/crud-api-in-go.git
```

## Examples 
Try the following commands on a Unix command line to interact with the API
```
curl -X POST http://localhost:8000/customers/create 
-d '{"name":"imma","email":"imma@yahoo.com","phone_number":"0743267419","address":"Nairobi"}'
```

```
curl -X GET http://localhost:8000/customers/get_all
```

```
curl -X GET http://localhost:8000/customers/get_email?email=imma@yahoo.com
```

```
curl -X PATCH http://localhost:8000/customers/update?email=imma@yahoo.com -d '{"name":"munene"}'
```

```
curl -X DELETE http://localhost:8000/customers/delete_email?email=imma@yahoo.com
```
