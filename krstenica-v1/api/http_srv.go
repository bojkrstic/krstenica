package api

import (
	"fmt"
	"log"

	"krstenica/krstenica-v1/dao"
	"krstenica/pkg/apiutil"
)

var db dao.UserDao

func initializeDatabases(connString string) {
	cpd := dao.NewHramDao(connString)
	db = cpd
}

func RunHTTPServer(configFilePath string) {
	pathRegistry := createPathRegistry(configFilePath)
	c, err := pathRegistry.Config.GetConf()
	if err != nil {
		log.Fatal()
	}
	//connections to the database
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.PostgresSQL.Username, c.PostgresSQL.Password, c.PostgresSQL.Server, c.PostgresSQL.Database)

	fmt.Println("Connection String:", connString)
	initializeDatabases(connString)

	// server
	httpsrv, err := pathRegistry.NewDefaultHTTPSrv(configFilePath)
	if err != nil {
		log.Fatal()
	}
	log.Fatal(httpsrv.ListenAndServe())

}

func createPathRegistry(configFilePath string) *apiutil.PathRegistry {
	pathRegistry, config := apiutil.NewPathRegistry(configFilePath)
	log.Println("Prefix", config.URIPrefix)
	// pathRegistry.Map(config.URIPrefix+"/users", apiutil.POST, new(UserAdd))
	pathRegistry.Map(config.URIPrefix+"/hram", apiutil.POST, new(HramAdd))
	pathRegistry.Map(config.URIPrefix+"/hram/$id", apiutil.GET, new(HramGet))
	pathRegistry.Map(config.URIPrefix+"/hram/$id", apiutil.DELETE, new(HramDelete))
	// pathRegistry.Map(config.URIPrefix+"/hram/$id", apiutil.PUT, new(HramUpdate))
	// pathRegistry.Map(config.URIPrefix+"/hram/$id", apiutil.DELETE, new(HramDelete))
	//eparhija
	// pathRegistry.Map(config.URIPrefix+"/eparhija", apiutil.POST, new(EparhijaAdd))
	return pathRegistry
}
