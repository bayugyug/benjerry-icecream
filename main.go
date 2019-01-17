package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/bayugyug/benjerry-icecream/config"
	"github.com/bayugyug/benjerry-icecream/controllers"
)

const (
	//VersionMajor main ver no.
	VersionMajor = "0.1"
	//VersionMinor sub  ver no.
	VersionMinor = "0"
)

var (
	//BuildTime pass during build time
	BuildTime string
	//ApiVersion is the app ver string
	ApiVersion string
	//Settings of the app
	Settings *config.ApiSettings
)

//internal system initialize
func init() {

	//uniqueness
	rand.Seed(time.Now().UnixNano())
	ApiVersion = "Ver: " + VersionMajor + "." + VersionMinor + "-" + BuildTime

}

func main() {

	start := time.Now()
	log.Println(ApiVersion)

	var err error

	//init
	Settings = config.NewAppSettings()

	//check
	if Settings.Config == nil {
		log.Fatal("Oops! Config missing")
	}

	//init service
	if controllers.ApiInstance, err = controllers.NewApiService(
		controllers.WithSvcOptAddress(":"+Settings.Config.HttpPort),
		controllers.WithSvcOptDbConf(&Settings.Config.Driver),
		controllers.WithSvcOptDumpFile(Settings.Config.DumpFile),
	); err != nil {
		log.Fatal("Oops! config might be missing", err)
	}
	//prep data
	controllers.ApiInstance.PrepareData()

	//run service
	controllers.ApiInstance.Run()
	log.Println("Since", time.Since(start))
	log.Println("Done")
}
