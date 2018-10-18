package initialization

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func logrusFatalHandler() {
	// gracefully shutdown something...
}
func init() {
	//txtFormatter := log.TextFormatter{}
	logFormatter := log.JSONFormatter{}
	logFormatter.PrettyPrint = true
	log.SetFormatter(&logFormatter)
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	log.RegisterExitHandler(logrusFatalHandler)

	// Only log the warning severity or above.
	log.SetLevel(log.ErrorLevel)
}
