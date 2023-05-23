package handler

import (
	"GO_APP/internal/model"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/stretchr/testify/assert"
)

// MockDB creates a mocked database and returns a *gorm.DB and a sqlmock.Sqlmock.
func MockDB() (*gorm.DB, sqlmock.Sqlmock, *sql.DB, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, mockDB, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, mockDB, err
	}
	return gormDB, mock, mockDB, nil
}

func TestGetServerOr404(t *testing.T) {
	db, mock, dbmock, err := MockDB()
	defer dbmock.Close()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
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

	cols := []string{"ip", "hostname", "active"}
	tests := []testCases{
		{
			name: "Existing Server",
			args: args{
				id: 1,
			},
			field: field{
				server: &model.Server{
					IP:       "192.168.0.1",
					Hostname: "test-server",
					Active:   true,
				},
			},
			wantServer: &model.Server{
				Hostname: "test-server",
				IP:       "192.168.0.1",
				Active:   true,
			},

			wantError:  nil,
			wantStatus: http.StatusOK,
			mock: func(testServer *model.Server) {
				id := 1
				mock.ExpectQuery((`SELECT (.+) FROM "servers" WHERE id = (.+) LIMIT 1`)).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows(cols).
						AddRow(testServer.IP, testServer.Hostname, testServer.Active))
			},
		},
		{
			name: "Non Existing Server",
			args: args{
				id: 10,
			},
			field: field{
				server: &model.Server{
					Hostname: "Server1",
					IP:       "192.168.0.1",
					Active:   true,
				},
			},
			wantServer: nil,

			wantError:  errors.New("sql: no rows in result set"),
			wantStatus: http.StatusNotFound,
			mock: func(serverDet *model.Server) {
				mock.ExpectQuery((`SELECT (.+) FROM "servers" WHERE id = (.+) LIMIT 1`)).WithArgs(10).WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.field.server)
			// Create the custom response writer
			w := &testResponseWriter{httptest.NewRecorder()}

			// Create a new Gin context with the custom response writer
			c, _ := gin.CreateTestContext(w)
			server, err := getServerOr404(db, tt.args.id, c)

			assert.Equal(t, tt.wantServer, server)
			assert.Equal(t, tt.wantError, err)
		})
	}
}

