package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
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

	faasCtxMsg, err := marshalEntity(faasCtx)
	if err != nil {
		return buildAPIGatewayProxyResponse(http.StatusInternalServerError, "marshal context failed"), nil
	}

	httpEventMsg, err := marshalEntity(httpEvent)
	if err != nil {
		return buildAPIGatewayProxyResponse(http.StatusInternalServerError, "marshal http event failed"), nil
	}

	if err := sendMessageToLark(ctx, fmt.Sprintf("context: %v, http event: %v", faasCtxMsg, httpEventMsg)); err != nil {
		return buildAPIGatewayProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return buildAPIGatewayProxyResponse(http.StatusOK, "{}"), nil
}

func buildAPIGatewayProxyResponse(statusCode int, msg string) scf.APIGatewayProxyResponse {
	body := map[string]string{
		"msg": msg,
	}
	bodyData, err := json.Marshal(body)
	if err != nil {
		bodyData = []byte{}
	}

	return scf.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			headers.ContentType: "application/json",
		},
		Body:            string(bodyData),
		IsBase64Encoded: false,
	}
}

func marshalEntity(entity any) (string, error) {
	data, err := json.Marshal(entity)
	return string(data), err
}

func sendMessageToLark(ctx context.Context, msg string) error {
	larkClient := lark.NewClient(larkAppId, larkAppSecret)

	content := larkim.NewTextMsgBuilder().Text(msg).Build()

	msgBody := larkim.NewCreateMessageReqBodyBuilder().
		ReceiveId(larkReceiverUserId).
		MsgType(larkim.MsgTypeText).
		Content(content).
		Build()
	req := larkim.NewCreateMessageReqBuilder().
		Body(msgBody).
		ReceiveIdType(larkim.ReceiveIdTypeUserId).
		Build()
	resp, err := larkClient.Im.Message.Create(ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 0 {
		return fmt.Errorf("status code: %v, message: %v", resp.StatusCode, resp.Msg)
	}

	return nil
}
