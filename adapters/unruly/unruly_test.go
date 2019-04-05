package unruly

import (
	"encoding/json"
	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"net/http"
	"reflect"
	"testing"
)

func TestReturnsNewUnrulyBidderWithParams(t *testing.T) {
	mockClient := &http.Client{}
	mockAdapter := &adapters.HTTPAdapter{Client: mockClient}
	actual := *(NewUnrulyBidder(mockClient, "http://mockEndpoint.com"))
	expected := UnrulyAdapter{mockAdapter, "http://mockEndpoint.com"}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("actual = %v expected = %v", actual, expected)
	}
}

func TestBuildRequest(t *testing.T) {
	request := openrtb.BidRequest{}
	expectedJson, _ := json.Marshal(request)
	mockHeaders := http.Header{}
	mockHeaders.Add("Content-Type", "application/json;charset=utf-8")
	mockHeaders.Add("Accept", "application/json")
	data := adapters.RequestData{
		Method:  "POST",
		Uri:     "http://mockEndpoint.com",
		Body:    expectedJson,
		Headers: mockHeaders,
	}

	adapter := UnrulyAdapter{URI: "http://mockEndpoint.com"}

	actual, _ := adapter.BuildRequest(&request)
	expected := data
	if !reflect.DeepEqual(expected, *actual) {
		t.Errorf("actual = %v expected = %v", actual, expected)
	}

}

func TestReplaceImp(t *testing.T) {
	imp1 := openrtb.Imp{ID: "hello1"}
	imp2 := openrtb.Imp{ID: "hello2"}
	imp3 := openrtb.Imp{ID: "hello3"}
	newImp := openrtb.Imp{ID: "hello4"}
	request := openrtb.BidRequest{Imp: []openrtb.Imp{imp1, imp2, imp3}}
	adapter := UnrulyAdapter{URI: "http://mockEndpoint.com"}
	newRequest := adapter.ReplaceImp(newImp, &request)

	if len(newRequest.Imp) != 1 {
		t.Errorf("Size of Imp Array should be 1")
	}
	if !reflect.DeepEqual(request, openrtb.BidRequest{Imp: []openrtb.Imp{imp1, imp2, imp3}}) {
		t.Errorf("actual = %v expected = %v", request, openrtb.BidRequest{Imp: []openrtb.Imp{imp1, imp2, imp3}})
	}
	if !reflect.DeepEqual(newImp, newRequest.Imp[0]) {
		t.Errorf("actual = %v expected = %v", newRequest.Imp[0], newImp)
	}
}

func TestCheckImpExtension(t *testing.T) {
	adapter := UnrulyAdapter{URI: "http://mockEndpoint.com"}

	imp := openrtb.Imp{Ext: json.RawMessage(`{"bidder": {}}`)}
	request := openrtb.BidRequest{Imp: []openrtb.Imp{imp}}

	actual := adapter.CheckImpExtension(&request)
	expected := true

	if actual != expected {
		t.Errorf("actual = %v expected = %v", actual, expected)
	}
}

func TestCheckImpExtensionWithBadInput(t *testing.T) {
	adapter := UnrulyAdapter{URI: "http://mockEndpoint.com"}

	imp := openrtb.Imp{Ext: json.RawMessage(`{"bidder": notjson}`)}
	request := openrtb.BidRequest{Imp: []openrtb.Imp{imp}}

	actual := adapter.CheckImpExtension(&request)
	expected := false

	if actual != expected {
		t.Errorf("actual = %v expected = %v", actual, expected)
	}
}
func TestMakeRequests(t *testing.T) {
	adapter := UnrulyAdapter{URI: "http://mockEndpoint.com"}

	imp1 := openrtb.Imp{ID: "hello1", Ext: json.RawMessage(`{"bidder1": {}}`)}
	imp2 := openrtb.Imp{ID: "hello2", Ext: json.RawMessage(`{"bidder2": {}}`)}
	imp3 := openrtb.Imp{ID: "hello3", Ext: json.RawMessage(`{"bidder3": {}}`)}

	imps := []openrtb.Imp{imp1, imp2, imp3}

	inputRequest := openrtb.BidRequest{Imp: []openrtb.Imp{imp1, imp2, imp3}}
	actualAdapterRequests, _ := adapter.MakeRequests(&inputRequest)
	mockHeaders := http.Header{}
	mockHeaders.Add("Content-Type", "application/json;charset=utf-8")
	mockHeaders.Add("Accept", "application/json")
	if len(actualAdapterRequests) != 3 {
		t.Errorf("should have 3 imps")
	}
	for n, imp := range imps {
		request := openrtb.BidRequest{Imp: []openrtb.Imp{imp}}
		expectedJson, _ := json.Marshal(request)
		data := adapters.RequestData{
			Method:  "POST",
			Uri:     "http://mockEndpoint.com",
			Body:    expectedJson,
			Headers: mockHeaders,
		}
		if !reflect.DeepEqual(data, *actualAdapterRequests[n]) {
			t.Errorf("actual = %v expected = %v", *actualAdapterRequests[0], data)
		}
	}
}
