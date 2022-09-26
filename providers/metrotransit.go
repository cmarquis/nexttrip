package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
)

// https://svc.metrotransit.org/swagger/index.html
type (
	// Route represents a transit route
	Route struct {
		RouteID    string `json:"route_id"`
		AgencyID   int64  `json:"agency_id"`
		RouteLabel string `json:"route_label"`
	}

	// Direction represents a compass direction
	Direction struct {
		DirectionID   int64  `json:"direction_id"`
		DirectionName string `json:"direction_name"`
	}

	// Stop represents a transit stop.
	Stop struct {
		StopID      int64   `json:"stop_id"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Description string  `json:"description"`
		PlaceCode   string  `json:"place_code"`
	}

	// Alert represents a alert from the transit authority
	Alert struct {
		StopClosed bool   `json:"stop_closed"`
		AlertText  string `json:"alert_text"`
	}

	// Departure represents a departure from a stop
	Departure struct {
		DepartureTime int64 `json:"departure_time"`
	}

	// Departures represents upcoming departurs from a stop
	Departures struct {
		Stops     []Stop      `json:"stops"`
		Alerts    []Alert     `json:"alerts"`
		Depatures []Departure `json:"departures"`
	}
)

// HTTPClient is the client used for contacting metrotransit this can be swapped
// for testing
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// MetroTransitProvider implements Provider for MetroTransitProvider
type MetroTransitProvider struct {
	APIClient  HTTPClient
	UseSandbox bool
}

func (p *MetroTransitProvider) getURL(endpoint string) string {
	baseURL, err := url.Parse("https://svc.metrotransit.org/nextripv2")
	if err != nil {
		panic(err)
	}
	baseURL.Path = path.Join(baseURL.Path, endpoint)
	return baseURL.String()
}

func (p *MetroTransitProvider) request(req *http.Request, out interface{}) error {
	resp, err := p.APIClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, out)
	if err != nil {
		return err
	}
	return nil
}

func (p *MetroTransitProvider) findRoute(routeString string) (Route, error) {
	var route Route
	req, err := http.NewRequest("GET", p.getURL("routes"), nil)
	if err != nil {
		return route, err
	}
	routes := []Route{}
	if err := p.request(req, &routes); err != nil {
		return route, err
	}
	for _, r := range routes {
		if strings.Contains(strings.ToLower(r.RouteLabel), strings.ToLower(routeString)) {
			return r, nil
		}
	}

	return route, fmt.Errorf("no route found that matches %s", routeString)
}

func (p *MetroTransitProvider) findDirection(directionString, routeID string) (Direction, error) {
	var direction Direction
	req, err := http.NewRequest("GET", p.getURL(fmt.Sprintf("directions/%s", routeID)), nil)
	if err != nil {
		return direction, err
	}
	directions := []Direction{}
	if err := p.request(req, &directions); err != nil {
		return direction, err
	}
	for _, r := range directions {
		if strings.Contains(strings.ToLower(r.DirectionName), strings.ToLower(directionString)) {
			return r, nil
		}
	}

	return direction, fmt.Errorf("no direction found that matches %s", directionString)
}

func (p *MetroTransitProvider) getStop(stopName, routeID, directionID string) (Stop, error) {
	var stop Stop
	req, err := http.NewRequest("GET", p.getURL(fmt.Sprintf("stops/%s/%s", routeID, directionID)), nil)
	if err != nil {
		return stop, err
	}
	stops := []Stop{}
	if err := p.request(req, &stops); err != nil {
		return stop, err
	}
	for _, r := range stops {
		if strings.Contains(strings.ToLower(r.Description), strings.ToLower(stopName)) {
			return r, nil
		}
	}

	return stop, fmt.Errorf("no stop found that matches %s", stopName)
}

func (p *MetroTransitProvider) getNextDepature(routeID, directionID, placeCode string) (Departure, error) {
	var depature Departure
	req, err := http.NewRequest("GET", p.getURL(fmt.Sprintf("%s/%s/%s", routeID, directionID, placeCode)), nil)
	if err != nil {
		return depature, err
	}
	depatures := Departures{}
	if err := p.request(req, &depatures); err != nil {
		return depature, err
	}
	if len(depatures.Alerts) > 0 {
		for _, a := range depatures.Alerts {
			if a.StopClosed {
				return depature, fmt.Errorf("the stop is closed due to %s", a.AlertText)
			}
		}
	}
	deps := depatures.Depatures
	if len(deps) == 0 {
		return depature, fmt.Errorf("no upcoming depatures for this route")
	}
	sort.Slice(deps, func(i, j int) bool {
		return deps[i].DepartureTime < deps[j].DepartureTime
	})

	return deps[0], nil
}

func (p *MetroTransitProvider) GetNextTrip(route, stop, direction string) (int64, error) {
	r, err := p.findRoute(route)
	if err != nil {
		return 0, err
	}

	d, err := p.findDirection(direction, r.RouteID)
	if err != nil {
		return 0, err
	}

	s, err := p.getStop(stop, r.RouteID, strconv.FormatInt(d.DirectionID, 10))
	if err != nil {
		return 0, err
	}

	departure, err := p.getNextDepature(r.RouteID, strconv.FormatInt(d.DirectionID, 10), s.PlaceCode)
	if err != nil {
		return 0, err
	}

	return departure.DepartureTime, nil
}
