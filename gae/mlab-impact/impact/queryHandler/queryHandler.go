/* The ns package is intended provides the querying functionality for M-Lab Impact */
package queryHandler

import ( // for docs http://golang.org/pkg/ pkgname
	"impact/data"
	"net/http"
	"strconv"
)

type ComparativeValues struct {
	Local float32
	State float32
	World float32
}

type SearchResult struct {
	Latitude, Longitude float64
	City                string
	State               string
	Country             string
	Throughput          ComparativeValues
	RTT                 ComparativeValues
	Packet_loss         ComparativeValues
}

func GetResult(r *http.Request) (map[string]interface{}, error) {

	clientGeolocation, err := location(r, "c")
	if err != nil {
		result := make(map[string]interface{})
		result["error"] = err.Error()
		return result, nil
	}
	serverGeolocation, err := location(r, "s")
	if err != nil {
		result := make(map[string]interface{})
		result["error"] = err.Error()
		return result, nil
	}
	return data.Query(r, clientGeolocation, serverGeolocation)
}

func location(r *http.Request, prefix string) (map[string]interface{}, error) {
	geo := data.DefaultGeolocator()
	geolocation := make(map[string]interface{})

	qType := r.FormValue(prefix + "Type")
	if qType == "latlng" {
		lat, err := strconv.ParseFloat(r.FormValue(prefix+"Lat"), 64)
		if err != nil {
			geolocation["Error"] = "Unparseable client latitude"
			return geolocation, nil
		}
		long, err := strconv.ParseFloat(r.FormValue(prefix+"Long"), 64)
		if err != nil {
			geolocation["Error"] = "Unparseable client longitude"
			return geolocation, nil
		}
		return geo.ByLatLong(r, lat, long)
	} else if qType == "cityregioncountry" {

		result := make(map[string]interface{})

		result["Country"] = r.FormValue(prefix + "Country")
		if result["Country"] == "" {
			delete(result, "Country")
		}
		result["Region"] = r.FormValue(prefix + "Region")
		if result["Region"] == "" {
			delete(result, "Region")
		}
		result["County"] = r.FormValue(prefix + "County")
		if result["County"] == "" {
			delete(result, "County")
		}
		result["City"] = r.FormValue(prefix + "City")
		if result["City"] == "" {
			delete(result, "City")
		}

		return result, nil
	}
	return make(map[string]interface{}), nil
}
