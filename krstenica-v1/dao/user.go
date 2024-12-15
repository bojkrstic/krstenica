package dao

import (
	"database/sql"
	"fmt"
	"log"
)

func (c *HramDaoPostgresSql) CreateHram(user *HramDo) (uint, error) {
	c.Connect()
	defer c.Disconect()

	var id int
	err := c.db.QueryRow("insert into public.hram (naziv_hrama, status, created_at) VALUES ($1, $2, $3) returning hram_id", user.HramName, user.Status, user.CreatedAt).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return uint(id), nil
}

func (c *HramDaoPostgresSql) UpdateHram(id uint, update map[string]interface{}) error {
	c.Connect()
	defer c.Disconect()
	args := []interface{}{}
	ind := 1
	//formiranje querija
	query := fmt.Sprintf("UPDATE public.hram SET")

	newStatus, ok := update["status"]
	if ok {
		args = append(args, newStatus)
		query += fmt.Sprintf(" status=$%d,", ind)
		ind++
	}
	nazivHrama, ok := update["naziv_hrama"]
	if ok {
		args = append(args, nazivHrama)
		query += fmt.Sprintf(" naziv_hrama=$%d,", ind)
		ind++
	}
	args = append(args, id)

	if len(args) == 0 {
		return fmt.Errorf("Nema polja za azuriranje")
	}

	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE hram_id = $%d", ind)
	log.Println(ind)
	log.Println(query)
	log.Println(args)
	_, err := c.db.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return ErrHramDubleValue
	}

	return nil

}

func (c *HramDaoPostgresSql) GetHram(id uint) (*HramDo, error) {
	c.Connect()
	defer c.Disconect()

	var hram HramDo
	err := c.db.QueryRow("select hram_id, naziv_hrama, status, created_at from public.hram where hram_id = $1", id).Scan(&hram.HramID, &hram.HramName, &hram.Status, &hram.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrHramNotFound
		}
		log.Println(err)
		return nil, err
	}

	return &hram, nil
}
func (c *HramDaoPostgresSql) DeleteHram(id uint) (*HramDo, error) {
	c.Connect()
	defer c.Disconect()

	sqlq := fmt.Sprintf("UPDATE public.hram SET status= $1 WHERE hram_id = $2")
	var hram HramDo
	_, err := c.db.Exec(sqlq, "deleted", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrHramNotFound
		}
		log.Println(err)
		return nil, err
	}

	err = c.db.QueryRow("select hram_id, naziv_hrama, status, created_at from public.hram where hram_id = $1", id).Scan(&hram.HramID, &hram.HramName, &hram.Status, &hram.CreatedAt)
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
