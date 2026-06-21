[![Unit Coverage](https://coveralls.io/repos/github/Sugyk/auth_service/badge.svg?branch=main&flag=unit)](https://coveralls.io/github/Sugyk/auth_service)
[![Integration Coverage](https://coveralls.io/repos/github/Sugyk/auth_service/badge.svg?branch=main&flag=integration)](https://coveralls.io/github/Sugyk/auth_service)

## Problems

* **Tests with txmanager can be used only happy paths.** TestTxManager does not commit nothing, so tests with tx can not be used when several transaction used in service, especially when they are work with the same data, because tests can not rollback and test bad cases.

