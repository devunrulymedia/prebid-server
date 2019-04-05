package unruly

import (
	"encoding/json"
	"github.com/mxmCherry/openrtb"
	"github.com/prebid/prebid-server/adapters"
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
	reqCopy.Imp = []openrtb.Imp{}
	reqCopy.Imp = append(reqCopy.Imp, imp)
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
