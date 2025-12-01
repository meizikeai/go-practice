// internal/pkg/tencent/tencent.go
package tencent

import (
	"fmt"

	"go-practice/internal/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Client struct {
	SMS *sms.Client
	SES *ses.Client
	CFG *config.TencentCloudInstance
}

func NewClient(cfg *config.TencentCloudInstance) (*Client, error) {
	region := cfg.Region
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)

	smsClient, err := sms.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, fmt.Errorf("create sms client failed: %w", err)
	}

	sesClient, err := ses.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, fmt.Errorf("create ses client failed: %w", err)
	}

	return &Client{
		SMS: smsClient,
		SES: sesClient,
		CFG: cfg,
	}, nil
}
