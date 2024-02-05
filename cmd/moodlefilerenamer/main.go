package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
	"github.com/tamada/moodle_file_renamer"
	"github.com/tamada/moodle_file_renamer/utils"
)

const Version = "1.0.0"

func printHelp(progName string) string {
	name := filepath.Base(progName)
	return fmt.Sprintf(`%s version %s
Usage: %s [OPTIONS] <DIR>
OPTIONS
  -r, --restoration      restores the renamed directory names to the original.
  -f, --format <FORMAT>  specifies the format of the resultant directory names.
                         default: "%%default" (is equals to "%%uid_%%lname_%%fname")
                         available variables: default, original, uid, fname, lname, name, sid, and note.
  -h, --help             prints this help message and exit.
ARGUMENTS
  DIR    the target directory containing the downloaded directories from Moodle.`, name, Version, name)
}

func performImpl(opts *options, md moodle_file_renamer.MoodleDir) error {
	if opts.restore {
		return md.Restore()
	}
	return md.Rename(opts.format)
}

func perform(opts *options) error {
	var errs []error
	for _, dir := range opts.args {
		if !utils.IsDir(dir) {
			continue
		}
		md, err := moodle_file_renamer.Open(dir)
		if err != nil {
			return err
		}
		defer md.Close()
		if err := performImpl(opts, md); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func goMain(args []string) int {
	opts, err := parseOptions(args)
	if err != nil {
		return printError(1, err, args[0])
	}
	if opts.helpFlag {
		fmt.Println(printHelp(args[0]))
		return 0
	}
	if err := perform(opts); err != nil {
		return printError(1, err, args[0])
	}
	return 0
}

type options struct {
	restore  bool
	format   string
	helpFlag bool
	args     []string
}

func validate(opts *options) (*options, error) {
	if opts.helpFlag {
		return opts, nil
	}
	if len(opts.args) == 0 {
		return nil, fmt.Errorf("no directory is specified")
	}
	return opts, nil
}

func parseOptions(args []string) (*options, error) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(printHelp(args[0])) }
	opts := &options{}
	flags.BoolVarP(&opts.restore, "restore", "r", false, "restore the renamed directory names to the original")
	flags.StringVarP(&opts.format, "format", "f", "%default", "specify the format of the resultant directory names")
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this help message and exit")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}
	opts.args = flags.Args()
	return validate(opts)
}

func printError(status int, err error, progName string) int {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		fmt.Println(printHelp(progName))
	}
	return status
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
