package census

import (
	"appengine"
	"appengine/urlfetch"

	"bytes"
	"errors"
	"fmt"
	"impact/data/secrets"
	"io/ioutil"
	"net/http"
	"strings"
)

type CensusRequester struct {
	key string
}

func DefaultCensusRequester() *CensusRequester {
	c := &CensusRequester{
		key: secrets.Keys().CensusKey,
	}
	return c
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (cr CensusRequester) FetchResults(r *http.Request, url string) (string, error) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Get(url)

	if err != nil {
		return "", errors.New("response")
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("body")
		return "", err
	}

	return bytes.NewBuffer(body).String(), nil
}

func (cr CensusRequester) getFieldNames(fields [][]string, fieldCodes []string) []string {

	fieldNames := make([]string, len(fieldCodes))

	for index, code := range fieldCodes {
		for _, entry := range fields {
			if entry[0] == code {
				fieldNames[index] = entry[1]
			}
		}
	}
	return fieldNames
}

func (cr *CensusRequester) buildSpace(spaceName string, parent map[string]interface{}) map[string]interface{} {
	index := strings.Index(spaceName, " - ")
	if index == -1 {
		space := make(map[string]interface{})
		parent[spaceName] = space
		return space
	}
	currSpace := spaceName[:index]
	if parent[currSpace] == nil {
		parent[currSpace] = make(map[string]interface{})
	}
	return cr.buildSpace(spaceName[index+3:], parent[currSpace].(map[string]interface{}))
}

func (cr *CensusRequester) ParseResults(fields [][]string, strResults string, result map[string]interface{}) {
	lines := strings.Split(strResults, "\n")

	fieldCodes := strings.Split(lines[0][3:], "\",\"")
	fieldNames := cr.getFieldNames(fields, fieldCodes)

	fieldVals := strings.Split(lines[1][2:], "\",\"")
	for index, name := range fieldNames {
		if name != "" {
			space := cr.buildSpace(name, result)
			space["total"] = fieldVals[index]
		}
	}
}

func (cr *CensusRequester) FormURL(r *http.Request, urlPrefix string, fields [][]string, countyCode string, stateCode string) string {

	url := fmt.Sprintf("%vkey=%v", urlPrefix, cr.key)
	if len(fields) > 0 {
		fieldsStr := ""
		for _, entry := range fields {
			fieldsStr = fmt.Sprintf("%v,%v", fieldsStr, entry[0])
		}
		fieldsStr = fieldsStr[1:]
		url = fmt.Sprintf("%v&get=%v", url, fieldsStr)
	}
	if stateCode != "" {
		if countyCode == "" {
			url = fmt.Sprintf("%v&for=state:%v", url, stateCode)
		} else {
			url = fmt.Sprintf("%v&in=state:%v", url, stateCode)
		}
	}
	if countyCode != "" {
		url = fmt.Sprintf("%v&for=county:%v", url, countyCode)
	}

	return url
}

func (cr *CensusRequester) AskApiInChunks(r *http.Request, urlPrefix string, fields [][]string, countyCode string, stateCode string, maxFields int) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for fieldStart := 0; fieldStart < len(fields); {

		fieldEnd := min(fieldStart+maxFields, len(fields))
		subFields := fields[fieldStart:fieldEnd]

		url := cr.FormURL(r, urlPrefix, subFields, countyCode, stateCode)

		censusResult, err := cr.FetchResults(r, url)
		if err != nil {
			return nil, err
		}
		cr.ParseResults(fields, censusResult, result)
		fieldStart = fieldEnd
	}
	return result, nil

}
