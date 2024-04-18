package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
)

const (
	appName = "checkcidr"
)

// CIでビルド時に値を埋め込む。
// 埋め込む値の設定は .goreleaser.yaml を参照。
var (
	version  = "dev"
	revision = "dev"
	logger   = slog.New(slog.NewTextHandler(os.Stderr, nil))
)

const (
	exitStatusOK = iota
	exitStatusCLIError
	exitStatusConvertError
	exitStatusInputFileError
	exitStatusOutputError
)

func main() {
	args, err := ParseArgs()
	if err != nil {
		logger.Error("failed to parse args", "err", err)
		os.Exit(exitStatusCLIError)
	}

	if args.Version {
		msg := fmt.Sprintf("%s %s (%s)", appName, version, revision)
		fmt.Println(msg)
		fmt.Println("")
		fmt.Println("author:     jiro")
		fmt.Println("repository: https://github.com/jiro4989/checkcidr")
		os.Exit(exitStatusOK)
	}

	cidr, err := ReadCIDR(args.Args[0])
	if err != nil {
		logger.Error("failed to read cidr file", "err", err)
		os.Exit(exitStatusInputFileError)
	}

	output := os.Stdout
	if args.OutFile != "" {
		output, err = os.Create(args.OutFile)
		if err != nil {
			logger.Error("failed to create a output file", "err", err)
			os.Exit(exitStatusInputFileError)
		}
		defer output.Close()
	}

	results := make([]result, 0)
	p := progress{
		noProgress: args.NoProgress,
		counter:    0,
	}
	for _, file := range args.Args[1:] {
		cp, err := os.Open(file)
		defer cp.Close()
		if err != nil {
			logger.Error("failed to read ip file", "err", err)
			os.Exit(exitStatusInputFileError)
		}

		sc := bufio.NewScanner(cp)
		for sc.Scan() {
			l := strings.TrimSpace(sc.Text())
			ip := net.ParseIP(l)
			for _, c := range cidr {
				p.increment()
				p.prints()

				contains := c.Contains(ip)
				r := result{
					IPFile:   file,
					CIDR:     c,
					IP:       ip,
					Contains: contains,
					Style:    args.Style,
				}
				if isLinePrinting(args.Style) {
					text, err := r.format()
					if err != nil {
						logger.Error("failed to marshal json", "err", err)
						os.Exit(exitStatusInputFileError)
					}
					fmt.Fprintln(output, text)
				} else {
					results = append(results, r)
				}
			}
		}
		if err := sc.Err(); err != nil {
			logger.Error("failed to read ip file", "err", err)
			os.Exit(exitStatusInputFileError)
		}
	}

	if !isLinePrinting(args.Style) {
		b, err := json.Marshal(results)
		if err != nil {
			logger.Error("failed to marshal json", "err", err)
			os.Exit(exitStatusInputFileError)
		}
		s := string(b)
		fmt.Fprintln(output, s)
	}

	os.Exit(exitStatusOK)
}

type result struct {
	IPFile   string     `json:"ip_file"`
	CIDR     *net.IPNet `json:"cidr"`
	IP       net.IP     `json:"ip"`
	Contains bool       `json:"contains"`
	Style    string     `json:"-"`
}

func (r *result) format() (string, error) {
	switch r.Style {
	case "free_text":
		return fmt.Sprintf("ip_file=%v cidr=%v ip=%v contains=%v", r.IPFile, r.CIDR, r.IP, r.Contains), nil
	}

	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type progress struct {
	noProgress bool
	counter    int
}

func (p *progress) increment() {
	p.counter++
}

func (p *progress) prints() {
	if p.noProgress {
		return
	}

	output := os.Stderr
	if p.counter%2500 == 0 {
		fmt.Fprintf(output, ".")
	}
	if p.counter%100000 == 0 {
		text := fmt.Sprintf(" %d", p.counter)
		fmt.Fprintln(output, text)
	}
}

func isLinePrinting(style string) bool {
	switch style {
	case "free_text", "json_stream":
		return true
	}
	return false
}
