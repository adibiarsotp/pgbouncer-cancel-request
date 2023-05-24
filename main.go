package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

func test_db(ctx context.Context, db *sql.DB, q string) {

	result, err := db.ExecContext(ctx, q)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Rows affected", rows)
}

func testhit() {

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 100*time.Millisecond)

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:4000", nil)
	req = req.WithContext(ctx)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	fmt.Println("Response received, status code:", res.StatusCode)

}

func initServer() {
	var dsn = flag.String("dsn", "postgres://postgres@localhost:6432?dbname=db&sslmode=disable", "PostgreSQL DSN postgres://postgres:password@localhost:5432?sslmode=disable")
	var q = flag.String("q", "SELECT pg_sleep(2);", "Query to execute example: \"SELECT TRUE;\"")
	flag.Parse()
	fmt.Println("Using PostgreSQL DSN:", *dsn)
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		panic(err)
	}

	err = db.PingContext(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to DB")

	go func() {
		http.ListenAndServe(":4000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			fmt.Fprint(os.Stdout, "processing request\n")
			test_db(ctx, db, *q)
			w.Write([]byte("request processed"))
			fmt.Fprint(os.Stderr, "request processed\n")
		}))
	}()

}

func initClient() {
	var wg sync.WaitGroup
	for e := 0; e < 40; e++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10000; i++ {
				testhit()
			}
			wg.Done()
		}()
	}
	wg.Wait()

}

func main() {

	initServer()

	time.Sleep(6 * time.Second)

	initClient()

}
