package main

import (
	"redis_test/RedisOpt"
	"vislog"

	log "github.com/Sirupsen/logrus"
)

func main() {
	var redisOpt RedisOpt.RedisOpt
	redisOpt.InitCluster([]string{}, "")
}

func initLog(logfile string, loglevel string, syslogAddr string) {
	if logfile == "" {
		logfile = "taskinterface.log"
	}

	hook, err := vislog.NewVislogHook(logfile)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.AddHook(hook)

	level, err := log.ParseLevel(loglevel)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	log.SetFormatter(&log.JSONFormatter{})

}
