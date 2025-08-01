package main

import (
	"database/sql"
	"fmt"
	"iter"
	"log"
)

type User struct{}

func doQuery(string) (sql.Rows, error)

func DoQuery[T any](query string) (iter.Seq2[T, error], error) {
	rows, err := doQuery(query)
	if err != nil {
		return nil, err
	}

	return func(yield func(T, error) bool) {
		defer rows.Close()
		for rows.Next() {
			var value T
			err := rows.Scan(&value)
			if !yield(value, err) {
				return
			}
		}
	}, nil
}

func main() {
	rows, err := DoQuery[User]("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}

	for user, err := range rows {
		fmt.Println(user, err)
	}
}
