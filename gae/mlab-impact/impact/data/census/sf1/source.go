package sf1

import (
	"impact/data/census"
	"net/http"
)

var (
	DefaultFields = [][]string{
		[]string{"P0030001", "Total population"},
		[]string{"P0030002", "White alone"},
		[]string{"P0030003", "Black or African American alone"},
		[]string{"P0030004", "American Indian and Alaska Native alone"},
		[]string{"P0030005", "Asian alone"},
		[]string{"P0030006", "Native Hawaiian and Other Pacific Islander alone"},
		[]string{"P0030007", "Some Other Race alone"},
		[]string{"P0030008", "Two or More Races"},

		[]string{"P0120001", "Total population"},
		[]string{"P0120002", "Male:"},
		[]string{"P0120003", "Male: - Under 5 years"},
		[]string{"P0120004", "Male: - 5 to 9 years"},
		[]string{"P0120005", "Male: - 10 to 14 years"},
		[]string{"P0120006", "Male: - 15 to 17 years"},
		[]string{"P0120007", "Male: - 18 and 19 years"},
		[]string{"P0120008", "Male: - 20 years"},
		[]string{"P0120009", "Male: - 21 years"},
		[]string{"P0120010", "Male: - 22 to 24 years"},
		[]string{"P0120011", "Male: - 25 to 29 years"},
		[]string{"P0120012", "Male: - 30 to 34 years"},
		[]string{"P0120013", "Male: - 35 to 39 years"},
		[]string{"P0120014", "Male: - 40 to 44 years"},
		[]string{"P0120015", "Male: - 45 to 49 years"},
		[]string{"P0120016", "Male: - 50 to 54 years"},
		[]string{"P0120017", "Male: - 55 to 59 years"},
		[]string{"P0120018", "Male: - 60 and 61 years"},
		[]string{"P0120019", "Male: - 62 to 64 years"},
		[]string{"P0120020", "Male: - 65 and 66 years"},
		[]string{"P0120021", "Male: - 67 to 69 years"},
		[]string{"P0120022", "Male: - 70 to 74 years"},
		[]string{"P0120023", "Male: - 75 to 79 years"},
		[]string{"P0120024", "Male: - 80 to 84 years"},
		[]string{"P0120025", "Male: - 85 years and over"},
		[]string{"P0120026", "Female:"},
		[]string{"P0120027", "Female: - Under 5 years"},
		[]string{"P0120028", "Female: - 5 to 9 years"},
		[]string{"P0120029", "Female: - 10 to 14 years"},
		[]string{"P0120030", "Female: - 15 to 17 years"},
		[]string{"P0120031", "Female: - 18 and 19 years"},
		[]string{"P0120032", "Female: - 20 years"},
		[]string{"P0120033", "Female: - 21 years"},
		[]string{"P0120034", "Female: - 22 to 24 years"},
		[]string{"P0120035", "Female: - 25 to 29 years"},
		[]string{"P0120036", "Female: - 30 to 34 years"},
		[]string{"P0120037", "Female: - 35 to 39 years"},
		[]string{"P0120038", "Female: - 40 to 44 years"},
		[]string{"P0120039", "Female: - 45 to 49 years"},
		[]string{"P0120040", "Female: - 50 to 54 years"},
		[]string{"P0120041", "Female: - 55 to 59 years"},
		[]string{"P0120042", "Female: - 60 and 61 years"},
		[]string{"P0120043", "Female: - 62 to 64 years"},
		[]string{"P0120044", "Female: - 65 and 66 years"},
		[]string{"P0120045", "Female: - 67 to 69 years"},
		[]string{"P0120046", "Female: - 70 to 74 years"},
		[]string{"P0120047", "Female: - 75 to 79 years"},
		[]string{"P0120048", "Female: - 80 to 84 years"},
		[]string{"P0120049", "Female: - 85 years and over"},

		[]string{"P0180001", "Households"},
		[]string{"P0180002", "Family households:"},
		[]string{"P0180003", "Family households: - Husband-wife family"},
		[]string{"P0180004", "Family households: - Other family:"},
		[]string{"P0180005", "Family households: - Other family: - Male householder, no wife present"},
		[]string{"P0180006", "Family households: - Other family: - Female householder, no husband present"},
		[]string{"P0180007", "Nonfamily households: - Other family: - Female householder, no husband present"},
		[]string{"P0180008", "Nonfamily households: - Householder living alone - Female householder, no husband present"},
		[]string{"P0180009", "Nonfamily households: - Householder not living alone - Female householder, no husband present"},
	}
)

type SF1 struct {
}

func SF1_Source() *SF1 {
	s := &SF1{}
	return s
}

func (sf1 *SF1) getSF1Results(r *http.Request, fields [][]string, loc map[string]interface{}) (map[string]interface{}, error) {

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
	return requester.AskApiInChunks(r, "http://api.census.gov/data/2010/sf1?", DefaultFields, county, state, 5)
}

func (sf1 *SF1) Query(r *http.Request, clientLoc map[string]interface{}, serverLoc map[string]interface{}) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	sf1_val := make(map[string]interface{})
	result["SF1"] = sf1_val

	client_sf1, err := sf1.getSF1Results(r, DefaultFields, clientLoc)
	if err != nil {
		return nil, err
	}
	sf1_val["client"] = client_sf1

	server_sf1, err := sf1.getSF1Results(r, DefaultFields, serverLoc)
	if err != nil {
		return nil, err
	}
	sf1_val["server"] = server_sf1

	return result, nil
}
