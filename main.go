package main

import (
	"context"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/functioncontext"
	"net/http"
)

func main() {
	cloudfunction.Start(httpHandler)
}

func httpHandler(ctx context.Context, httpEvent scf.APIGatewayProxyRequest) (resp scf.APIGatewayProxyResponse, err error) {
	faasCtx, ok := functioncontext.FromContext(ctx)
	if ok {
		fmt.Printf("%+v", faasCtx)
	}

	fmt.Printf("%+v", httpEvent)

	resp = scf.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			headers.ContentType: "application/json",
		},
		Body:            "{}",
		IsBase64Encoded: false,
	}

	return resp, nil
}
