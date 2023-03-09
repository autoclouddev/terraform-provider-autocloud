package logger

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
)

/*

example of usage

log := Create(log.Fields{"resource": "test()"}) or logger.Create(map[string]interface{}{"resource": "test"})
log.Info("upload")
log.Info("upload complete")
log.Warn("upload retry")
log.WithError(errors.New("unauthorized")).Error("upload failed")
log.WithFields(log.Fields{"method": "getSomething()"}).info("getting something")
*/

func Create(fields log.Fields) *log.Entry {
	level, err := log.ParseLevel(os.Getenv("TF_LOG")) // using same levels as terraform
	if err != nil {
		level = 1 //info
	}

	var defaultHandler log.Handler

	jsonOutput := json.New(os.Stdout)
	cliOutput := cli.New(os.Stdout)

	if os.Getenv("LOG_HANDLER") == "JSON" {
		defaultHandler = jsonOutput
	} else {
		defaultHandler = cliOutput
	}

	logger := log.Logger{
		Handler: defaultHandler,
		Level:   level,
	}

	log := logger.WithFields(fields)
	return log
}
