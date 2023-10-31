package analyze

import (
	"internal/server"

	"reflect"

	"github.com/tom-draper/log-analyzer/pkg/parse"
)

func dataTypeBreakdown(extraction []parse.Extraction) map[string]map[string]int {
	dataTypes := make(map[string]map[string]int)
	for _, e := range extraction {
		for k, p := range e.Params {
			dataType := reflect.TypeOf(p)
			if _, ok := dataTypes[k]; !ok {
				dataTypes[k] = make(map[string]int)
			}
			dataTypes[k][dataType.String()] += 1
		}
	}
	return dataTypes
}

func failedLines(extraction []parse.Extraction) map[int]string {
	failedLines := make(map[int]string)
	for _, e := range extraction {
		if len(e.Params) == 0 {
			failedLines[e.LineNumber] = e.Line
		}
	}
	return failedLines
}

func NewData(extraction []parse.Extraction, config *parse.Config) *server.Data {
	dataTypes := dataTypeBreakdown(extraction)
	failed := failedLines(extraction)
	data := server.Data{Extraction: extraction, DataTypes: dataTypes, Failed: failed, Config: config}
	return &data
}

func Run(extraction []parse.Extraction, config *parse.Config) {
	data := NewData(extraction, config)
	server.Start(data)
}
