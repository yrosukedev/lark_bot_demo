package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-http-utils/headers"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkvc "github.com/larksuite/oapi-sdk-go/v3/service/vc/v1"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/functioncontext"
	"net/http"
	"strings"
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

	DebugLogger.Printf("[show http event] HTTP Event: %+v", httpEvent)

	eventHandler := dispatcher.
		NewEventDispatcher(larkVerificationToken, larkEventEncryptKey).
		OnP2MeetingAllMeetingEndedV1(func(ctx context.Context, event *larkvc.P2MeetingAllMeetingEndedV1) error {

			InfoLogger.Printf("[eventHandler OnP2MeetingAllMeetingEndedV1] begin event handler")

			meetingEventMsg, err := marshalEntity(event)
			if err != nil {
				ErrorLogger.Printf("[eventHandler OnP2MeetingAllMeetingEndedV1] marshal event vc.meeting.all_meeting_ended_v1 failed, error: %+v", err)
				return fmt.Errorf("marshal event vc.meeting.all_meeting_ended_v1 failed, error: %+v", err)
			}

			InfoLogger.Printf("[eventHandler OnP2MeetingAllMeetingEndedV1] before sending message to lark, lark event: %v", meetingEventMsg)

			if err := sendMessageToLark(ctx, fmt.Sprintf("lark event: %v", meetingEventMsg)); err != nil {
				ErrorLogger.Printf("[httpHandler] sending message to lark failed, error: %+v", err)
				return fmt.Errorf("sending message to lark failed, error: %+v", err)
			}

			InfoLogger.Printf("[eventHandler OnP2MeetingAllMeetingEndedV1] after sending message to lark")

			return nil
		}).
		OnP2MeetingEndedV1(func(ctx context.Context, event *larkvc.P2MeetingEndedV1) error {

			InfoLogger.Printf("[eventHandler OnP2MeetingEndedV1] begin event handler")

			meetingEventMsg, err := marshalEntity(event)
			if err != nil {
				ErrorLogger.Printf("[eventHandler OnP2MeetingEndedV1] marshal event vc.meeting.meeting_ended_v1 failed, error: %+v", err)
				return fmt.Errorf("marshal event vc.meeting.all_meeting_ended_v1 failed, error: %+v", err)
			}

			InfoLogger.Printf("[eventHandler OnP2MeetingEndedV1] before sending message to lark, lark event: %v", meetingEventMsg)

			if err := sendMessageToLark(ctx, fmt.Sprintf("lark event: %v", meetingEventMsg)); err != nil {
				ErrorLogger.Printf("[httpHandler] sending message to lark failed, error: %+v", err)
				return fmt.Errorf("sending message to lark failed, error: %+v", err)
			}

			InfoLogger.Printf("[eventHandler OnP2MeetingEndedV1] after sending message to lark")

			return nil
		})

	InfoLogger.Printf("[httpHandler] lark event handler created")

	larkEventReq := &larkevent.EventReq{
		Header:     larkEventHeadersFrom(httpEvent.Headers),
		RequestURI: httpEvent.Path,
		Body:       []byte(httpEvent.Body),
	}

	InfoLogger.Printf("[httpHandler] before handling lark request, request: %+v", larkEventReq)

	larkResp := eventHandler.Handle(ctx, larkEventReq)

	InfoLogger.Printf("[httpHandler] after handling lark request, response: %+v", larkResp)

	return buildAPIGatewayProxyResponseWithLarkEventResponse(larkResp), nil
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

func buildAPIGatewayProxyResponseWithLarkEventResponse(larkEventResp *larkevent.EventResp) scf.APIGatewayProxyResponse {
	return scf.APIGatewayProxyResponse{
		StatusCode:      larkEventResp.StatusCode,
		Headers:         larkEventHeadersTo(larkEventResp.Header),
		Body:            string(larkEventResp.Body),
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

func larkEventHeadersFrom(headers map[string]string) http.Header {
	larkHeaders := make(http.Header)

	for k, v := range headers {
		values := strings.Split(v, ",")

		for _, v := range values {
			larkHeaders.Add(k, strings.TrimSpace(v))
		}
	}

	return larkHeaders
}

func larkEventHeadersTo(larkHeaders http.Header) map[string]string {
	headers := make(map[string]string)

	for k, v := range larkHeaders {
		headers[k] = strings.Join(v, ", ")
	}

	return headers
}
