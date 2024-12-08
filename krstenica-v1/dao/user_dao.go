package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type UserDao interface {
	CreateHram(user *HramDo) (uint, error)
	GetHram(id uint) (*HramDo, error)
	DeleteHram(id uint) (*HramDo, error)
	GetHramByName(name string) (*HramDo, error)
}

// here we have connection to PostgresSql
type HramDaoPostgresSql struct {
	db         *sql.DB
	connString string
}

func NewHramDao(connectionString string) UserDao {
	return &HramDaoPostgresSql{
		connString: connectionString,
	}
}

func (u *HramDaoPostgresSql) Connect() error {
	fmt.Println("Connection string **************************", u.connString)
	db, err := sql.Open("postgres", u.connString)
	if err != nil {
		log.Println(err)
		return err
	}

	u.db = db

	return nil
}
func (u *HramDaoPostgresSql) Disconect() {
	u.db.Close()
}
