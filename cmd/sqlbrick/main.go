// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"anbillon.com/sqlbrick/cmd/sqlbrick/parser"
)

var (
	workDir     string
	outputDir   string
	packageName string
)

func genFromSql(g *Generator, brickName string, sourceFilename string,
	outputFilename string, statements []parser.Statement, syntaxes []parser.Syntax) {
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
	flag.StringVar(&packageName, "p", "models", "The package name of generated source code")
	flag.Parse()

	if len(outputDir) == 0 {
		outputDir = filepath.Join(workDir, packageName)
	}

	if files, err := getSqlFiles(workDir); err != nil {
		log.Fatalf("error: generate from file: %s", err)
	} else if len(files) == 0 {
		log.Printf("error: no sql files found in current dir")
		flag.Usage()
	} else {
		var bricks []string
		var txMap = make(map[string]bool)
		for _, value := range files {
			b := getBrickName(value)
			bricks = append(bricks, b)
			p := parser.NewParser()
			statements, syntaxes, err := p.LoadSqlFile(value)
			if err != nil {
				log.Fatalf("parse sql file fail: %s", err)
				break
			}
			g := NewGenerator(outputDir, packageName)
			hasTx := g.CheckTx(statements)
			if hasTx {
				txMap[b] = hasTx
			}

			genFromSql(g, b, getFileName(value), getSourceName(value)+".go",
				statements, syntaxes)
		}

		bricksGen := NewGenerator(outputDir, packageName)
		bricksGen.header("")
		bricksGen.GenerateSqlBrick(bricks, txMap)
		if err := bricksGen.Output("sqlbrick.go"); err != nil {
			log.Fatalf("error: writing output: %s", err)
		}
	}
}
