package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bayugyug/benjerry-icecream/driver"
	"github.com/bayugyug/benjerry-icecream/utils"
)

const (
	//status
	usageConfig = "use to set the config file parameter with db-userinfos/http-port"
)

var (
	//Settings of the app
	Settings *ApiSettings
)

//internal system initialize
func init() {
	//uniqueness
	rand.Seed(time.Now().UnixNano())
}

//ParameterConfig optional parameter structure
type ParameterConfig struct {
	Driver   driver.DbConnectorConfig `json:"driver"`
	HttpPort string                   `json:"http_port"`
	Showlog  bool                     `json:"showlog"`
	DumpFile string                   `json:"dump_file"`
}

//AppSettings app mapping on its config
type ApiSettings struct {
	Config    *ParameterConfig
	CmdParams string
	EnvVars   map[string]*string
}

type Setup func(*ApiSettings)

func WithSetupConfig(r *ParameterConfig) Setup {
	return func(args *ApiSettings) {
		args.Config = r
	}
}

func WithSetupCmdParams(r string) Setup {
	return func(args *ApiSettings) {
		args.CmdParams = r
	}
}

func WithSetupEnvVars(r map[string]*string) Setup {
	return func(args *ApiSettings) {
		args.EnvVars = r
	}
}

//NewAppSettings main entry for config
func NewAppSettings(setters ...Setup) *ApiSettings {
	//set default
	cfg := &ApiSettings{
		EnvVars: make(map[string]*string),
	}
	//maybe export from envt
	cfg.EnvVars = map[string]*string{
		"BENJERRY_ICECREAM_CONFIG": &cfg.CmdParams,
	}
	//chk the passed params
	for _, setter := range setters {
		setter(cfg)
	}
	//start
	cfg.Initializer()
	return cfg
}

//InitRecov is for dumpIng segv in
func (g *ApiSettings) InitRecov() {
	//might help u
	defer func() {
		recvr := recover()
		if recvr != nil {
			fmt.Println("MAIN-RECOV-INIT: ", recvr)
		}
	}()
}

//InitEnvParams enable all OS envt vars to reload internally
func (g *ApiSettings) InitEnvParams() {
	//just in-case, over-write from ENV
	for k, v := range g.EnvVars {
		if os.Getenv(k) != "" {
			*v = os.Getenv(k)
		}
	}
	//get options
	flag.StringVar(&g.CmdParams, "config", g.CmdParams, usageConfig)
	flag.Parse()
}

//Initializer set defaults for initial reqmts
func (g *ApiSettings) Initializer() {
	//prepare
	g.InitRecov()
	g.InitEnvParams()
	log.Println("CmdParams:", g.CmdParams)

	//try to reconfigure if there is passed params, otherwise use show err
	if g.CmdParams != "" {
		g.Config = g.FormatParameterConfig(g.CmdParams)
	}

	//check defaults
	if g.Config == nil {
		return
	}
	//set dump flag
	utils.ShowMeLog = g.Config.Showlog

}

//FormatParameterConfig new ParameterConfig
func (g *ApiSettings) FormatParameterConfig(s string) *ParameterConfig {
	var cfg ParameterConfig
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		log.Println("FormatParameterConfig", err)
		return nil
	}
	return &cfg
}
