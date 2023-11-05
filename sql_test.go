package belajar_golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id,name) VALUES('joko','Joko')"
	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id,name FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}

	//iterate to get data on rows
	for rows.Next() {
		var id, name string // initial variable to get data

		// to get data, parameter send must pointer. and also this function return error to check.
		err = rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("Id", id)     // print id
		fmt.Println("Name", name) // print name
	}

	defer rows.Close()
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id,name,email,balance,rating,birth_date,married,created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}

	//iterate to get data on rows
	for rows.Next() {
		var id, name string      // initial variable to get data
		var email sql.NullString // initial variable if column database can insert NULL value
		var balance int32
		var rating float64
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool
		// to get data, parameter send must pointer. and also this function return error to check.
		err = rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("===============")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
		if email.Valid {
			fmt.Println("Email:", email.String)
		}
		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)
		if birthDate.Valid {
			fmt.Println("Birth Date:", birthDate.Time)
		}
		fmt.Println("Married:", married)
		fmt.Println("Created At:", createdAt)
	}

	defer rows.Close()
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// dummy data is like from user
	username := "admin'; #" // sql injection
	password := "salah"

	script := "SELECT username FROM user WHERE username='" + username + "' AND password='" + password + "'"
	fmt.Println(script) // print SQL query
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var username string

		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Login Success", username) // login success

	} else {
		fmt.Println("Login Fail") // login failed
	}

	defer rows.Close()
}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// dummy data is like from user
	username := "admin'; #" // sql injection
	password := "salah"

	script := "SELECT username FROM user WHERE username=? AND password=?"
	fmt.Println(script) // print SQL query
	rows, err := db.QueryContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var username string

		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Login Success", username) // login success

	} else {
		fmt.Println("Login Fail") // login failed
	}

	defer rows.Close()
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "rahmat"
	password := "rahmat"

	script := "INSERT INTO user(username,password) VALUES(?,?)"
	_, err := db.ExecContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new user")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "rahmat@rahmat.com"
	comment := "test comment"

	script := "INSERT INTO comments(email,comment) VALUES(?,?)"
	result, err := db.ExecContext(ctx, script, email, comment)
	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new comment with id", insertId)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "INSERT INTO comments(email,comment) VALUES(?,?)"
	statement, err := db.PrepareContext(ctx, script)

	if err != nil {
		panic(err)
	}

	defer statement.Close()

	for i := 1; i <= 10; i++ {
		email := "rahmat" + strconv.Itoa(i) + "@gmail.com"
		comment := "Comment to " + strconv.Itoa(i)

		// execution the statement
		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()

		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id", id)

	}

}

func TestTransaction(t *testing.T) {
	// make connection
	db := GetConnection()
	defer db.Close() // close connection on end

	ctx := context.Background() // make context

	tx, err := db.Begin() // start transaction mode

	if err != nil {
		panic(err)
	}

	script := "INSERT INTO comments(email,comment) VALUES(?,?)"

	for i := 1; i <= 10; i++ {
		email := "budi" + strconv.Itoa(i) + "@gmail.com"
		comment := "Comment budi to " + strconv.Itoa(i)

		// execution the statement
		result, err := tx.ExecContext(ctx, script, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()

		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id", id)

	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
