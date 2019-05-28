// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var (
	workDir   string
	outputDir string
)

func genFromSql(g *Generator, brickName string, sourceFilename string,
	outputFilename string, statements []Statement, syntaxes []Syntax) {
	g.header(sourceFilename)
	g.GenerateBrick(sourceFilename, brickName, syntaxes, statements)
	for _, value := range statements {
		g.Generate(brickName, value)
	}

	if err := g.Output(outputFilename); err != nil {
		log.Fatalf("error: generator file: %s", err)
	}
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Generate code fail: %v", err)
		return
	}

	flag.StringVar(&workDir, "w", dir, "The work directory to search sql files")
	flag.StringVar(&outputDir, "o", "", "The output directory of generated source code")
	flag.Parse()

	if len(outputDir) == 0 {
		outputDir = workDir
	}

	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		log.Fatalf("error: check directory: %s", err)
	}

	if files, err := getSqlFiles(workDir); err != nil {
		log.Fatalf("error: generate from file: %s", err)
	} else if len(files) == 0 {
		log.Printf("error: no sqb files found in current dir")
		flag.Usage()
	} else {
		var bricks []string
		var txMap = make(map[string]bool)
		for _, value := range files {
			b := getBrickName(value)
			bricks = append(bricks, b)
			p := NewParser()
			statements, syntaxes, err := p.LoadSqlFile(value)
			if err != nil {
				log.Fatalf("parse sql file fail: %s", err)
				break
			}

			inputFileName := getFileName(value)
			outputFileName := getSourceName(value) + ".go"

			entityGen := newGenerator(filepath.Join(outputDir, "entity"), "entity")
			entityGen.GenerateEntity(inputFileName, outputFileName, b, syntaxes)

			g := newGenerator(filepath.Join(outputDir, "brick"), "brick")
			hasTx := g.CheckTx(statements)
			if hasTx {
				txMap[b] = hasTx
			}

			genFromSql(g, b, inputFileName, outputFileName, statements, syntaxes)
		}

		bricksGen := newGenerator(filepath.Join(outputDir, "brick"), "brick")
		bricksGen.header("")
		bricksGen.GenerateSqlBrick(bricks, txMap)
		if err := bricksGen.Output("sqlbrick.go"); err != nil {
			log.Fatalf("error: writing output: %s", err)
		} else {
			log.Printf("SqlBrick has wrote to: %s", outputDir)
		}
	}
}
