package unruly

import (
	"encoding/json"
	"fmt"
	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
	"github.com/prebid/prebid-server/errortypes"
	"github.com/prebid/prebid-server/openrtb_ext"
	"net/http"
)

type UnrulyAdapter struct {
	http *adapters.HTTPAdapter
	URI  string
}

func (a *UnrulyAdapter) Name() string {
	return "unruly"
}

func (a *UnrulyAdapter) SkipNoCookies() bool {
	return false
}

func GetClient(config *adapters.HTTPAdapterConfig) *http.Client {
	return adapters.NewHTTPAdapter(config).Client
}

func NewUnrulyAdapter(config *adapters.HTTPAdapterConfig, endpoint string) *UnrulyAdapter {
	return NewUnrulyBidder(GetClient(config), endpoint)
}

func NewUnrulyBidder(client *http.Client, endpoint string) *UnrulyAdapter {
	clientAdapter := &adapters.HTTPAdapter{Client: client}

	return &UnrulyAdapter{
		http: clientAdapter,
		URI:  endpoint,
	}
}

func (a *UnrulyAdapter) ReplaceImp(imp openrtb.Imp, request *openrtb.BidRequest) *openrtb.BidRequest {
	reqCopy := *request
	reqCopy.Imp = append(make([]openrtb.Imp, 0, 1), imp)
	return &reqCopy
}

func (a *UnrulyAdapter) CheckImpExtension(request *openrtb.BidRequest) bool {
	var bidderExt adapters.ExtImpBidder
	return json.Unmarshal(request.Imp[0].Ext, &bidderExt) == nil
}

func (a *UnrulyAdapter) BuildRequest(request *openrtb.BidRequest) (*adapters.RequestData, []error) {
	var errs []error
	reqJSON, err := json.Marshal(request)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}
	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	headers.Add("Accept", "application/json")
	return &adapters.RequestData{
		Method:  "POST",
		Uri:     a.URI,
		Body:    reqJSON,
		Headers: headers,
	}, errs
}

func (a *UnrulyAdapter) MakeRequests(request *openrtb.BidRequest) ([]*adapters.RequestData, []error) {
	var errs []error
	var adapterRequests []*adapters.RequestData
	for _, imp := range request.Imp {
		newRequest := a.ReplaceImp(imp, request)
		if a.CheckImpExtension(newRequest) {
			adapterReq, errors := a.BuildRequest(newRequest)
			if adapterReq != nil {
				adapterRequests = append(adapterRequests, adapterReq)
			}
			errs = append(errs, errors...)
		}
	}
	return adapterRequests, errs
}

//func (a *UnrulyAdapter) MakeBids(internalRequest *openrtb.BidRequest, externalRequest *adapters.RequestData, response *adapters.ResponseData) (*adapters.BidderResponse, []error) {
//
//}

func getMediaTypeForImpWithId(impID string, imps []openrtb.Imp) (openrtb_ext.BidType, error) {
	for _, imp := range imps {
		if imp.ID == impID {
			if imp.Banner != nil {
				return openrtb_ext.BidTypeBanner, nil
			} else {
				return openrtb_ext.BidTypeVideo, nil
			}
		}
	}
	return "", &errortypes.BadInput{
		Message: fmt.Sprintf("Failed to find impression \"%s\" ", impID),
	}
}

func CheckResponse(response *adapters.ResponseData) error {
	if response.StatusCode != http.StatusOK {
		return &errortypes.BadServerResponse{
			Message: fmt.Sprintf("Unexpected status code: %d. Run with request.debug = 1 for more info", response.StatusCode),
		}
	}
	return nil
}
