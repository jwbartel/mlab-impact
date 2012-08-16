package acs

import (
	"impact/data/census"
	"net/http"
)

var (
	DefaultFields = [][]string{
		[]string{"B19101_002E", "Income - Less than $10,000"},
		[]string{"B19101_003E", "Income - $10,000 to $14,999"},
		[]string{"B19101_004E", "Income - $15,000 to $19,999"},
		[]string{"B19101_005E", "Income - $20,000 to $24,999"},
		[]string{"B19101_006E", "Income - $25,000 to $29,999"},
		[]string{"B19101_007E", "Income - $30,000 to $34,999"},
		[]string{"B19101_008E", "Income - $35,000 to $39,999"},
		[]string{"B19101_009E", "Income - $40,000 to $44,999"},
		[]string{"B19101_010E", "Income - $45,000 to $49,999"},
		[]string{"B19101_011E", "Income - $50,000 to $59,999"},
		[]string{"B19101_012E", "Income - $60,000 to $74,999"},
		[]string{"B19101_013E", "Income - $75,000 to $99,999"},
		[]string{"B19101_014E", "Income - $100,000 to $124,999"},
		[]string{"B19101_015E", "Income - $125,000 to $149,999"},
		[]string{"B19101_016E", "Income - $150,000 to $199,999"},
		[]string{"B19101_017E", "Income - $200,000 or more"},

		[]string{"B05003F_001E", "Sex by Age by Citizenship Status - Total"},
		[]string{"B05003F_002E", "Sex by Age by Citizenship Status - Male"},
		[]string{"B05003F_003E", "Sex by Age by Citizenship Status - Male - Under 18 years"},
		[]string{"B05003F_004E", "Sex by Age by Citizenship Status - Male - Under 18 years - Native"},
		[]string{"B05003F_005E", "Sex by Age by Citizenship Status - Male - Under 18 years - Foreign born"},
		[]string{"B05003F_006E", "Sex by Age by Citizenship Status - Male - Under 18 years - Foreign born - Naturalized U.S. citizen"},
		[]string{"B05003F_007E", "Sex by Age by Citizenship Status - Male - Under 18 years - Foreign born - Not a U.S. citizen"},
		[]string{"B05003F_008E", "Sex by Age by Citizenship Status - Male - 18 years and over"},
		[]string{"B05003F_009E", "Sex by Age by Citizenship Status - Male - 18 years and over - Native"},
		[]string{"B05003F_010E", "Sex by Age by Citizenship Status - Male - 18 years and over - Foreign born"},
		[]string{"B05003F_011E", "Sex by Age by Citizenship Status - Male - 18 years and over - Foreign born - Naturalized U.S. citizen"},
		[]string{"B05003F_012E", "Sex by Age by Citizenship Status - Male - 18 years and over - Foreign born - Not a U.S. citizen"},
		[]string{"B05003F_013E", "Sex by Age by Citizenship Status - Female"},
		[]string{"B05003F_014E", "Sex by Age by Citizenship Status - Female - Under 18 years"},
		[]string{"B05003F_015E", "Sex by Age by Citizenship Status - Female - Under 18 years - Native"},
		[]string{"B05003F_016E", "Sex by Age by Citizenship Status - Female - Under 18 years - Foreign born"},
		[]string{"B05003F_017E", "Sex by Age by Citizenship Status - Female - Under 18 years - Foreign born - Naturalized U.S. citizen"},
		[]string{"B05003F_018E", "Sex by Age by Citizenship Status - Female - Under 18 years - Foreign born - Not a U.S. citizen"},
		[]string{"B05003F_019E", "Sex by Age by Citizenship Status - Female - 18 years and over"},
		[]string{"B05003F_020E", "Sex by Age by Citizenship Status - Female - 18 years and over - Native"},
		[]string{"B05003F_021E", "Sex by Age by Citizenship Status - Female - 18 years and over - Foreign born"},
		[]string{"B05003F_022E", "Sex by Age by Citizenship Status - Female - 18 years and over - Foreign born - Naturalized U.S. citizen"},
		[]string{"B05003F_023E", "Sex by Age by Citizenship Status - Female - 18 years and over - Foreign born - Not a U.S. citizen"},
	}
)

type ACS struct {
}

func ACS_Source() *ACS {
	a := &ACS{}
	return a
}

func (acs *ACS) getACSResults(r *http.Request, fields [][]string, loc map[string]interface{}) (map[string]interface{}, error) {

	if loc["Country"] != "United States" {
		return make(map[string]interface{}), nil
	}

	countyName := loc["County"]
	if countyName == nil {
		countyName = ""
	}
	stateName := loc["Region"]
	if stateName == nil {
		stateName = ""
	}
	county, state := census.CountyAndStateCodes(countyName.(string), stateName.(string))
	requester := census.DefaultCensusRequester()
	return requester.AskApiInChunks(r, "http://api.census.gov/data/2010/acs5?", DefaultFields, county, state, 5)
}

func (acs *ACS) Query(r *http.Request, clientLoc map[string]interface{}, serverLoc map[string]interface{}) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	acs_val := make(map[string]interface{})
	result["ACS"] = acs_val

	client_acs, err := acs.getACSResults(r, DefaultFields, clientLoc)
	if err != nil {
		return nil, err
	}
	acs_val["client"] = client_acs

	server_acs, err := acs.getACSResults(r, DefaultFields, serverLoc)
	if err != nil {
		return nil, err
	}
	acs_val["server"] = server_acs

	return result, nil
}
