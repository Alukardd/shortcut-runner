package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

var (
	conf struct {
		Program     string `yaml:"program"`
		Args        string `yaml:"args"`
		Timeout     uint   `yaml:"timeout"`
		InputDevice string `yaml:"InputDevice"`
		Shortcut    string `yaml:"shortcut"`
	}
	configFile string
	l          *syslog.Writer
)

func init() {
	// parse command line arguments
	flag.StringVar(&conf.Program, "program", "/bin/true", "specify full path to executed program")
	flag.StringVar(&conf.Args, "args", "", "specify arguments for executed program")
	flag.UintVar(&conf.Timeout, "timeout", 180, "set timeout during which shortcut are expected")
	flag.StringVar(&conf.InputDevice, "inputdev", "/dev/input/event0", "specify input device path")
	flag.StringVar(&conf.Shortcut, "shortcut", "KEY_LEFTALT KEY_A", "specify arguments for executed program")
	flag.StringVar(&configFile, "config", "/etc/shortcut-runner/shortcut-runner.yml", "specify config file")

	flag.Parse()
}

func main() {
	var err error
	l, err = syslog.Dial("", "", syslog.LOG_DAEMON|syslog.LOG_INFO, "shortcut-runner")
	if err != nil {
		log.Fatalf("Failed to connect to syslog: %s\n", err.Error())
		os.Exit(1)
	}
	defer l.Info("exiting")
	defer l.Close()

	// read config, if exist
	var buf []byte
	if buf, err = ioutil.ReadFile(configFile); err == nil {
		if err = yaml.Unmarshal(buf, &conf); err != nil {
			l.Info("Config file is malformed")
		}
	} else {
		l.Info("Config file can not be read")
	}

	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			l.Info(fmt.Sprint("program was terminated by signal(", sig, ")"))
			l.Close()
			os.Exit(2)
		}
	}()

	// global programm timer
	go func() {
		time.Sleep(time.Duration(conf.Timeout) * time.Second)
		l.Info("program was terminated by timeout")
		l.Close()
		os.Exit(0)
	}()

	if err = waitShortcut(conf.InputDevice, strings.Split(conf.Shortcut, " ")); err == nil {
		l.Info(fmt.Sprint("Keystroke (", conf.Shortcut, ") was captured"))
		l.Info(fmt.Sprint("Executing a program: ", conf.Program, " ", conf.Args))
		cmd := exec.Command(conf.Program, conf.Args)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			l.Err(fmt.Sprint("Error running the program: ", err))
		}
	}
}
