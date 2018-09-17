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
	GetDrinks() ([]*Drink, error)
	GetDrinksByDateRange(start, end time.Time) ([]*Drink, error)
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

	/// should validate uniqueness of name
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

	sqlStatement := `
		INSERT INTO drinks(percent, oz, stddrink, imbibed_on, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, percent, oz, stddrink, imbibed_on, user_id`

	err := store.db.QueryRow(sqlStatement,
		drink.Percent,
		drink.Oz,
		drink.Stddrink,
		drink.ImbibedOn,
		userID).Scan(&id, &percent, &oz, &stddrink, &imbibedOn, &dbUserID)
	if err != nil {
		return nil, err
	}

	drink.ID = id
	drink.Percent = percent
	drink.Oz = oz
	drink.Stddrink = stddrink
	drink.ImbibedOn = imbibedOn
	drink.UserID = dbUserID
	return drink, err
}

func (store *dbStore) GetDrinks() ([]*Drink, error) {
	sqlStatement := `
		SELECT id, percent, oz, stddrink, imbibed_on
		FROM drinks
		WHERE user_id = $1
	`
	rows, err := store.db.Query(sqlStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	drinks := []*Drink{}
	for rows.Next() {
		drink := &Drink{}
		if err := rows.Scan(&drink.ID, &drink.Percent, &drink.Oz, &drink.Stddrink, &drink.ImbibedOn); err != nil {
			return nil, err
		}

		drinks = append(drinks, drink)
	}

	return drinks, nil
}

func (store *dbStore) GetDrinksByDateRange(start, end time.Time) ([]*Drink, error) {
	// start := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)
	// end := time.Now()
	// drinks, err := store.GetDrinksByDateRange(start, end)
	//
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//
	// fmt.Println(drinks)

	sqlStatement := `
		SELECT id, percent, oz, stddrink, imbibed_on
		FROM drinks
		WHERE imbibed_on
		BETWEEN $1 AND $2
	`
	/// by user_id
	rows, err := store.db.Query(sqlStatement, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	drinks := []*Drink{}

	for rows.Next() {
		drink := &Drink{}
		if err := rows.Scan(
			&drink.ID,
			&drink.Percent,
			&drink.Oz,
			&drink.Stddrink,
			&drink.ImbibedOn); err != nil {
			return nil, err
		}

		drinks = append(drinks, drink)
	}

	return drinks, nil
}

var store Store

func InitStore(s Store) {
	store = s
}
