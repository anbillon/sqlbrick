// Copyright (c) 2019 Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

const version = "0.2.0"

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the SQLBrick version info",
		Run:   runVersion,
	}
}

func runVersion(_ *cobra.Command, _ []string) {
	fmt.Printf(`SQLBrick:
 Version:       %v
 Go version:    %v`,
		version, runtime.Version())
}
