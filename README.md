## This project is skeleton for golang project follow Clean Architecture  

Generate mock-repo  
```
make mock-repo
```

Run unit test  
```
make test
```

### Create transaction  
POST http://localhost:50051/api/users/1/transactions
```
{
  "account_id": 2,
  "amount": 100000.00,
  "transaction_type": "deposit"
}
```

### Update Transaction  
PUT http://localhost:50051/api/users/1/transactions/:transaction_id  
```
{
  "amount": 100000.00,
}
```

### Delete Transaction
DELETE http://localhost:50051/api/users/1/transactions/:transaction_id 

### Find Transaction
GET http://localhost:50051/api/users/1/transactions?account_id=2
