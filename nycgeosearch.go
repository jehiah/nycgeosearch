// package nycgeosearch for calling geosearch.planninglabs.nyc API
package nycgeosearch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	geojson "github.com/paulmach/go.geojson"
)

// Client implements the pelias format against the NYC Planning labs GeoSearch endpoint
//
// https://github.com/pelias/documentation
// https://geosearch.planninglabs.nyc/docs/
type Client string

var PlanningLabs = Client("https://geosearch.planninglabs.nyc")

type Location struct {
	Lat, Lng float64
}

type Options struct {
	Size int // defaults to 10
}

func (o Options) getSize() string {
	if o.Size < 1 {
		return "10"
	}
	return strconv.Itoa(o.Size)
}

// https://github.com/pelias/documentation/blob/master/search.md
func (c Client) Search(ctx context.Context, address string, opt Options) (*geojson.FeatureCollection, error) {
	return c.call(ctx, "/v2/search", &url.Values{
		"size":    []string{opt.getSize()},
		"address": []string{address},
	})
}

// https://github.com/pelias/documentation/blob/master/autocomplete.md
func (c Client) Autocomplete(ctx context.Context, address string, opt Options) (*geojson.FeatureCollection, error) {
	return c.call(ctx, "/v2/autocomplete", &url.Values{
		"size":    []string{opt.getSize()},
		"address": []string{address},
	})
}

// ReverseGeocode looks up addresses from a Lat/Lng
// https://github.com/pelias/documentation/blob/master/reverse.md
func (c Client) ReverseGeocode(ctx context.Context, l Location, opt Options) (*geojson.FeatureCollection, error) {
	return c.call(ctx, "/v2/reverse", &url.Values{
		"size":      []string{opt.getSize()},
		"point.lat": []string{fmt.Sprintf("%f", l.Lat)},
		"point.lon": []string{fmt.Sprintf("%f", l.Lng)},
	})
}

func (c Client) call(ctx context.Context, path string, params *url.Values) (*geojson.FeatureCollection, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", string(c)+path+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected HTTP response %v", resp.StatusCode)
	}
	var data geojson.FeatureCollection
	return &data, json.NewDecoder(resp.Body).Decode(&data)
}
