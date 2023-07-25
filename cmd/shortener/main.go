package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuildInfo("Build version", buildVersion)
	printBuildInfo("Build date", buildDate)
	printBuildInfo("Build commit", buildCommit)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go runService(ctx)
	<-ctx.Done()
}

func printBuildInfo(label, value string) {
	if len(value) == 0 {
		value = "N/A"
	}
	fmt.Printf("%s: %s\n", label, value)
}
