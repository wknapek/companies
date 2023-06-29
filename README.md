# companies
To run application as a parameter provide configuration in JSON format.
## example of configuration 
```
{
  "user": "test",
  "password": "pwd",
  "dburi": "mongodb://localhost:27017",
  "sessionTime": 15,
  "port": "3001",
  "db_user": "companies",
  "db_passwd": "S3cret"
  "JWT_Key": "my_secret_key"
}

```
* user -> to authenticate via JWT
* password -> to authenticate via JWT
* JWT_Key -> to signs tokens
* port -> to set port for application
##
T0 execute automated tests run ```docker-compose up```