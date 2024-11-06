package dao

import (
	"database/sql"
	"log"
)

func (c *HramDaoPostgresSql) CreateHram(user *HramDo) (uint, error) {
	c.Connect()
	defer c.Disconect()

	var id int
	err := c.db.QueryRow("insert into public.hram (naziv_hrama, created_at) VALUES ($1, $2) returning hram_id", user.HramID, user.CreatedAt).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return uint(id), nil
}

func (c *HramDaoPostgresSql) GetHram(id uint) (*HramDo, error) {
	c.Connect()
	defer c.Disconect()

	var hram HramDo
	err := c.db.QueryRow("select hram_id, naziv_hrama, created_at from public.hram where user_id = $1", id).Scan(&hram.HramID, &hram.HramName, &hram.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrHramNotFound
		}
		log.Println(err)
		return nil, err
	}

	return &hram, nil
}

func (c *HramDaoPostgresSql) GetHramByName(naziv string) (*HramDo, error) {
	c.Connect()
	defer c.Disconect()

	var hram HramDo

	err := c.db.QueryRow("select hram_id,naziv_hrama,created_at from public.hram where naziv_hrama = $1", naziv).Scan(&hram.HramID, &hram.HramName, &hram.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrHramNotFound
		}
		log.Println(err)
		return nil, err
	}

	return &hram, nil
}
