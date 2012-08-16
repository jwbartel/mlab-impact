package ndt

import (
	"code.google.com/p/google-api-go-client/bigquery/v2"
	"fmt"
	"net/http"
	"time"
)

type NDT struct {
	ProjectID string
	DatasetID string
}

func NDT_Source() *NDT {

	n := &NDT{
		ProjectID: "measurement-lab",
		DatasetID: "m_lab",
	}
	return n
}

var (
	DefaultDatasetID = "m_lab"
	DefaultFields    = []string{
		"MSSSent",
		"MSSRcvd",
		"SegsOut",
		"DataSegsOut",
		"DataOctetsOut",
		"SegsIn",
		"DataSegsIn",
		"DataOctetsIn",
		"Duration",
		"CurMSS",
		"SampleRTT",
		"SmoothedRTT",
		"RTTVar",
		"MaxRTT",
		"MinRTT",
		"SumRTT",
		"CountRTT",
		"CurRTO",
		"MaxRTO",
		"MinRTO",
		"SlowStart",
		"CongAvoid",
		"CongSignals",
		"OtherReductions",
		"CongOverCount",
		"CurCwnd",
		"MaxSsCwnd",
		"MaxCaCwnd",
		"MaxSsthresh",
		"MinSsthresh",
		"LimSsthresh",
		"Timeouts",
		"SubsequentTimeouts",
		"RcvRTT",
		"CurRwinRcvd",
		"MaxRwinRcvd",
	}
)

func (ndt *NDT) getQueryWhere(clientLoc map[string]interface{}, serverLoc map[string]interface{}) string {
	whereStr := ""

	if clientLoc["Country"] != nil {
		whereStr = fmt.Sprintf("connection_spec.client_geolocation.country_name=\"%v\"", clientLoc["Country"])
		if clientLoc["State/Region"] != nil {
			whereStr = fmt.Sprintf("%v AND connection_spec.client_geolocation.region=\"%v\"", whereStr, clientLoc["State/Region"])
		}
		if clientLoc["City"] != nil {
			whereStr = fmt.Sprintf("%v AND connection_spec.client_geolocation.city=\"%v\"", whereStr, clientLoc["City"])
		}
	}

	if len(whereStr) > 0 {
		whereStr = fmt.Sprintf("WHERE %v", whereStr)
	}
	return whereStr
}

func (ndt *NDT) getQueryFields(fields []string) string {

	fieldPart := "SELECT count(*)"
	for _, field := range fields {
		field = fmt.Sprintf("web100_log_entry.snap.%v", field)
		fieldPart = fmt.Sprintf("%v, AVG(%v), STDDEV(%v)", fieldPart, field, field)
	}
	return fieldPart
}

func (ndt *NDT) getQueryTable(year int, month int) string {

	monthStr := fmt.Sprintf("%v", month)
	for i := 0; i < 2-len(monthStr); i++ {
		monthStr = fmt.Sprintf("0%v", monthStr)
	}
	tablePart := fmt.Sprintf("FROM %v.%v_%v", ndt.DatasetID, year, monthStr)
	return tablePart
}

func (ndt *NDT) parseRows(fields []string, rows []*bigquery.TableRow, result map[string]interface{}) {

	if len(rows) > 0 {
		row := rows[0].F
		result["sample size"] = row[0].V
		for pos, fieldName := range fields {
			stats := make(map[string]interface{})
			stats["average"] = row[2*pos+1].V
			stats["stdev"] = row[2*pos+2].V
			result[fieldName] = stats
		}
	}

}

func (ndt *NDT) GetData(r *http.Request, fields []string, year int, month int, clientLoc map[string]interface{}, serverLoc map[string]interface{}) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	fieldPart := ndt.getQueryFields(fields)
	tablePart := ndt.getQueryTable(year, month)
	wherePart := ndt.getQueryWhere(clientLoc, serverLoc)

	query := fmt.Sprintf("%v %v %v", fieldPart, tablePart, wherePart)

	queryResponse, err := ndt.askBigQuery(r, query)

	if err != nil {
		return nil, err
	}
	dataResult := make(map[string]interface{})
	dataResult["complete"] = queryResponse.JobComplete
	if queryResponse.JobComplete && queryResponse.TotalRows > 0 {
		ndt.parseRows(fields, queryResponse.Rows, dataResult)
	} else {
		dataResult["jobID"] = queryResponse.JobReference.JobId
	}
	result["network data"] = dataResult
	return result, nil
}

func (ndt *NDT) askBigQuery(r *http.Request, query string) (*bigquery.QueryResponse, error) {

	client := getJWTClient(r)
	bigqueryService, err := bigquery.New(client)
	if err != nil {
		return nil, err
	}

	queryRequest := &bigquery.QueryRequest{
		Query: query,
	}

	queryCall := bigqueryService.Jobs.Query(ndt.ProjectID, queryRequest)
	return queryCall.Do()

}

func (ndt *NDT) JobResult(r *http.Request, jobID string) (map[string]interface{}, error) {

	client := getJWTClient(r)
	bigqueryService, err := bigquery.New(client)
	if err != nil {
		return nil, err
	}

	jobsService := bigqueryService.Jobs.GetQueryResults(ndt.ProjectID, jobID)
	jobsService.TimeoutMs(5000)
	response, err := jobsService.Do()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["complete"] = response.JobComplete
	if response.JobComplete {
		ndt.parseRows(DefaultFields, response.Rows, result)
	} else {
		result["jobID"] = response.JobReference.JobId
	}
	return result, nil

}

func (ndt *NDT) Query(r *http.Request, clientLoc map[string]interface{}, serverLoc map[string]interface{}) (map[string]interface{}, error) {

	fields := DefaultFields
	year, month, _ := time.Now().Date()
	return ndt.GetData(r, fields, year, int(month), clientLoc, serverLoc)
}
