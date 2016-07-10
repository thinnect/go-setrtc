// Author  Raido Pahtma
// License MIT

package setrtc

//import "fmt"
import "time"
import "bytes"
import "encoding/binary"

import "github.com/proactivity-lab/go-loggers"
import "github.com/proactivity-lab/go-sfconnection"

import "github.com/thinnect/go-devparam"

import "github.com/beevik/ntp"

const STARTUP_DELAY = 3 * time.Second

type RealTimeClockSetter struct {
	loggers.DIWEloggers

	sfc *sfconnection.SfConnection
	dpm *deviceparameters.DeviceParameterManager

	host string

	Exit chan bool
}

func NewRealTimeClockSetter(sfc *sfconnection.SfConnection, host string) *RealTimeClockSetter {
	ss := new(RealTimeClockSetter)
	ss.InitLoggers()

	ss.sfc = sfc
	ss.dpm = deviceparameters.NewDeviceParameterManager(sfc)

	ss.host = host

	ss.Exit = make(chan bool)

	return ss
}

func (self *RealTimeClockSetter) AnnounceTime(offset int64) error {
	oldval, err := self.dpm.GetValue("unix_time")
	if err != nil {
		self.Warning.Printf("Get current failed")
		return err
	}

	ntpr, err := ntp.Query(self.host, 4)
	if err != nil {
		self.Warning.Printf("NTP query failed %s", err)
		return err
	}
	self.Debug.Printf("NTP %d stratum %d RTT %s offset %s", ntpr.Time.Unix(), ntpr.Stratum, ntpr.RTT, ntpr.ClockOffset)

	t := int64(ntpr.Time.Unix() + offset)
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, t); err != nil {
		return err
	}

	newval, err := self.dpm.SetValue("unix_time", buf.Bytes())
	if err == nil {
		if offset != 0 {
			self.Warning.Printf("Using RTC offset %d\n", oldval, newval, offset)
		}
		self.Info.Printf("%s->%s\n", oldval, newval)
	} else {
		self.Warning.Printf("Setting RTC failed: %s", err)
	}

	return err
}

func (self *RealTimeClockSetter) Run(period time.Duration, offset int64) {
	self.Debug.Printf("run\n")
	self.AnnounceTime(offset)
	for {
		select {
		case <-self.Exit:
			self.Debug.Printf("Exit.\n")
			self.dpm.Close()
			self.sfc.Disconnect()
		case <-time.After(period):
			self.AnnounceTime(offset)
		}
	}
}
