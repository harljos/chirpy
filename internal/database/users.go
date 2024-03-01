package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

var ErrAlreadyExists = errors.New("already exists")

func (db *DB) CreateUser(email string, password string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		Email:    email,
		ID:       id,
		Password: hashed,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(id int, email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user.Email = email
	user.Password = hashed
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) LoginUser(email string, password string) (User, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserById(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, ErrNotExist
}
