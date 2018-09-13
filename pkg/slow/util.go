package slow

import (
	"github.com/golang/glog"
)

func runDelay(label string, d DelayRange) {
	glog.Infof("%s delay %s", label, d.ToString())
	d.Wait()
}
