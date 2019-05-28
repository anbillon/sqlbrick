// Copyright (c) 2019 Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sqlbrick",
		Short: "SQLBrick generates golang function from your SQL statements",
	}
	cmd.AddCommand(newVersionCmd(), newGenerateCmd())

	return cmd
}
