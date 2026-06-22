![Unit Coverage](https://coveralls.io/repos/github/ВАШ_ЮЗЕР/РЕПО/badge.svg?flag=unit)
![Integration Coverage](https://coveralls.io/repos/github/ВАШ_ЮЗЕР/РЕПО/badge.svg?flag=integration)

## Problems

* **Tests with txmanager can be used only happy paths.** TestTxManager does not commit nothing, so tests with tx can not be used when several transaction used in service, especially when they are work with the same data, because tests can not rollback and test bad cases.

