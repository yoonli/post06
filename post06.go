package post06

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type UserData struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

var (
	Hostname = ""
	Port     = 0
	Username = ""
	Password = ""
	Database = ""
)

func openConnection() (*sql.DB, error) {
	//connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)
	//open db connection
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Function to return the User ID of the username
func exists(username string) int {
	username = strings.ToLower(username)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	userID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "users" WHERE "username" = '%s'`, username)
	rows, err := db.Query(statement)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		userID = id
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)
	return userID
}

// function to add a user to the database
// Return new userid
// -1 if there was an error
func AddUser(d UserData) int {
	d.Username = strings.ToLower(d.Username)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists: ", userID)
	}
	insertStatement := fmt.Sprintf(`INSERT INTO "users" ("username") values ($1)`)
	_, err = db.Exec(insertStatement, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	userID = exists(d.Username)
	if userID != -1 {
		return userID
	}
	insertStatement = `insert into "userdata" ("userid", "name", "surname", "description") values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}
	return userID
}

// DeleteUser Function to delete user
func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	// Does the user exist
	statement := fmt.Sprintf(`SELECT "username" FROM "users" where id = %d`, id)
	rows, err := db.Query(statement)
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	if exists(username) != id {
		return fmt.Errorf("user with ID %d does not exist", id)
	}
	// Delete from Userdata
	deleteStatement := `delete from "userdata" where userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	deleteStatement = `delete from "users" where id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func ListUsers() ([]UserData, error) {
	Data := []UserData{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)
	rows, err := db.Query(`SELECT 
 						"id","username","name","surname","description"
 						FROM "users","userdata"
 						WHERE users.id = userdata.userid`)
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if err != nil {
		return Data, err
	}
	for rows.Next() {
		var id int
		var username string
		var name string
		var description string
		err = rows.Scan(&id, &username, &name, &description)
		temp := UserData{ID: id, Username: username, Name: name, Description: description}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	return Data, nil
}
