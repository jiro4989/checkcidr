package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type CmdArgs struct {
	Version    bool
	OutFile    string
	Style      string
	NoProgress bool
	Args       []string
}

const (
	helpMsgHelp       = "print help"
	helpMsgVersion    = "print version"
	helpMsgOutFile    = "output file"
	helpMsgStyle      = "printing style format [free_text | json | json_stream]"
	helpMsgNoProgress = "disable printing progress"
)

func ParseArgs() (*CmdArgs, error) {
	args := CmdArgs{}

	flag.Usage = flagHelpMessage
	flag.BoolVar(&args.Version, "v", false, helpMsgVersion)
	flag.StringVar(&args.OutFile, "o", "", helpMsgOutFile)
	flag.StringVar(&args.Style, "style", "free_text", helpMsgStyle)
	flag.BoolVar(&args.NoProgress, "noprogress", false, helpMsgNoProgress)
	flag.Parse()
	args.Args = flag.Args()

	if err := args.validate(); err != nil {
		return nil, err
	}

	return &args, nil
}

func flagHelpMessage() {
	cmd := os.Args[0]
	out := os.Stderr
	fmt.Fprintln(out, fmt.Sprintf("%s convert text to '%s' style.", cmd, appName))
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out, fmt.Sprintf("  %s [OPTIONS] [files...]", cmd))
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Examples:")
	fmt.Fprintln(out, fmt.Sprintf("  %s sample.txt", cmd))
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Options:")

	flag.PrintDefaults()
}

func (c *CmdArgs) validate() error {
	switch c.Style {
	case "free_text", "json", "json_stream":
		// nothing to do
	default:
		return errors.New("style must be 'free_text', 'json' or 'json_stream'")
	}

	if c.Version {
		return nil
	}

	if len(c.Args) < 2 {
		return errors.New("args must need 2 files or more. CIDR list file and IP addresses list files")
	}

	return nil
}
