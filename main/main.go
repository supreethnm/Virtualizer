package main

import (
	"log"
	"net/http"
	"virtualizer/db"
	"virtualizer/route"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"

	c "virtualizer/configuration"
	cn "virtualizer/constants"
	r "virtualizer/route"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	var config c.Config
	if _, err := toml.DecodeFile(cn.CONFIG_FILE, &config); err != nil {
		logrus.WithFields(logrus.Fields{}).Panic(err.Error())
		panic(err)
	} //For TOML

	var dbConfig c.Config
	if _, err := toml.Decode(db.GetDBEndpoints(), &dbConfig); err != nil {
		logrus.WithFields(logrus.Fields{}).Panic(err.Error())
		panic(err)
	} //For TOML

	// config + dbConfig
	config.Services = append(config.Services, dbConfig.Services...)

	route.InitializeRoutes(config.Services) //For TOML
	router := r.NewRouter()
	logrus.WithFields(logrus.Fields{}).Info("Listening on port " + cn.PORT)
	log.Fatal(http.ListenAndServe(":"+cn.PORT, router))

}
