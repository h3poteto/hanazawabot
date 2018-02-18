package logging

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/johntdyer/slackrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	log.Hooks.Add(&slackrus.SlackrusHook{
		HookURL:        os.Getenv("SLACK_URL"),
		AcceptedLevels: slackrus.LevelThreshold(logrus.ErrorLevel),
		Channel:        "#hanazawabot",
		IconEmoji:      ":bapho:",
		Username:       "logrus",
	})

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

func (self *Logging) MethodInfoWithStacktrace(pkg string, err error) *logrus.Entry {
	stackErr, ok := err.(Stacktrace)
	if !ok {
		panic("oops, err does not implement Stacktrace")
	}
	st := stackErr.Stacktrace()
	traceLength := len(st)
	if traceLength > 5 {
		traceLength = 5
	}

	return self.Log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"pkg":        pkg,
		"stacktrace": fmt.Sprintf("%+v", st[0:traceLength]),
	})
}

// PanicRecover send error and stacktrace
func (u *Logging) PanicRecover() *logrus.Entry {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	return u.Log.WithFields(logrus.Fields{
		"time":       time.Now(),
		"pkg":        "main",
		"stacktrace": string(buf),
	})
}
