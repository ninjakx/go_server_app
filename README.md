# mta-hosting-optimizer

Using gin, GORM, gocron and JWT tokens.
### API

```go
	router.GET("/servers/get_hostname/:thresh", a.GetServerHostname)
	router.GET("/servers", a.GetAllServer)
	router.GET("/server/:id", a.GetServer)
	router.POST("/servers/create", a.CreateServer)
	router.PUT("/servers/:id/update_server", a.UpdateServer)
	router.PUT("/servers/:id/disable", a.DisableServer)
	router.PUT("/servers/:id/enable", a.EnableServer)
	router.DELETE("/servers/:id", a.DeleteServer)
```

all the api with examples can be found under postman collection file.

### CURL

**Create server:**
```bash
curl --location 'http://localhost:8004/servers/create' \
--header 'Content-Type: text/plain' \
--data '{
	"Ip":"127.0.0.8",
	"Hostname":"mta-prod-5",
	"Active": false
}'
```

**Search server by id:**
```bash
curl --location 'http://localhost:8004/server/2'
```

### RUN:

`go run main.go`

**To continuously connect to the application server, run the following command**

#### to run server:
nodemon --exec go run cmd/server/main.go --signal SIGTERM
#### to run cron:
nodemon --exec go run cmd/cron/cronjob.go --signal SIGTERM

**To run test:**

with coverage:

`go test ./... -cover` 

## PostgresSQL DB:

**To make them sorted and in an order by ID:**

```bash
UPDATE servers m
SET id = sub.rn
from (SELECT id, row_number() OVER (ORDER BY id, id) AS rn FROM servers)sub
WHERE  m.id = sub.id;
```

**DB:**

![](https://github.com/ninjakx/go_server_app/blob/bc43e9c47ee3533fbb7b37994aaa5125821be6c9/Images4Readme/psql_db.png?raw=true)

**API for getting hostnames with threshold:**

![](https://github.com/ninjakx/go_server_app/blob/main/Images4Readme/query_thresh.png?raw=true)

**Code coverage:**

Due to logging the error: Unit test coverage can't be more than `61.6%`. It's good to log the error so introducing those will reduce the coverage as can be seen from the image (red part-> not covered in unit test)

![](https://github.com/ninjakx/go_server_app/blob/main/Images4Readme/code_coverage.png?raw=true)
