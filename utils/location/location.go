package location

import "time"

type location struct {
	pst *time.Location
	bjs *time.Location
}

var loc location

// New the constructor function but can not be used directly
// It must be registered in plugins by plugin.Register() and Loaded at begining
func New() *location {
	loc = location{}
	return &loc
}

func (l *location) Load() error {
	var err error
	l.pst, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return err
	}
	l.bjs, err = time.LoadLocation("Asia/Chongqing")
	if err != nil {
		return err
	}
	return nil
}

func (l location) Loaded() bool {
	return l.bjs != nil
}

func (l location) Name() string {
	return "hrpc-location"
}

func (l location) DependsOn() []string {
	return nil
}

// PST can be used for LA
func PST() *time.Location {
	return loc.pst
}

// BJS can be used for CN
func BJS() *time.Location {
	return loc.bjs
}
