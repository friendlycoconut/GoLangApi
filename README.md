# GoLang Api

In process of creation of this REST API were used: PostgreSQL, Go-Chi, 


## DB configuration

Use the change the connection settings of your DB

```go
package main

const (
	DB_USER     = "user"
	DB_PASSWORD = "password"
	DB_username = "dbname"
)
```

## Usage (Main class)

```go
package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	a  App
)

func main() {
	a = App{}
	a.Initialize()

	a.Run(":yourport")
}

```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
