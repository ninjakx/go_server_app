package queries

const (
	CreateDB = `
		CREATE TABLE IF NOT EXISTS Servers (
			id SERIAL PRIMARY KEY,
			ip TEXT NOT NULL,
			hostname TEXT NOT NULL,
			active BOOLEAN NOT NULL
		);
	`
	QueryInsertServerData = `
		INSERT INTO Servers (ip, hostname, active) VALUES(:ip, :hostname, :active);
	`

	QueryGetAllHostnameWithThresh = `
		SELECT hostname as Hostnames
		FROM Servers
		GROUP BY hostname
		HAVING COUNT(CASE WHEN active THEN 1 END)<=$1;
	`

	QueryFindServer = `SELECT * FROM Servers WHERE id=$1;`

	QueryAllserver = `SELECT * FROM Servers;`

	QueryUpdateServer = `UPDATE Servers SET ip=:ip, hostname=:hostname, active=:active WHERE id=:id;`

	QueryUpdateServerStatus = `UPDATE Servers SET active=:active WHERE id=:id;`

	QueryDeleteServer = `DELETE FROM Servers WHERE id=$1;`
)
