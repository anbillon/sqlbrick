// Copyright (c) 2019 Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/anbillon/sqlbrick/internal"
	"github.com/spf13/cobra"
)

var (
	workDir     string
	outputDir   string
	withContext bool
)

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate SQLBrick files from input sqb files",
		RunE:  runGeneration,
	}
	initFlags(cmd)

	return cmd
}

func initFlags(cmd *cobra.Command) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("%v", err)
	}

	cmd.Flags().StringVarP(&workDir, "work-dir", "w", dir,
		fmt.Sprintf(`The work directory to search sql files (default "%v")`, dir))
	cmd.Flags().StringVarP(&outputDir, "output-dir", "o", "",
		"The output directory of generated source code")
	cmd.Flags().BoolVarP(&withContext, "with-context", "c", false,
		"Should SQLBrick support context or not (default false)")
}

func runGeneration(_ *cobra.Command, _ []string) error {
	if len(outputDir) == 0 {
		outputDir = workDir
	}

	var err error
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return err
	}

	if files, err := internal.GetSqlFiles(workDir); err != nil {
		return err
	} else if len(files) == 0 {
		return errors.New("error: no sqb files found in current dir")
	} else {
		var bricks []string
		var txMap = make(map[string]bool)
		for _, value := range files {
			b := internal.GetBrickName(value)
			bricks = append(bricks, b)
			p := internal.NewParser()
			statements, syntaxes, err := p.LoadSqlFile(value)
			if err != nil {
				log.Fatalf("parse sql file fail: %s", err)
				break
			}

			inputFileName := internal.GetFileName(value)
			outputFileName := internal.GetSourceName(value) + ".go"

			entityGen := internal.NewGenerator(filepath.Join(outputDir, "entity"), "entity")
			entityGen.GenerateEntity(inputFileName, outputFileName, b, syntaxes)

			g := internal.NewGenerator(filepath.Join(outputDir, "brick"), "brick")
			hasTx := g.CheckTx(statements)
			if hasTx {
				txMap[b] = hasTx
			}

			genFromSql(g, b, inputFileName, outputFileName, statements, syntaxes)
		}

		bricksGen := internal.NewGenerator(filepath.Join(outputDir, "brick"), "brick")
		bricksGen.Header("")
		bricksGen.GenerateSqlBrick(bricks, txMap)
		if err := bricksGen.Output("sqlbrick.go"); err != nil {
			log.Fatalf("error: writing output: %s", err)
		} else {
			fmt.Printf("SQLBrick has wrote to: %s", outputDir)
		}
	}

	return nil
}

func genFromSql(g *internal.Generator, brickName string, sourceFilename string,
	outputFilename string, statements []internal.Statement, syntaxes []internal.Syntax) {
	g.Header(sourceFilename)
	g.GenerateBrick(sourceFilename, brickName, syntaxes, statements)
	for _, value := range statements {
		g.GenerateSqlFunc(brickName, withContext, value)
	}

	if err := g.Output(outputFilename); err != nil {
		log.Fatalf("error: generator file: %s", err)
	}
}
