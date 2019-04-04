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

func TestMakeRequest(t *testing.T) {
	request := openrtb.BidRequest{
		ID: "testID",
	}
	reqJSON, _ := json.Marshal(&request)
	mockHeaders := http.Header{}
	mockHeaders.Add("Content-Type", "application/json;charset=utf-8")
	mockHeaders.Add("Accept", "application/json")
	data := adapters.RequestData{
		Method:  "POST",
		Uri:     "http://mockEndpoint.com",
		Body:    reqJSON,
		Headers: mockHeaders,
	}

	mockClient := &http.Client{}
	mockAdapter := &adapters.HTTPAdapter{Client: mockClient}

	adapter := UnrulyAdapter{mockAdapter, "http://mockEndpoint.com"}

	actual, _ := adapter.makeRequest(&request)
	expected := data

	if !reflect.DeepEqual(expected, *actual) {
		t.Errorf("actual = %v expected = %v", actual, expected)
	}

}
