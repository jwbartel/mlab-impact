//Package defines how data sources are accessed by impact
package data

import ( // for docs http://golang.org/pkg/ pkgname
	"impact/data/census/acs"
	"impact/data/census/sf1"
	"impact/data/ndt"
	"net/http"
	"reflect"
)

/*
 * Interface each source must implement
 * 
 * Results of interfaces are maps of strings to values
 * 
 * The keys identify the title of keys and values should be primitives
 * or maps determing substructures, all of which can be interpreted as
 * JSON objects by json.Marshal
 */
type Source interface {
	Query(r *http.Request, clientLoc map[string]interface{}, serverLoc map[string]interface{}) (map[string]interface{}, error)
}

//Sources used for querying
var sources = []Source{acs.ACS_Source(), sf1.SF1_Source(), ndt.NDT_Source()} //TestSource{}} 

//Required fields before a query is determined to be complete
var requiredFields = []string{"Throughput", "RTT", "Packet Loss"}

func Query(r *http.Request, clientLoc map[string]interface{}, serverLoc map[string]interface{}) (map[string]interface{}, error) {
	result := DefaultResult()
	for i := 0; i < len(sources); i++ {
		newResult, err := sources[i].Query(r, clientLoc, serverLoc)

		if err != nil {
			result["err"] = err.Error()
			return result, nil
		}

		result = merge(result, newResult)

		if resultIsComplete(result) {
			return result, nil
		}
	}
	result["client"] = clientLoc
	result["server"] = serverLoc
	return result, nil
}

func mergeValues(oldValue interface{}, newValue interface{}) interface{} {
	if reflect.TypeOf(newValue).String() == "map[string]interface {}" {
		oldMap, oldOk := oldValue.(map[string]interface{})
		newMap, newOk := newValue.(map[string]interface{})
		if oldOk && newOk {
			merge(oldMap, newMap)
		}
	}
	return oldValue
}

func merge(oldResult map[string]interface{}, newResult map[string]interface{}) map[string]interface{} {

	for key, value := range newResult {
		if oldResult[key] == nil {
			oldResult[key] = value
		} else {
			oldResult[key] = mergeValues(oldResult[key], value)
		}
	}

	return oldResult
}

func resultIsComplete(r map[string]interface{}) bool {

	for _, val := range requiredFields {
		if r[val] == nil {
			return false
		}
	}

	return true
}

func DefaultResult() map[string]interface{} {
	return make(map[string]interface{})
}
