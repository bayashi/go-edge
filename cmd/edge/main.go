package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

const (
	cmd string = "edge"

	exitOK  int = 0
	exitErr int = 1
)

type options struct {
	file           string
	grep           []string
	showTotalCount bool
}

var (
	version     = ""
	installFrom = "Source"
)

func main() {
	err := Run()
	if err != nil {
		putErr(fmt.Sprintf("Err: %s", err))
		os.Exit(exitErr)
	}
	os.Exit(exitOK)
}

func putErr(message ...interface{}) {
	fmt.Fprintln(os.Stderr, message...)
}

func putUsage() {
	putErr(fmt.Sprintf("Usage: %s [OPTIONS] FILE", cmd))
}

func putHelp(message string) {
	putErr(message)
	putUsage()
	putErr("Options:")
	flag.PrintDefaults()
	os.Exit(exitOK)
}

func Run() error {
	o := &options{}

	o.parseArgs()
	o.getTargetFile()

	err := o.pickLogs()
	if err != nil {
		return err
	}

	return nil
}

func (o *options) pickLogs() error {
	f, err := os.Open(o.file)
	if err != nil {
		return errors.Wrap(err, "Could not open file")
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	c := 0
	total := 0
	needGrep := len(o.grep) > 0
	var (
		line     string
		lastLine *string
	)

	for s.Scan() {
		line = s.Text()
		total++
		if needGrep && !o.isMatchedLine(&line) {
			continue
		}
		c++
		if c == 1 {
			fmt.Println(fmt.Sprintf("%d: %s", c, line))
		} else {
			lastLine = &line
		}
	}

	if lastLine != nil {
		fmt.Println(fmt.Sprintf("%d: %s", c, *lastLine))
	}

	if s.Err() != nil {
		return errors.Wrap(s.Err(), "Happened error during reading file")
	}

	if o.showTotalCount {
		plural := "s"
		if total == 1 {
			plural = ""
		}
		fmt.Println(fmt.Sprintf("total: %d line%s", total, plural))
	}

	return nil
}

func (o *options) isMatchedLine(line *string) bool {
	bline := []byte(*line)
	for _, word := range o.grep {
		if matched, _ := regexp.Match(regexp.QuoteMeta(word), bline); matched {
			return true
		}
	}

	return false
}

func (o *options) parseArgs() {
	var flagHelp bool
	var flagVersion bool
	flag.BoolVarP(&flagHelp, "help", "h", false, "Show help (This message) and exit")
	flag.BoolVarP(&flagVersion, "version", "v", false, "Show version and build info and exit")

	flag.StringArrayVarP(&o.grep, "grep", "g", []string{}, "Search PATTERN in each line")
	flag.BoolVarP(&o.showTotalCount, "c", "c", false, "Show tatl count of a file")

	flag.Parse()

	if flagHelp {
		putHelp(fmt.Sprintf("[%s] Version %s", cmd, getVersion()))
	}

	if flagVersion {
		putErr(versionDetails())
		os.Exit(exitOK)
	}
}

func versionDetails() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	compiler := runtime.Version()

	return fmt.Sprintf(
		"Version %s - %s.%s (compiled:%s, %s)",
		getVersion(),
		goos,
		goarch,
		compiler,
		installFrom,
	)
}

func getVersion() string {
	if version != "" {
		return version
	}
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return "Unknown"
	}

	return i.Main.Version
}

func (o *options) getTargetFile() {
	for _, arg := range flag.Args() {
		if o.file != "" {
			putHelp(fmt.Sprintf("Err: Wrong args. Unnecessary arg [%s]", arg))
		}
		if arg == "-" {
			continue
		}
		o.file = arg
	}

	if o.file == "" {
		putHelp("Err: Wrong args. You should specify a FILE")
	}
}
