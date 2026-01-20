package client

type BitwardenClient interface {
	ReadSecrets() (map[string]string, error)
}

type BitwardenProvider interface {
	NewBitwardenClient(apiURL string, identityURL string, accessToken string, organizationID string, stateFile string) (BitwardenClient, error)
}
