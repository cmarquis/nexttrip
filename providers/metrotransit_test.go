package providers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testAPIClient struct {
	json string
	err  error
	req  *http.Request
}

func (client *testAPIClient) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		Body:       ioutil.NopCloser(strings.NewReader(client.json)),
		StatusCode: http.StatusOK,
	}
	client.req = req
	return resp, client.err
}

func TestFindRoute(t *testing.T) {
	testClient := &testAPIClient{
		json: `
			[
				{
					"route_id": "4",
					"agency_id": 6,
					"route_label": "not right"
				},
				{
					"route_id": "1",
					"agency_id": 2,
					"route_label": "1234 Testroute 1234"
				}
			]
		`,
	}

	provider := MetroTransitProvider{
		APIClient: testClient,
	}

	r, err := provider.findRoute("testroute")
	require.Nil(t, err)

	require.Equal(t, "1", r.RouteID)
	require.Equal(t, int64(2), r.AgencyID)
	require.Equal(t, "1234 Testroute 1234", r.RouteLabel)

	r, err = provider.findRoute("bad route")
	require.EqualError(t, err, "no route found that matches bad route")
}

func TestFindDirection(t *testing.T) {

	testClient := &testAPIClient{
		json: `
			[
				{
					"direction_id": 0,
					"direction_name": "Northbound"
				},
				{
					"direction_id": 1,
					"direction_name": "Southbound"
				}
			]
		`,
	}

	provider := MetroTransitProvider{
		APIClient: testClient,
	}

	r, err := provider.findDirection("south", "1")
	require.Nil(t, err)

	require.Equal(t, int64(1), r.DirectionID)
	require.Equal(t, "Southbound", r.DirectionName)

	r, err = provider.findDirection("bad", "1")
	require.EqualError(t, err, "no direction found that matches bad")
}

func TestGetStop(t *testing.T) {

	testClient := &testAPIClient{
		json: `
			[
				{
					"place_code": "TF2",
					"description": "Target Field Station Platform 2"
				},
				{
					"place_code": "TF1",
					"description": "Target Field Station Platform 1"
				},
				{
					"place_code": "WARE",
					"description": "Warehouse District/ Hennepin Ave Station"
				}
			]
		`,
	}

	provider := MetroTransitProvider{
		APIClient: testClient,
	}

	r, err := provider.getStop("platform 2", "1", "1")
	require.Nil(t, err)

	require.Equal(t, "TF2", r.PlaceCode)

	r, err = provider.getStop("bad", "1", "1")
	require.EqualError(t, err, "no stop found that matches bad")
}

func TestGetNextDeparture(t *testing.T) {
	testClient := &testAPIClient{
		json: `
			{
				"stops": [
					{
						"stop_id": 56335,
						"latitude": 44.982905,
						"longitude": -93.277396,
						"description": "Target Field Station Platform 1"
					}
				],
				"alerts": [],
				"departures": [
					{
						"actual": false,
						"trip_id": "22847851-AUG22-RAIL-Weekday-03",
						"stop_id": 56335,
						"departure_text": "5:03",
						"departure_time": 1664229780,
						"description": "to Mall of America",
						"gate": "1",
						"route_id": "901",
						"route_short_name": "Blue",
						"direction_id": 1,
						"direction_text": "SB",
						"schedule_relationship": "NoData"
					},
					{
						"actual": false,
						"trip_id": "22847855-AUG22-RAIL-Weekday-03",
						"stop_id": 56335,
						"departure_text": "5:18",
						"departure_time": 1664230680,
						"description": "to Mall of America",
						"gate": "1",
						"route_id": "901",
						"route_short_name": "Blue",
						"direction_id": 1,
						"direction_text": "SB",
						"schedule_relationship": "NoData"
					}
				]
			}
		`,
	}

	provider := MetroTransitProvider{
		APIClient: testClient,
	}

	r, err := provider.getNextDepature("1", "1", "1")
	require.Nil(t, err)

	require.Equal(t, int64(1664229780), r.DepartureTime)
}

func TestGetNextDepartureNoDepartures(t *testing.T) {
	testClient := &testAPIClient{
		json: `
			{
				"stops": [
					{
						"stop_id": 56335,
						"latitude": 44.982905,
						"longitude": -93.277396,
						"description": "Target Field Station Platform 1"
					}
				],
				"alerts": [],
				"departures": []
			}
		`,
	}

	provider := MetroTransitProvider{
		APIClient: testClient,
	}

	_, err := provider.getNextDepature("1", "1", "1")
	require.EqualError(t, err, "no upcoming departures for this route")
}

func TestGetNextDepartureStopClosed(t *testing.T) {

	testClient := &testAPIClient{
		json: `
			{
				"stops": [
					{
						"stop_id": 56335,
						"latitude": 44.982905,
						"longitude": -93.277396,
						"description": "Target Field Station Platform 1"
					}
				],
				"alerts": [{
					"stop_closed": true,
					"alert_text": "bad weather"
				}],
				"departures": [
					{
						"actual": false,
						"trip_id": "22847851-AUG22-RAIL-Weekday-03",
						"stop_id": 56335,
						"departure_text": "5:03",
						"departure_time": 1664229780,
						"description": "to Mall of America",
						"gate": "1",
						"route_id": "901",
						"route_short_name": "Blue",
						"direction_id": 1,
						"direction_text": "SB",
						"schedule_relationship": "NoData"
					},
					{
						"actual": false,
						"trip_id": "22847855-AUG22-RAIL-Weekday-03",
						"stop_id": 56335,
						"departure_text": "5:18",
						"departure_time": 1664230680,
						"description": "to Mall of America",
						"gate": "1",
						"route_id": "901",
						"route_short_name": "Blue",
						"direction_id": 1,
						"direction_text": "SB",
						"schedule_relationship": "NoData"
					}
				]
			}
		`,
	}

	provider := MetroTransitProvider{
		APIClient: testClient,
	}

	_, err := provider.getNextDepature("1", "1", "1")
	require.EqualError(t, err, "the stop is closed due to bad weather")
}
