// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
)

func getSqlFiles(dir string) ([]string, error) {
	var sqlFiles []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, value := range files {
		if value.IsDir() || !strings.HasSuffix(value.Name(), ".sqb") {
			continue
		}
		sqlFiles = append(sqlFiles, filepath.Join(workDir, value.Name()))
	}

	return sqlFiles, nil
}

func getBrickName(sqlFilePath string) string {
	bn := getSourceName(sqlFilePath)
	return strcase.ToCamel(bn)
}

func getSourceName(sqlFilePath string) string {
	name := getFileName(sqlFilePath)
	dotIndex := strings.LastIndex(name, ".")
	if dotIndex <= 0 {
		return ""
	}
	return name[:dotIndex]
}

func getFileName(sqlFilePath string) string {
	index := strings.LastIndex(sqlFilePath, string(os.PathSeparator))
	if index <= 0 {
		return sqlFilePath
	}
	return sqlFilePath[index+1:]
}