func TestGetServerHostName(t *testing.T) {
	db, mock, dbmock, err := MockDB()
	defer dbmock.Close()

	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	type args struct {
		c *gin.Context
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
				c: &gin.Context{
					Params: gin.Params{gin.Param{Key: "thresh", Value: "50"}},
				},
			},
			field:      field{},
			want:       []string{"mta-prod-1", "mta-prod-2"},
			wantStatus: http.StatusOK,
			mock: func() {
				rows := sqlmock.NewRows([]string{"hostname"}).AddRow("mta-prod-1").AddRow("mta-prod-2")
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT (.+)").WithArgs(50).WillReturnRows(rows)
				mock.ExpectCommit()
			},
		},
		{
			name: "Get server hostname : bad request on thresh pass default value",
			args: args{
				c: &gin.Context{
					Params: gin.Params{gin.Param{Key: "thresh", Value: ""}},
				},
			},
			field:      field{},
			want:       []string{"mta-prod-1", "mta-prod-2"},
			wantStatus: http.StatusOK,
			mock: func() {
				rows := sqlmock.NewRows([]string{"hostname"}).AddRow("mta-prod-1").AddRow("mta-prod-2")
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT (.+)").WithArgs(1).WillReturnRows(rows)
				mock.ExpectCommit()
			},
		},
		// {
		// 	name: "Get server hostname : status code -> 404",
		// 	args: args{
		// 		c: &gin.Context{
		// 			Params: gin.Params{gin.Param{Key: "thresh", Value: ""}},
		// 		},
		// 	},
		// 	field:      field{},
		// 	want:       []string{"mta-prod-1", "mta-prod-2"},
		// 	wantStatus: http.StatusInternalServerError,
		// 	mock: func() {
		// 		mock.ExpectBegin()
		// 		mock.ExpectQuery("SELECT (.+)").WithArgs(1).WillReturnError(sql.ErrTxDone)
		// 		mock.ExpectCommit()
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			w := &testResponseWriter{httptest.NewRecorder()}
			// Create a new Gin context with the custom response writer
			c, _ := gin.CreateTestContext(w)
			c.Params = tt.args.c.Params

			GetServerHostName(db, c)
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
	db, mock, dbmock, err := MockDB()
	defer dbmock.Close()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}

	w := &testResponseWriter{httptest.NewRecorder()}

	// Create a new Gin context with the custom response writer
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	// Create expected server
	expectedServer := &model.Server{Hostname: "Test Server"}

	// Set up database mock to return the expected server
	mock.ExpectQuery(`SELECT (.+) FROM "servers" WHERE id = (.+) LIMIT 1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"hostname"}).
			AddRow(expectedServer.Hostname))

	// Call the function being tested
	GetServer(db, c)

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
	db, mock, dbmock, err := MockDB()
	defer dbmock.Close()
	if err != nil {
		t.Fatalf("error creating mock database: %s", err)
	}

	w := &testResponseWriter{httptest.NewRecorder()}

	// Create a new Gin context with the custom response writer
	c, _ := gin.CreateTestContext(w)

	// Create expected server
	expectedServer := []model.Server{
		{
			Hostname: "Test Server", IP: "127.0.0.1", Active: true,
		},
		{
			Hostname: "Test Server", IP: "127.0.0.1", Active: true,
		},
	}

	cols := []string{"hostname", "ip", "active"}
	// Set up database mock to return the expected server
	rows := sqlmock.NewRows(cols)
	for i := range expectedServer {
		rows.AddRow(expectedServer[i].Hostname, expectedServer[i].IP, expectedServer[i].Active)
	}

	mock.ExpectQuery(`SELECT (.+) FROM "servers"`).WillReturnRows(rows)

	// Call the function being tested
	GetAllServer(db, c)

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
	db, mock, dbmock, _ := MockDB()
	defer dbmock.Close()
	server := model.Server{
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   true,
	}

	cols := []string{"hostname", "ip", "active"}
	// Expect a single insert query
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT (.+)").
		WillReturnRows(
			sqlmock.NewRows(cols).
				AddRow(server.Hostname, server.IP, server.Active))
	mock.ExpectCommit()

	reqBody, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}
	// rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/servers/create", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	params := gin.Params{gin.Param{Key: "id", Value: "1"}}
	context := gin.Context{Request: req, Params: params}

	CreateServer(db, &context)

	// if rr.Code != http.StatusOK {
	// 	t.Errorf("Expected response code %d, but got %d", http.StatusOK, rr.Code)
	// }

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestUpdateServer(t *testing.T) {
	db, mock, dbmock, _ := MockDB()
	defer dbmock.Close()
	server := model.Server{
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   true,
	}
	id := 1

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"ip", "hostname", "active"}).AddRow(server.IP, server.Hostname, server.Active))
	mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body, err := json.Marshal(server)
	if err != nil {
		t.Fatalf("Error marshaling server: %v", err)
	}
	// rr := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/servers/1/update_server", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	params := gin.Params{gin.Param{Key: "id", Value: "1"}}
	context := gin.Context{Request: req, Params: params}

	UpdateServer(db, &context)

	// if rr.Code != http.StatusOK {
	// 	t.Errorf("Expected response code %d, but got %d", http.StatusOK, rr.Code)
	// }
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDisableServer(t *testing.T) {
	db, mock, dbmock, _ := MockDB()
	defer dbmock.Close()
	server := model.Server{
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   false,
	}

	id := 1

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"ip", "hostname", "active"}).AddRow(server.IP, server.Hostname, server.Active))
	mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body, err := json.Marshal(server)
	if err != nil {
		t.Fatalf("Error marshaling server: %v", err)
	}
	// rr := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/servers/1/disable", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	params := gin.Params{gin.Param{Key: "id", Value: "1"}}
	context := gin.Context{Request: req, Params: params}

	DisableServer(db, &context)

	// if rr.Code != http.StatusOK {
	// 	t.Errorf("Expected response code %d, but got %d", http.StatusOK, rr.Code)
	// }
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEnableServer(t *testing.T) {
	db, mock, dbmock, _ := MockDB()
	defer dbmock.Close()
	server := model.Server{
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   false,
	}

	id := 1
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"ip", "hostname", "active"}).AddRow(server.IP, server.Hostname, server.Active))
	mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body, err := json.Marshal(server)
	if err != nil {
		t.Fatalf("Error marshaling server: %v", err)
	}
	// rr := httptest.NewRecorder()
	req, err := http.NewRequest("PUT", "/servers/1/enable", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	params := gin.Params{gin.Param{Key: "id", Value: "1"}}
	context := gin.Context{Request: req, Params: params}

	EnableServer(db, &context)

	// if rr.Code != http.StatusOK {
	// 	t.Errorf("Expected response code %d, but got %d", http.StatusOK, rr.Code)
	// }
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteServer(t *testing.T) {
	db, mock, dbmock, _ := MockDB()
	defer dbmock.Close()
	server := model.Server{
		IP:       "192.168.1.1",
		Hostname: "test.com",
		Active:   true,
	}

	id := 1

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "ip", "hostname", "active"}).AddRow(id, server.IP, server.Hostname, server.Active))
	mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body, err := json.Marshal(server)
	if err != nil {
		t.Fatalf("Error marshaling server: %v", err)
	}
	// rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/servers/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	params := gin.Params{gin.Param{Key: "id", Value: "1"}}
	context := gin.Context{Request: req, Params: params}

	DeleteServer(db, &context)

	// if rr.Code != http.StatusOK {
	// 	t.Errorf("Expected response code %d, but got %d", http.StatusOK, rr.Code)
	// }
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
