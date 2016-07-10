// Author  Raido Pahtma
// License MIT

// setrtc executable
package main

import "fmt"
import "time"
import "os"
import "os/signal"

import "github.com/proactivity-lab/go-loggers"
import "github.com/proactivity-lab/go-sfconnection"

import "github.com/thinnect/go-setrtc"

import "github.com/jessevdk/go-flags"

const ApplicationVersionMajor = 0
const ApplicationVersionMinor = 1
const ApplicationVersionPatch = 0

var ApplicationBuildDate string
var ApplicationBuildDistro string

type Options struct {
	Positional struct {
		ConnectionString string `description:"Connectionstring sf@HOST:PORT"`
	} `positional-args:"yes"`

	NtpHost string `short:"n" long:"ntp-host" default:"0.pool.ntp.org" description:"NTP server address"`

	Offset int64 `short:"o" long:"offset" default:"0"   description:"Time offset"`
	Period int   `short:"p" long:"period" default:"900" description:"Announcement period"`

	Debug       []bool `short:"D" long:"debug"   description:"Debug mode, print raw packets"`
	ShowVersion func() `short:"V" long:"version" description:"Show application version"`
}

func mainfunction() int {

	var opts Options
	opts.ShowVersion = func() {
		if ApplicationBuildDate == "" {
			ApplicationBuildDate = "YYYY-mm-dd_HH:MM:SS"
		}
		if ApplicationBuildDistro == "" {
			ApplicationBuildDistro = "unknown"
		}
		fmt.Printf("setrtc %d.%d.%d (%s %s)\n", ApplicationVersionMajor, ApplicationVersionMinor, ApplicationVersionPatch, ApplicationBuildDate, ApplicationBuildDistro)
		os.Exit(0)
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		flagserr := err.(*flags.Error)
		if flagserr.Type != flags.ErrHelp {
			if len(opts.Debug) > 0 {
				fmt.Printf("Argument parser error: %s\n", err)
			}
			return 1
		}
		return 0
	}

	host, port, err := sfconnection.ParseSfConnectionString(opts.Positional.ConnectionString)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return 1
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)

	sfc := sfconnection.NewSfConnection()
	ss := setrtc.NewRealTimeClockSetter(sfc, opts.NtpHost)

	logger := loggers.BasicLogSetup(len(opts.Debug))
	if len(opts.Debug) > 0 {
		sfc.SetLoggers(logger)
	}
	ss.SetLoggers(logger)

	sfc.Autoconnect(host, port, 30*time.Second)

	time.Sleep(time.Second)

	go ss.Run(time.Duration(opts.Period)*time.Second, opts.Offset)

	for interrupted := false; interrupted == false; {
		select {
		case sig := <-signals:
			signal.Stop(signals)
			logger.Debug.Printf("signal %s\n", sig)
			sfc.Disconnect()
			interrupted = true
			ss.Exit <- true
		}
	}

	time.Sleep(100 * time.Millisecond)
	return 0
}

func main() {
	os.Exit(mainfunction())
}
