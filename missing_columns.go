package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/smallnest/gen/dbmeta"
)

func generateMissingColumns(conf *dbmeta.Config) error {
	var err error

	modelDir := filepath.Join(*outDir, *modelPackageName)
	var JsonTmpl *dbmeta.GenTemplate

	if JsonTmpl, err = LoadTemplate("colInfo.json.tmpl"); err != nil {
		fmt.Print(au.Red(fmt.Sprintf("Error loading template %v\n", err)))
		return err
	}

	var tableNames []string
	for tableName, _ := range tableInfos {
		tableNames = append(tableNames, tableName)
	}

	var data = make(map[string]interface{})
	data["tableNames"] = tableNames
	data["colNames"] = dbmeta.MissingColumns

	jsonFile := filepath.Join(modelDir, "mapping.json")
	err = conf.InstantiateTemplate(JsonTmpl, data, jsonFile)
	if err != nil {
		fmt.Print(au.Red(fmt.Sprintf("Error writing file: %v\n", err)))
		os.Exit(1)
	}

	return nil
}

func loadMissingColumns() error {
	modelDir := filepath.Join(*outDir, *modelPackageName)
	file := filepath.Join(modelDir, "mapping.json")

	if jsonFile, err := os.Open(file); err != nil {
		return err
	} else {
		defer jsonFile.Close()
		if byteValue, err := ioutil.ReadAll(jsonFile); err != nil {
			return err
		} else {
			if err := json.Unmarshal(byteValue, &dbmeta.ColInfos); err != nil {
				return err
			}
		}
	}

	return nil
}
