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

func genFromSql(brickName string, sqlFilePath string,
	sourceFilename string, filename string) {
	p := parser.NewParser()
	g := NewGenerator(outputDir, packageName)
	g.header(sourceFilename)
	keys, statements, syntaxes, err := p.LoadSqlFile(sqlFilePath)
	if err != nil {
		log.Fatalf("error: parse sql file: %s", err)
		return
	}

	g.GenerateBrick(sourceFilename, brickName, syntaxes)
	for _, value := range keys {
		g.Generate(brickName, value, statements[value])
	}

	if err := g.Output(filename); err != nil {
		log.Fatalf("error: generator filename: %s", err)
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

	if outputDir == "" {
		outputDir = filepath.Join(workDir, packageName)
	}

	if files, err := getSqlFiles(workDir); err != nil {
		log.Fatalf("error: generate from file: %s", err)
	} else if len(files) == 0 {
		log.Printf("error: no sql files found in current dir")
		flag.Usage()
	} else {
		var bricks []string
		for _, value := range files {
			b := getBrickName(value)
			bricks = append(bricks, b)
			genFromSql(b, value, getFileName(value), getSourceName(value)+".go")
		}

		bricksGen := NewGenerator(outputDir, packageName)
		bricksGen.header("")
		bricksGen.GenerateSqlBrick(bricks)
		if err := bricksGen.Output("sqlbrick.go"); err != nil {
			log.Fatalf("error: writing output: %s", err)
		}
	}
}
