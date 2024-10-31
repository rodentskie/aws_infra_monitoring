package awsquery

import (
	"context"
	"log"
	"packages/library/go/structs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
)

func GetACMCertificates(ctx context.Context, region string) ([]structs.ACMCertificate, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	client := acm.NewFromConfig(cfg)

	input := &acm.ListCertificatesInput{
		MaxItems: aws.Int32(100),
	}

	var certificates []structs.ACMCertificate
	paginator := acm.NewListCertificatesPaginator(client, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, cert := range output.CertificateSummaryList {
			certDetail, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
				CertificateArn: cert.CertificateArn,
			})
			if err != nil {
				log.Printf("Error getting details for certificate %s: %v", *cert.CertificateArn, err)
				continue
			}

			tagsOutput, err := client.ListTagsForCertificate(ctx, &acm.ListTagsForCertificateInput{
				CertificateArn: cert.CertificateArn,
			})

			tags := make(map[string]string)
			if err == nil {
				for _, tag := range tagsOutput.Tags {
					tags[*tag.Key] = *tag.Value
				}
			}

			certificate := structs.ACMCertificate{
				CertificateArn:     *cert.CertificateArn,
				DomainName:         *cert.DomainName,
				RenewalEligibility: string(certDetail.Certificate.RenewalEligibility),
				Status:             string(certDetail.Certificate.Status),
				Type:               string(certDetail.Certificate.Type),
				Tags:               tags,
			}

			if certDetail.Certificate.NotBefore != nil {
				notBefore := certDetail.Certificate.NotBefore.Format("2006-01-02")
				certificate.NotBefore = &notBefore
			}
			if certDetail.Certificate.NotAfter != nil {
				notAfter := certDetail.Certificate.NotAfter.Format("2006-01-02")
				certificate.NotAfter = &notAfter
			}

			certificates = append(certificates, certificate)
		}
	}

	return certificates, nil
}
