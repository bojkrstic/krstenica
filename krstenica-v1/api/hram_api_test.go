package api

import (
	"fmt"
	"krstenica/pkg/apiutil"
	"log"
	"strconv"
	"testing"
)

var testPathRegistry *apiutil.PathRegistry

func init() {
	confFilePath := "/home/krle/develop/horisen/Krstenica/Go/Krstenica/krstenica/doc/user_api_conf.json"

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
	// randomInt := rand.Int()

	req := HramWo{
		NazivHrama: "Crkva Presvete Bogorodice",
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
		t.Fatal(err)
	}
	t.Log(response)

}
