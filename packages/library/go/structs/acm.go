package structs

type ACMCertificate struct {
	CertificateArn     string            `json:"certificateArn"`
	DomainName         string            `json:"domainName"`
	Status             string            `json:"status"`
	Type               string            `json:"type"`
	RenewalEligibility string            `json:"renewalEligibility"`
	NotBefore          *string           `json:"notBefore"`
	NotAfter           *string           `json:"notAfter"`
	Tags               map[string]string `json:"tags"`
}

type ACMCertificateResponse struct {
	Data []ACMCertificate `json:"data"`
}
