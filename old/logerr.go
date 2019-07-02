package store

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Logerr struct {
	*logrus.Logger
}

// Logging is off by default, the caller can turn it on by passing
// us a logger (hmm. perhaps an io.Writer or a channel?).
var (
	l  *Logerr
	lg *logrus.Logger = nil
)

func SetLogerr(ll *logrus.Logger) {
	lg = ll // Logging is On!
}

func (l *Logerr) Println(args ...string) {
	if l != nil {
		l.Errorln(args)
	}
}

// IndexError will create an Error for the caller
func (l *Logerr) IndexError(arg string, msg string) error {
	errmsg := "Index Error"
	err := fmt.Errorf("%s %s %s", errmsg, arg, msg)
	l.Print(err.Error())
	return err
}

// IndexError will create an Error for the caller
func (l *Logerr) JSONError(arg string, msg string) error {
	errmsg := "JSON Error"
	err := fmt.Errorf("%s %s %s", errmsg, arg, msg)
	l.Print(err)
	return err
}

// IndexError will create an Error for the caller
func (l *Logerr) FetchError(arg string, msg string) error {
	errmsg := "Fetch Error"
	err := fmt.Errorf("%s %s %s", errmsg, arg, msg)
	l.Print(err)
	return err
}
