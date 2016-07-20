package logging

import (
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

type Logging struct {
	Log *logrus.Logger
}

type Stacktrace interface {
	Stacktrace() []errors.Frame
}

var sharedInstance *Logging = New()

func New() *Logging {
	log := logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	return &Logging{Log: log}
}

func SharedInstance() *Logging {
	return sharedInstance
}

func (self *Logging) MethodInfo(pkg string) *logrus.Entry {
	return self.Log.WithFields(logrus.Fields{
		"time": time.Now(),
		"pkg":  pkg,
	})
}
