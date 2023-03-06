package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

const (
	cmd     string = "edge"
	version string = "0.0.1"

	exitOK  int = 0
	exitErr int = 1
)

type options struct {
	file           string
	grep           []string
	showTotalCount bool
}

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
		line string
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
	flag.BoolVarP(&flagHelp, "help", "h", false, "Show help (This message) and exit")

	flag.StringArrayVarP(&o.grep, "grep", "g", []string{}, "Search PATTERN in each line")
	flag.BoolVarP(&o.showTotalCount, "c", "c", false, "Show tatl count of a file")

	flag.Parse()

	if flagHelp {
		putHelp(fmt.Sprintf("[%s] Version v%s", cmd, version))
	}
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
