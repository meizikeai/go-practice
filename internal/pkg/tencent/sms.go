// internal/pkg/tencent/sms.go
package tencent

import (
	"context"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type SMS struct {
	client *Client
}

func (c *Client) SendSMS(ctx context.Context, phone, templateID string, params []string) (reqID string, err error) {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = common.StringPtr(c.CFG.SMS.SdkAppID)
	req.SignName = common.StringPtr(c.CFG.SMS.SignName)
	req.TemplateId = common.StringPtr(templateID)
	req.PhoneNumberSet = common.StringPtrs([]string{phone})
	req.TemplateParamSet = common.StringPtrs(params)

	resp, err := c.SMS.SendSmsWithContext(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Response.SendStatusSet) == 0 {
		return "", fmt.Errorf("empty response")
	}

	status := resp.Response.SendStatusSet[0]
	if *status.Code != "Ok" {
		return *status.SerialNo, fmt.Errorf("send failed: %s", *status.Message)
	}

	return *status.SerialNo, nil
}
