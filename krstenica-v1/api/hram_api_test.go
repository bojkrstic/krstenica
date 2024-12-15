package api

import (
	"fmt"
	"krstenica/pkg/apiutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

var testPathRegistry *apiutil.PathRegistry

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Greška prilikom dobijanja home direktorijuma:", err)
		return
	}
	confFilePath := filepath.Join(homeDir, "develop", "horisen", "Krstenica-new", "krstenica", "config", "krstenica_api_conf.json")
	// Ispis putanje
	fmt.Println("Putanja do konfiguracione datoteke:", confFilePath)
	// confFilePath := "/home/bojan/develop/horisen/krstenica-new/krstenica/config/krstenica_api_conf.json"
	// confFilePath := "/home/krle/develop/horisen/krstenica-new/krstenica/config/krstenica_api_conf.json"
	// confFilePath := "$HOME/develop/horisen/Krstenica-new/krstenica/config/krstenica_api_conf.json"
	//need create pathregistry
	pathRegistry := createPathRegistry(confFilePath)
	c, err := pathRegistry.Config.GetConf()
	if err != nil {
		log.Fatal(err)

	}
	// connections to the database
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.PostgresSQL.Username, c.PostgresSQL.Password, c.PostgresSQL.Server, c.PostgresSQL.Database)
	initializeDatabases(connString)

	testPathRegistry = pathRegistry

}

func TestUserPostFirst(t *testing.T) {
	//randomInt := rand.Int()
	rand.Seed(uint64(time.Now().UnixNano())) //first we need to set up seed , than we get random number
	randomInt := rand.Intn(1000)

	req := HramWo{
		NazivHrama: fmt.Sprintf("Bojan Krstic %d", randomInt),
	}
	var usr HramCrtResWo
	err := apiutil.PerformApiTest(testPathRegistry, "POST", "/hram", req, &usr, nil)
	if err != nil {
		t.Fatal(err)
	}

}
func TestUserGetFirst(t *testing.T) {
	// randomInt := rand.Int()
	sid := strconv.Itoa(int(9))

	var response HramCrtResWo
	err := apiutil.PerformApiTest(testPathRegistry, "GET", "/hram/"+sid, nil, &response, nil)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(response)

}

func TestUserLists(t *testing.T) {
	// randomInt := rand.Int()
	//sid := strconv.Itoa(int(9))

	var response HramCrtResWo
	err := apiutil.PerformApiTest(testPathRegistry, "GET", "/hram", nil, &response, nil)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(response)

}

func TestUserDelete(t *testing.T) {
	// randomInt := rand.Int()
	sid := strconv.Itoa(int(5))

	var response HramCrtResWo
	err := apiutil.PerformApiTest(testPathRegistry, "DELETE", "/hram/"+sid, nil, &response, nil)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(response)

}

func TestUserUpdate(t *testing.T) {
	// randomInt := rand.Int()
	sid := strconv.Itoa(int(2))
	nazivHrama := "Sveti Sava"
	status := "active"
	updateData := HramUpdateData{
		NazivHrama: &nazivHrama,
		Status:     &status,
	}

	// updateData := map[string]interface{}{
	// 	"status":      "active",
	// 	"naziv_hrama": "Sveti Sava",
	// }

	var response HramCrtResWo
	err := apiutil.PerformApiTest(testPathRegistry, "PUT", "/hram/"+sid, &updateData, &response, nil)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(response)

}
