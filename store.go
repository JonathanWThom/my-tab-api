package main

import (
	"database/sql"
	"github.com/jonathanwthom/my-tab-api/stddrink"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Store interface {
	CreateDrink(drink *Drink) (*Drink, error)
	CreateUser(user *User) (*User, error)
	DeleteDrink(id string) error
	GetDrinks(start, end string) ([]*Drink, error)
	LoginUser(user *User) (*User, error)
}

type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateUser(user *User) (*User, error) {
	var uuid string
	var username string
	var id int

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, err
	}
	sqlStatement := `
		INSERT INTO users(username, password)
		VALUES($1, $2)
		RETURNING uuid, username, id
	`

	// TODO: Should validate uniqueness of name
	err = store.db.QueryRow(
		sqlStatement,
		user.Username,
		string(hashedPassword)).Scan(&uuid, &username, &id)
	if err != nil {
		return nil, err
	}

	user.UUID = uuid
	user.Username = username
	user.ID = id

	return user, nil
}

func (store *dbStore) LoginUser(user *User) (*User, error) {
	storedUser := User{}

	sqlStatement := `SELECT password, uuid, id FROM users WHERE username=$1;`

	row := store.db.QueryRow(sqlStatement, user.Username)
	err := row.Scan(&storedUser.Password, &storedUser.UUID, &storedUser.ID)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return nil, err
	}

	return &storedUser, nil
}

func (store *dbStore) CreateDrink(drink *Drink) (*Drink, error) {
	drink.Stddrink = stddrink.Calculate(drink.Percent, drink.Oz)
	var id, dbUserID int
	var percent, oz, stddrink float64
	var imbibedOn time.Time
	var name sql.NullString

	sqlStatement := `
		INSERT INTO drinks(percent, oz, stddrink, imbibed_on, user_id, name)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, percent, oz, stddrink, imbibed_on, user_id, name`

	err := store.db.QueryRow(sqlStatement,
		drink.Percent,
		drink.Oz,
		drink.Stddrink,
		drink.ImbibedOn,
		userID,
		drink.Name).Scan(&id, &percent, &oz, &stddrink, &imbibedOn, &dbUserID, &name)
	if err != nil {
		return nil, err
	}

	drink.ID = id
	drink.Percent = percent
	drink.Oz = oz
	drink.Stddrink = stddrink
	drink.ImbibedOn = imbibedOn
	drink.UserID = dbUserID
	if name.Valid {
		drink.Name = name.String
	}
	return drink, err
}

func (store *dbStore) GetDrinks(start, end string) ([]*Drink, error) {
	var rows *sql.Rows
	var err error

	if start == "" || end == "" {
		sqlStatement := `
			SELECT id, percent, oz, stddrink, imbibed_on, user_id, name
			FROM drinks
			WHERE user_id = $1
			ORDER BY imbibed_on DESC
		`
		rows, err = store.db.Query(sqlStatement, userID)
	} else {
		times, err := stringsToTimes([]string{start, end})
		if err != nil {
			return nil, err
		}

		sqlStatement := `
			SELECT id, percent, oz, stddrink, imbibed_on, user_id, name
			FROM drinks
			WHERE imbibed_on
			BETWEEN $1 AND $2
			AND
			user_id = $3
			ORDER BY imbibed_on DESC
		`
		rows, err = store.db.Query(sqlStatement, times[0], times[1], userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	drinks := []*Drink{}
	for rows.Next() {
		drink := &Drink{}
		var name sql.NullString
		if err := rows.Scan(
			&drink.ID,
			&drink.Percent,
			&drink.Oz,
			&drink.Stddrink,
			&drink.ImbibedOn,
			&drink.UserID,
			&name); err != nil {
			return nil, err
		}

		if name.Valid {
			drink.Name = name.String
		}
		drinks = append(drinks, drink)
	}

	return drinks, nil
}

func (store *dbStore) DeleteDrink(id string) error {
	sqlStatement := `
		DELETE FROM drinks
		WHERE id = $1
		AND user_id = $2
	`
	_, err := store.db.Query(sqlStatement, id, userID)
	if err != nil {
		return err
	}
	// should something be returned here?
	return err
}

var store Store

func InitStore(s Store) {
	store = s
}
