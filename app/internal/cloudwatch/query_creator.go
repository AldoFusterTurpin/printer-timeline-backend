package cloudwatch

import (
	"bytes"
	"html/template"
)

// CreateQuery receives a template string and a maputil of values and
// returns the resulting query of replacing the mapValues in queryTemplateStr.
// It also returns an error, if any.
func CreateQuery(queryTemplateStr string, mapValues map[string]string) (query string, err error) {
	queryTemplate, err := template.New("queryTemplate").Parse(queryTemplateStr)
	if err != nil {
		return
	}

	var queryBuffer bytes.Buffer
	err = queryTemplate.Execute(&queryBuffer, mapValues)
	if err != nil {
		return
	}
	return queryBuffer.String(), nil
}
