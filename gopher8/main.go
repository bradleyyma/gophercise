package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "bma"
	password = ""
	dbname   = "gophercise8"
)

type phone struct {
	id     int
	number string
}

func main() {
	// First connect to postgres database to create/reset gophercise8
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("host=%s port=%d user=%s sslmode=disable",
		host, port, user))
	if err != nil {
		log.Fatal(err)
	}

	if err := resetDB(conn, dbname); err != nil {
		log.Fatal("Error resetting database:", err)
	}
	conn.Close(context.Background())

	// Now connect to the newly created gophercise8 database
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	conn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	if err := createPhoneNumberTable(conn); err != nil {
		panic(err)
	}

	_, err = insertPhone(conn, "1234567890")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(conn, "123 456 7891")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(conn, "(123) 456 7992")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(conn, "(123) 456-7893")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(conn, "123-456-7894")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(conn, "123-456-7980")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(conn, "1234567892")
	if err != nil {
		log.Fatal(err)
	}
	_, err = insertPhone(conn, "(123)456-7892")
	if err != nil {
		log.Fatal(err)
	}

	phone, err := getPhone(conn, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Phone number with ID 3:", phone)

	phones, err := allPhones(conn)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range phones {
		fmt.Printf("Phone: %s\n", p.number)
		number := normalize(p.number)
		fmt.Printf("Normalized number: %s\n", number)
		if number != p.number {
			fmt.Println("Inserting normalized number...")
			existingID, err := findPhone(conn, number)
			if err != nil {
				panic(err)
			}
			if existingID == -1 {
				fmt.Println("Inserting normalized number:", number)
				updatePhone(conn, p.id, number)
			} else {
				fmt.Println("Normalized number already exists, deleting this entry.")
				deletePhone(conn, p.id)
			}
		} else {
			fmt.Println("Number already normalized, skipping insert.")
		}
	}
}

func resetDB(db *pgx.Conn, name string) error {
	_, err := db.Exec(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", name))
	if err != nil {
		return err
	}
	_, err = db.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", name))
	return err
}

func createPhoneNumberTable(db *pgx.Conn) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (id SERIAL PRIMARY KEY, phone VARCHAR(255))
	`
	_, err := db.Exec(context.Background(), statement)
	return err
}

func getPhone(db *pgx.Conn, id int) (string, error) {
	var phone string
	err := db.QueryRow(context.Background(),
		"SELECT phone FROM phone_numbers WHERE id=$1", id).Scan(&phone)
	if err != nil {
		return "", err
	}
	return phone, nil
}

func findPhone(db *pgx.Conn, number string) (int, error) {
	var id int
	err := db.QueryRow(context.Background(),
		"SELECT id FROM phone_numbers WHERE phone=$1", number).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	return id, nil
}

func updatePhone(db *pgx.Conn, id int, newNumber string) error {
	_, err := db.Exec(context.Background(),
		"UPDATE phone_numbers SET phone=$1 WHERE id=$2", newNumber, id)
	return err
}

func deletePhone(db *pgx.Conn, id int) error {
	_, err := db.Exec(context.Background(),
		"DELETE FROM phone_numbers WHERE id=$1", id)
	return err
}

func allPhones(db *pgx.Conn) ([]phone, error) {
	rows, err := db.Query(context.Background(), "SELECT id, phone FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phones []phone
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		phones = append(phones, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return phones, nil
}

func insertPhone(db *pgx.Conn, phone string) (int, error) {
	var id int
	err := db.QueryRow(context.Background(),
		"INSERT INTO phone_numbers (phone) VALUES ($1) RETURNING id", phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			buf.WriteRune(r)

		}
	}
	return buf.String()
}
