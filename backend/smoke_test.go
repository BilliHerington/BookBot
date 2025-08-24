package main

import (
	"database/sql"
	"net/http"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://postgres:22848@postgres_db:5432/mydatabase?sslmode=disable")
	if err != nil {
		t.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Fatalf("Postgress ping failed: %v", err)
	}
}

func TestHTTPServer(t *testing.T) {
	resp, err := http.Get("http://go_backend:8080/ping")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("failed to connect to http server: %v", err)
	}
}
