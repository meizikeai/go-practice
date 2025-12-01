// internal/pkg/tencent/ses.go
package tencent

import (
	"context"
	"encoding/json"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
)

type Email struct {
	client *Client
}

func (c *Client) SendEmail(ctx context.Context, to []string, templateID int64, variables map[string]string, subject string) (reqID string, err error) {
	req := ses.NewSendEmailRequest()
	req.FromEmailAddress = common.StringPtr(c.CFG.SES.FromEmail)
	req.Destination = common.StringPtrs(to)
	req.Template = &ses.Template{
		TemplateID: common.Uint64Ptr(uint64(templateID)),
	}

	if len(variables) > 0 {
		data := make(map[string]*string)
		for k, v := range variables {
			data[k] = common.StringPtr(v)
		}
		req.Template.TemplateData = common.StringPtr(marshal(data))
	}

	req.Subject = common.StringPtr(subject)

	resp, err := c.SES.SendEmailWithContext(ctx, req)
	if err != nil {
		return "", err
	}

	return *resp.Response.RequestId, nil
}

func marshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
