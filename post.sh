# Create user in test DB
#curl -d "DB=test&ID=2&Firstname=hans" -H "Content-Type: application/json" -X POST http://192.168.2.70:3000/api/user

curl -d '{"DBName": "test", "ID": "2", "Firstname": "Mr", "Lastname": "Andersen"}' -H "Content-Type: application/json" -X POST http://192.168.2.70:3000/api/user

