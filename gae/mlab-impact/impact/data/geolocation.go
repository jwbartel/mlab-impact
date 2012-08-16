package data

import (
	"appengine"
	"appengine/urlfetch"

	"encoding/json"
	"errors"
	"fmt"
	"impact/data/secrets"
	"io/ioutil"
	"net/http"
	"strings"
)

type Geolocator struct {
	apiKey string
	url    string
	sensor bool
}

func DefaultGeolocator() *Geolocator {
	return &Geolocator{
		apiKey: secrets.Keys().GoogleAPIKey,
		url:    "https://maps.googleapis.com/maps/api/geocode/json?",
		sensor: true,
	}
}

func (geo *Geolocator) ByLatLong(r *http.Request, latitude float64, longitude float64) (map[string]interface{}, error) {

	parameters := fmt.Sprintf("latlng=%f,%f&sensor=%t",
		latitude, longitude, geo.sensor)

	result := make(map[string]interface{})
	result["lat"] = latitude
	result["lng"] = longitude
	err := geo.askGoogle(r, parameters, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (geo *Geolocator) ByCityRegionCountry(r *http.Request, city string, region string, country string) (map[string]interface{}, error) {

	parameters := "address="
	if city != "" {
		parameters += strings.Replace(city, " ", "+", -1)
	}
	parameters += ",+"
	if region != "" {
		parameters += strings.Replace(region, " ", "+", -1)
	}
	parameters += "+"
	if country != "" {
		parameters += strings.Replace(country, " ", "+", -1)
	}
	parameters += fmt.Sprintf("&sensor=%t", geo.sensor)

	result := make(map[string]interface{})
	if city != "" {
		result["City"] = city
	}
	if region != "" {
		result["State/Region"] = region
	}
	if country != "" {
		result["Country"] = country
	}

	err := geo.askGoogle(r, parameters, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (geo *Geolocator) askGoogle(r *http.Request, parameters string, result map[string]interface{}) error {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	resp, err := client.Get(geo.url + parameters)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jsonVal := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonVal)
	if err != nil {
		return err
	}

	switch jsonVal["status"] {
	case "ZERO_RESULTS":
		return errors.New("No geolocation results found")
	case "OVER_QUERY_LIMIT":
		return errors.New("Geolocation query limit reached") //TODO: stop making geolocation querys for the remainder of the day
	case "REQUEST_DENIED":
		return errors.New("Geolocation request denied")
	case "INVALID_REQUEST":
		return errors.New("Geolocation request invalid")
	default:
		break
	}

	possibleLocations := jsonVal["results"]
	geo.parseResults(possibleLocations.([]interface{}), result)
	return nil
	result = make(map[string]interface{})
	result["location"] = (jsonVal["results"].([]interface{}))[0]
	return nil
}

func (geo *Geolocator) parseResults(jsonVal []interface{}, result map[string]interface{}) {

	for _, location := range jsonVal {
		locationMap := location.(map[string]interface{})

		address := locationMap["address_components"]
		if address != nil {
			addressArray := address.([]interface{})
			geo.parseAddress(addressArray, result)
		}

		geometry := locationMap["geometry"]
		if geometry != nil {
			geometryMap := geometry.(map[string]interface{})
			geo.parseGeometry(geometryMap, result)
		}
	}
}

func (geo *Geolocator) parseGeometry(jsonVal map[string]interface{}, result map[string]interface{}) {
	location := jsonVal["location"]
	if location != nil {
		locationMap := location.(map[string]interface{})
		if result["lat"] == nil && result["lng"] == nil &&
			locationMap["lat"] != nil && locationMap["lng"] != nil {

			result["lat"] = locationMap["lat"]
			result["lng"] = locationMap["lng"]
		}
	}
}

func (geo *Geolocator) parseAddress(jsonVal []interface{}, result map[string]interface{}) {
	for _, component := range jsonVal {
		componentMap := component.(map[string]interface{})

		label := ""
		typeStr := geo.getType(componentMap)
		name := geo.getLongName(componentMap)

		switch typeStr {
		case "country":
			label = "Country"
		case "administrative_area_level_1":
			label = "State/Region"
		case "administrative_area_level_2":
			label = "County"
		case "locality":
			label = "City"
		case "postal_code":
			label = "Zip"
		default:
			label = ""
		}
		if label != "" && result[label] == nil && name != nil {
			result[label] = name
		}
	}

}

func (geo *Geolocator) getType(componentMap map[string]interface{}) string {
	types := componentMap["types"].([]interface{})
	if types == nil {
		return ""
	}
	if types[0] == nil {
		return ""
	}
	return types[0].(string)
}

func (geo *Geolocator) getLongName(componentMap map[string]interface{}) interface{} {
	return componentMap["long_name"]
}
