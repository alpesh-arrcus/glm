package config

import (
	"flag"
	//"fmt"
	"github.com/go-ini/ini"
	"github.com/rs/zerolog"
	. "github.com/toravir/glm/context"
	"log"
	"os"
)

var (
	dEFAULT_LOG_LEVEL = zerolog.ErrorLevel
)

type configState struct {
	iniCfg     *ini.File
	logger     zerolog.Logger
	listenAddr string
	dbSrc      string
	isSecure   bool
	serverKey  string // filename
	serverCrt  string // filename
}

func ParseCmdLineArgs() Context {
	ctx := CreateContext()
	cfgSt := configState{
		logger: zerolog.New(os.Stdout).Level(zerolog.ErrorLevel).With().Timestamp().Logger(),
	}
	ctx.Config = &cfgSt

	config := flag.String("config", "", "specify file which contains config for GLM")
	flag.Parse()

	if *config == "" {
		log.Fatal("Please specify a config file! Exiting..")
	}

	iniCfg, err := ini.Load(*config)
	if iniCfg == nil || err != nil {
		log.Fatal("Cannot load Config file:!!", err)
	}
	cfgSt.iniCfg = iniCfg
	parseLogConfig(&cfgSt)
	parseGlobalConfig(&cfgSt)
	parseDbConfig(&cfgSt)
	return ctx
}

func parseLogConfig(cfgSt *configState) {
	logCfg := cfgSt.iniCfg.Section("glm_logger")
	logLevelCfg := logCfg.Key("level").String()
	logDestCfg := logCfg.Key("destination").String()
	logLevel := dEFAULT_LOG_LEVEL
	if logLevelCfg != "" {
		var err error
		logLevel, err = zerolog.ParseLevel(logLevelCfg)
		if err != nil {
			//Invalid LoglevelCfg - so use default log level
			logLevel = dEFAULT_LOG_LEVEL
		}
	}
	logDest := os.Stdout
	switch logDestCfg {
	case "":
		fallthrough
	case "<stdout>":
		break
	case "<stderr>":
		logDest = os.Stderr
	default:
		fil, err := os.Create(logDestCfg)
		if err == nil {
			logDest = fil
		}
	}
	cfgSt.logger = zerolog.New(logDest).Level(logLevel).With().Timestamp().Logger()
}

func parseGlobalConfig(cfgSt *configState) {
	globalCfg := cfgSt.iniCfg.Section("global")
	saddr := globalCfg.Key("listenAddress").String()
	cfgSt.logger.Debug().Str("listenAddr", saddr).Msg("Loaded Config")
	if saddr == "" {
		//Crash and burn
		log.Fatal("Please specify listenAddress in the config file !")
	}
	isHttps, _ := globalCfg.Key("https").Bool()
	keyFile := globalCfg.Key("serverKey").String()
	crtFile := globalCfg.Key("serverCert").String()
	cfgSt.listenAddr = saddr
	cfgSt.isSecure = isHttps
	cfgSt.serverKey = keyFile
	cfgSt.serverCrt = crtFile
}

func parseDbConfig(cfgSt *configState) {
	dbCfg := cfgSt.iniCfg.Section("glm_database")
	dbSrc := dbCfg.Key("databaseName").String()
	cfgSt.dbSrc = dbSrc
}

func GetGLMListenAddress(ctx Context) string {
	cfgSt, _ := ctx.Config.(*configState)
	return cfgSt.listenAddr
}

func GetLogger(ctx Context) *zerolog.Logger {
	cfgSt, _ := ctx.Config.(*configState)
	return &cfgSt.logger
}

func GetDBSourceName(ctx Context) string {
	cfgSt, _ := ctx.Config.(*configState)
	return cfgSt.dbSrc
}

func GetHttpConfig(ctx Context) (isHttps bool, key, crt string) {
	cfgSt, _ := ctx.Config.(*configState)
	return cfgSt.isSecure, cfgSt.serverKey, cfgSt.serverCrt
}
