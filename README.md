# Auth Microservice

---

### Features

* Register users in PostgreSQL
* Set JWT-tokens to users
* Add logins to the Redis blacklist to protect it from bruteforce

### Technologies

* Redis
* PostgreSQL
* Docker Compose

### Testing

Testing implemented with python. To start testing on ready service run `python3 test.py`.

Or to start testing from zero just run script `startTest.sh`. It will rebuild compose services, up them and run test.py