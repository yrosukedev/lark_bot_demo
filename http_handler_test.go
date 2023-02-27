package main

import (
	"context"
	"fmt"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"github.com/tencentyun/scf-go-lib/functioncontext"
	"testing"
)

func TestHttpHandler(t *testing.T) {

	faasCtx := functioncontext.FunctionContext{
		RequestID:       "23a1b77297ab9167d5aa049021f6555d",
		Namespace:       "default",
		FunctionName:    "lark_notification-1676701322",
		FunctionVersion: "$LATEST",
		MemoryLimitInMb: 128,
		TimeLimitInMs:   3000,
		Environment: map[string]string{
			"SCF_NAMESPACE": "default",
		},
		TencentcloudRegion: "ap-beijing",
		TencentcloudAppID:  "1316902468",
		TencentcloudUin:    "100029648803",
	}

	httpEvent := scf.APIGatewayProxyRequest{
		Path:       "/lark_notification-1676701322",
		HTTPMethod: "POST",
		Headers: map[string]string{
			"accept-encoding":          "gzip",
			"content-length":           "1274",
			"content-type":             "application/json;charset=utf-8",
			"endpoint-timeout":         "1800",
			"host":                     "service-cbfce0qx-1316902468.bj.apigw.tencentcs.com",
			"requestsource":            "APIGW",
			"unit":                     "eu_nc user-agent:Go-http-client/1.1",
			"x-api-requestid":          "23a1b77297ab9167d5aa049021f6555d",
			"x-api-scheme":             "https",
			"x-b3-traceid":             "23a1b77297ab9167d5aa049021f6555d",
			"x-lark-request-nonce":     "533226633",
			"x-lark-request-timestamp": "1677254849",
			"x-lark-signature":         "9e0de55dbd5e56a11e747bcede2044af53e8fca472fa817db7b40360c80e776d",
			"x-qualifier":              "$DEFAULT",
			"x-request-id":             "1e81ee74-2c3c-473c-8c17-58304cbcbafb",
		},
		RequestContext: scf.APIGatewayProxyRequestContext{
			ServiceID:       "service-cbfce0qx",
			Path:            "/lark_notification-1676701322",
			HTTPMethod:      "ANY",
			Stage:           "test",
			SourceIP:        "123.58.10.238",
			WebsocketEnable: false,
		},
		Body:            "{\"encrypt\":\"dAG1vSy7t3kK371EjftOlhAEDUFWigoiaRu96sh4MJF49CLWcALmaHkbVzO91RZ/wXSsNEnCOuc2CL6w7c3+h1ECGf0MGvweAG3jgGfd8sbwvSOYG5k4cOx1bWVGJPS2QQT/0eGR1WLqvTYzBZkDIFQN1hJa9BwbJ4+nKlhAOZIVU6lZn1L/CGLd6pr/oyVKtc9fiG5GWbvg1famE6JA2gVho49zqPMB9BCx7DIqTHEijcVDssCdv5WBXSNXCoSE6gyxajqFwTihetBF8/YDFsXywfI+X5u080FFCD9rbF7DryAC7mEJFrfu3A3go1gBbQFc2HaP9UTseLxHlytS/nDEJafejLQKaIOg1pdQK4zC8f4NKDzZf5YD/TX4/vNFXD4mlnEheorQzyWMy3yp8m+1pl/eGpG1N8wgDfNSYBntKaK8Ceah2q2NSt0XfhlXWBf8gmUgj0cq6yIQqwgsXCb0XoZst8BLdRxmDwsfzRmZrihRkBrTZnirU2aPaRir5zHkh8EHSuCM09vQ9fTRpvsqXxuyNksRfSEntA7wwUm1i3S/Tg8iMAx5bdZ/MTmhM+vqoLkcrFYAjIOpZodpNQYIRg9RUN6RHECEjccdecBWbu+6kyZyBflm2pCwNu8SZJgWziMXdWfDvAKjOf77DXacgMhpY1doh06TVEym61Cs//Oj8jBf58SIS+fe8ktkC1u4Ams+rGOSOvXZYDisEJv4adWzD1+SkuceHWZJK+y9SHj93WVjb3GUqyxZQ1pqktiobgMFgngrLDmNgJLiewRtGNwMYgZnmdD+eXC1UkqZuQUijpcU6A7TvLH2/Pxe3dGZSCkc0LSwuw0C+N1okknas6THGj0+quIunTzPb7uwHV9+dfLpLB5axIrVDbj9H8Np60ohZS0u1YrXw0j4VKXFq2O8fs1NxZd174XCG1HbKBNitJoGKRHZoSTnLYjpRJtm5AKTF4GSApAXPjniFHnu5l4hIibPg0xfSiM+sl7QyHusS0mag3p+J08iO4C2xmU7qkmxGA3X5tyfWl4ostOcdQ5pU6EksSbS+3Y3rVITec2oacBbnf23BM9IA2dICkLJrq89Fmhifa92x9GLhzzeHFFC7YrVOAdSizWz939TdIDTVr9csCd5Yweg4Ew/9Fa3ybRkoEpwTgvn8La8NKXbFoCzcVCwNJ1ImEK3zE3kXjbmwK5YkF7vAT+Ai3xPmldZl52JnAqunljsKNKAsB+jHhPcd0ppJJSQfGClwXg=\"}",
		IsBase64Encoded: false,
	}

	resp, err := httpHandler(functioncontext.NewContext(context.Background(), &faasCtx), httpEvent)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("resp: %+v", resp)
}
