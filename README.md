# GO_APP
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

`nodemon --exec go run main.go --signal SIGTERM`

**To run test:**

with coverage:

`go test ./... -cover` 

## PostgresSQL DB:

**DB:**

![]()

