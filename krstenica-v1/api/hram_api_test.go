package api

import (
	"fmt"
	"krstenica/pkg/apiutil"
	"log"
	"strconv"
	"testing"

	"golang.org/x/exp/rand"
)

var testPathRegistry *apiutil.PathRegistry

func init() {
	confFilePath := "/home/krle/develop/horisen/Krstenica-new/krstenica/config/krstenica_api_conf.json"

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
	randomInt := rand.Int()
	// randomInt := rand.Intn(1000)

	req := HramWo{
		NazivHrama: fmt.Sprintf("Ilija Krstic %d", randomInt),
	}
	var usr HramCrtResWo
	err := apiutil.PerformApiTest(testPathRegistry, "POST", "/hram", req, &usr, nil)
	if err != nil {
		t.Fatal(err)
	}

}
func TestUserGetFirst(t *testing.T) {
	// randomInt := rand.Int()
	sid := strconv.Itoa(int(1))

	var response HramCrtResWo
	err := apiutil.PerformApiTest(testPathRegistry, "GET", "/hram/"+sid, nil, &response, nil)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(response)

}
