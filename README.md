<p align="right">
 <img src="https://github.com/hellgate75/go-services/workflows/Go/badge.svg?branch=master"></img>
&nbsp;&nbsp;<img src="https://api.travis-ci.com/hellgate75/go-services.svg?branch=master" alt="trevis-ci" width="98" height="20" />&nbsp;&nbsp;<a href="https://travis-ci.com/hellgate75/go-services">Check last build on Travis-CI</a>
 </p>

# go-services

Services library including drivers for some infrastructure dependencies.

## Service

Database drivers:

* Model - Database infrastructure interfaces
* MySQL - Database drivers that allows to connect with MySql database servers
* MongoDB - Database drivers that allows to connect with MongoDb database servers

### Model

Base components are :

* [Driver](/database/database.go) - Allows multiple service connection
* [DriverConfig](/database/database.go) - Allows service connection configuration
* [Connection](/database/database.go) - Represents the service connection instance


### MySQL

Instance will is provided by `GetDatabaseDriver` or `GetDatabaseDriverByName`, it accepts the database.MySQLDriver
variable or with the driver name `mysql`.



### MongoDB

Instance will is provided by `GetDatabaseDriver` or `GetDatabaseDriverByName`, it accepts the database.MongoDbDriver
variable or with the driver name `mongodb`.


### Get the library

Library is available running:

```
go get -u github.com/hellgate75/go-services
```


## License

The library is licensed with [LGPL v. 3.0](/LICENSE) clauses, with prior authorization of author before any production or commercial use. Use of this library or any extension is prohibited due to high risk of damages due to improper use. No warranty is provided for improper or unauthorized use of this library or any implementation.

Any request can be prompted to the author [Fabrizio Torelli](https://www.linkedin.com/in/fabriziotorelli) at the following email address:

[hellgate75@gmail.com](mailto:hellgate75@gmail.com)
 
