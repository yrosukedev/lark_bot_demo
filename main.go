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
	InfoLogger.Printf("[httpHandler] begin http handler")

	faasCtx, ok := functioncontext.FromContext(ctx)
	if !ok {
		ErrorLogger.Printf("[show context] failed to convert FunctionContext from context.")
		return buildAPIGatewayProxyResponse(http.StatusInternalServerError, "failed to convert FunctionContext from context"), nil
	}

	DebugLogger.Printf("[show context] Function Context: %+v", faasCtx)

	faasCtxMsg, err := marshalEntity(faasCtx)
	if err != nil {
		ErrorLogger.Printf("[httpHandler] marshal context failed, error: %+v", err)
		return buildAPIGatewayProxyResponse(http.StatusInternalServerError, "marshal context failed"), nil
	}

	DebugLogger.Printf("[show http event] HTTP Event: %+v", httpEvent)

	httpEventMsg, err := marshalEntity(httpEvent)
	if err != nil {
		ErrorLogger.Printf("[show http event] marshal httpEvent failed, error: %+v", err)
		return buildAPIGatewayProxyResponse(http.StatusInternalServerError, "marshal http event failed"), nil
	}

	InfoLogger.Printf("[send message] before sending message to lark, context: %v, http event: %v", faasCtxMsg, httpEventMsg)

	if err := sendMessageToLark(ctx, fmt.Sprintf("context: %v, http event: %v", faasCtxMsg, httpEventMsg)); err != nil {
		ErrorLogger.Printf("[httpHandler] sending message to lark failed, error: %+v", err)
		return buildAPIGatewayProxyResponse(http.StatusInternalServerError, err.Error()), nil
	}

	InfoLogger.Printf("[send message] after sending message to lark")

	return buildAPIGatewayProxyResponse(http.StatusOK, ""), nil
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

	content, err := marshalEntity(map[string]string{
		"text": msg,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal msg to json format content, %+v", err)
	}

	DebugLogger.Printf("[send message to lark] content: %v", content)

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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %v, code: %v, message: %v", resp.StatusCode, resp.Code, resp.Msg)
	}

	if resp.Code != 0 {
		return fmt.Errorf("status code: %v, code: %v, message: %v", resp.StatusCode, resp.Code, resp.Msg)
	}

	return nil
}
