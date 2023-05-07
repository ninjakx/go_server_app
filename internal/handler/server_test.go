package handler

import (
	"GO_APP/internal/model"
	"GO_APP/internal/queries"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// MockDB creates a mocked database and returns a sqlx.DB and a sqlmock.Sqlmock.
func MockDB() (*sqlx.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock, nil
}

func TestGetServerOr404(t *testing.T) {
	db, mock, err := MockDB()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()

	type args struct {
		id int
	}
	type field struct {
		server *model.Server
	}
	type testCases struct {
		name       string
		args       args
		field      field
		wantServer *model.Server
		wantError  error
		wantStatus int
		mock       func(*model.Server)
	}

	cols := []string{"id", "hostname", "ip"}
	tests := []testCases{
		{
			name: "Existing Server",
			args: args{
				id: 10,
			},
			field: field{
				server: &model.Server{
					Id:       10,
					Hostname: "Server1",
					IP:       "192.168.0.1",
				},
			},
			wantServer: &model.Server{
				Id:       10,
				Hostname: "Server1",
				IP:       "192.168.0.1",
			},

			wantError:  nil,
			wantStatus: http.StatusOK,
			mock: func(serverDet *model.Server) {
				rows := sqlmock.NewRows(cols).AddRow(serverDet.Id, serverDet.Hostname, serverDet.IP)
				mock.ExpectQuery("SELECT \\* FROM Servers WHERE id=\\$1").WithArgs(10).WithArgs(serverDet.Id).WillReturnRows(rows)
			},
		},
		{
			name: "Non Existing Server",
			args: args{
				id: 10,
			},
			field: field{
				server: &model.Server{
					Id:       10,
					Hostname: "Server1",
					IP:       "192.168.0.1",
				},
			},
			wantServer: &model.Server{},

			wantError:  errors.New("sql: no rows in result set"),
			wantStatus: http.StatusNotFound,
			mock: func(serverDet *model.Server) {
				mock.ExpectQuery("SELECT \\* FROM Servers WHERE id=\\$1").WithArgs(10).WithArgs(serverDet.Id).WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.field.server)
			w := httptest.NewRecorder()

			server, err := getServerOr404(db, tt.args.id, w)

			assert.Equal(t, tt.wantServer, server)
			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetServerHostName(t *testing.T) {
	db, mock, err := MockDB()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer db.Close()
	type args struct {
		ps httprouter.Params
	}
	type field struct {
	}
	type testCases struct {
		name       string
		args       args
		field      field
		want       []string
		wantStatus int
		mock       func()
	}

	tests := []testCases{
		{
			name: "Get server hostname : status code -> 200",
			args: args{
				ps: httprouter.Params{httprouter.Param{Key: "thresh", Value: "50"}},
			},
			field:      field{},
			want:       []string{"mta-prod-1", "mta-prod-2"},
			wantStatus: http.StatusOK,
			mock: func() {
				rows := sqlmock.NewRows([]string{"hostname"}).AddRow("mta-prod-1").AddRow("mta-prod-2")
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queries.QueryGetAllHostnameWithThresh)).WithArgs(50).WillReturnRows(rows)
				mock.ExpectCommit()
			},
		},
		{
			name: "Get server hostname : bad request on thresh pass default value",
			args: args{
				ps: httprouter.Params{httprouter.Param{Key: "thresh", Value: ""}},
			},
			field:      field{},
			want:       []string{"mta-prod-1", "mta-prod-2"},
			wantStatus: http.StatusOK,
			mock: func() {
				rows := sqlmock.NewRows([]string{"hostname"}).AddRow("mta-prod-1").AddRow("mta-prod-2")
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queries.QueryGetAllHostnameWithThresh)).WithArgs(1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
		},
		// {
		// 	name: "Get server hostname : status code -> 404",
		// 	args: args{
		// 		ps: httprouter.Params{httprouter.Param{Key: "thresh", Value: "50"}},
		// 	},
		// 	field: field{},
		// 	want:  []string{},

		// 	wantStatus: http.StatusNotFound,
		// 	mock: func() {
		// 		mock.ExpectBegin()
		// 		mock.ExpectQuery(regexp.QuoteMeta(queries.QueryGetAllHostnameWithThresh)).WithArgs(50).WillReturnError(sql.ErrTxDone)
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			w := httptest.NewRecorder()

			GetServerHostName(db, w, tt.args.ps)
			// Check the response status code
			if w.Code != tt.wantStatus {
				t.Errorf("Expected status code %v but got %v", tt.wantStatus, w.Code)
			}

			// Check the response body
			expectedBody, _ := json.Marshal(tt.want)
			if !bytes.Equal(w.Body.Bytes(), expectedBody) {
				t.Errorf("Expected response body %s but got %s", expectedBody, w.Body.String())
			}

			// Check the mock expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetServer(t *testing.T) {
	// Create a new mock database and HTTP request/response
	db, mock, err := MockDB()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}
	defer db.Close()

	w := httptest.NewRecorder()

	// Create expected server
	expectedServer := &model.Server{Id: 10, Hostname: "Test Server"}

	// Set up database mock to return the expected server
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM Servers WHERE id=$1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "hostname"}).
			AddRow(expectedServer.Id, expectedServer.Hostname))

	// Call the function being tested
	GetServer(db, w, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}})

	// Check that the response status code is correct
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code: got %d, expected %d", w.Code, http.StatusOK)
	}

	// Check that the response body is correct
	var server model.Server
	err = json.NewDecoder(w.Body).Decode(&server)
	if err != nil {
		t.Errorf("error decoding response body: %s", err)
	}
	if !reflect.DeepEqual(server, *expectedServer) {
		t.Errorf("unexpected server: got %+v, expected %+v", server, *expectedServer)
	}

	// Check that there were no unexpected database interactions
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestGetAllServer(t *testing.T) {
	// Create a new mock database and HTTP request/response
	db, mock, err := MockDB()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}
	defer db.Close()

	// req := httptest.NewRequest("GET", "/servers/1", nil)
	w := httptest.NewRecorder()

	// Create expected server
	expectedServer := []model.Server{
		{
			Id: 10, Hostname: "Test Server", IP: "127.0.0.1", Active: true,
		},
		{
			Id: 11, Hostname: "Test Server", IP: "127.0.0.1", Active: true,
		},
	}

	cols := []string{"id", "hostname", "ip", "active"}
	// Set up database mock to return the expected server
	rows := sqlmock.NewRows(cols)
	for i := range expectedServer {
		rows.AddRow(expectedServer[i].Id, expectedServer[i].Hostname, expectedServer[i].IP, expectedServer[i].Active)
	}

	mock.ExpectQuery(regexp.QuoteMeta(queries.QueryAllserver)).WillReturnRows(rows)

	// Call the function being tested
	GetAllServer(db, w)

	// Check that the response status code is correct
	if w.Code != http.StatusOK {
		t.Errorf("unexpected status code: got %d, expected %d", w.Code, http.StatusOK)
	}

	// Check that the response body is correct
	var server []model.Server
	err = json.NewDecoder(w.Body).Decode(&server)
	if err != nil {
		t.Errorf("error decoding response body: %s", err)
	}
	if !reflect.DeepEqual(server, expectedServer) {
		t.Errorf("unexpected server: got %+v, expected %+v", server, expectedServer)
	}

	// Check that there were no unexpected database interactions
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestCreateServer(t *testing.T) {
	db, mock, _ := MockDB()
	defer db.Close()
	server := model.Server{
		Id:       1,
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   true,
	}

	// Expect a single insert query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := httprouter.New()
	router.POST("/servers/create", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		CreateServer(db, w, r)
	})

	reqBody, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/servers/create", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateServer(t *testing.T) {
	db, mock, _ := MockDB()
	defer db.Close()
	server := model.Server{
		Id:       1,
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(server.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "ip", "hostname", "active"}).AddRow(server.Id, server.IP, server.Hostname, server.Active))
	mock.ExpectExec("UPDATE").WithArgs(server.IP, server.Hostname, server.Active, server.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := httprouter.New()
	router.PUT("/servers/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		UpdateServer(db, w, r, ps)
	})

	reqBody, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "/servers/1", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDisableServer(t *testing.T) {
	db, mock, _ := MockDB()
	defer db.Close()
	server := model.Server{
		Id:       1,
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   false,
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(server.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "ip", "hostname", "active"}).AddRow(server.Id, server.IP, server.Hostname, server.Active))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE Servers SET active=? WHERE id=?")).WithArgs(server.Active, server.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := httprouter.New()
	router.PUT("/servers/:id/disable", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		DisableServer(db, w, r, ps)
	})

	reqBody, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "/servers/1/disable", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEnableServer(t *testing.T) {
	db, mock, _ := MockDB()
	defer db.Close()
	server := model.Server{
		Id:       1,
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(server.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "ip", "hostname", "active"}).AddRow(server.Id, server.IP, server.Hostname, server.Active))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE Servers SET active=? WHERE id=?")).WithArgs(server.Active, server.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := httprouter.New()
	router.PUT("/servers/:id/enable", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		EnableServer(db, w, r, ps)
	})

	reqBody, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "/servers/1/enable", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteServer(t *testing.T) {
	db, mock, _ := MockDB()
	defer db.Close()
	server := model.Server{
		Id:       1,
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(server.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "ip", "hostname", "active"}).AddRow(server.Id, server.IP, server.Hostname, server.Active))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Servers WHERE id=$1;")).WithArgs(server.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := httprouter.New()
	router.DELETE("/servers/:id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		DeleteServer(db, w, r, ps)
	})

	reqBody, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("DELETE", "/servers/1", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
