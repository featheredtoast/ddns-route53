package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type IpUpdater struct {
	Ip string
	RecordName string
	ZoneId string
}
func (u *IpUpdater) UpdateIp() (*route53.ChangeResourceRecordSetsOutput, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-west-2"),
		},
	})
	if err != nil {
		return nil, err
	}
	svc := route53.New(sess)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(u.RecordName),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(u.Ip),
							},
						},
						TTL:  aws.Int64(60),
						Type: aws.String("A"),
					},
				},
			},
			Comment: aws.String("DDNS update"),
		},
		HostedZoneId: aws.String(u.ZoneId),
	}

	result, err := svc.ChangeResourceRecordSets(input)
	if err != nil {
		return nil, err
	}
	return result, nil

}
