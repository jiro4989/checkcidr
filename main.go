package main

import (
	"fmt"
	"log/slog"
	"os"
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

	os.Exit(exitStatusOK)
}
