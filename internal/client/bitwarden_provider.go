package client

var _ BitwardenProvider = bitwardenAPIClientProvider{}

type bitwardenAPIClientProvider struct{}

func NewBitwardenAPIClientProvider() bitwardenAPIClientProvider {
	return bitwardenAPIClientProvider{}
}

func (p bitwardenAPIClientProvider) NewBitwardenClient(apiURL string, identityURL string, accessToken string, organizationID string, stateFile string) (BitwardenClient, error) {
	return newBitwardenAPIClient(apiURL, identityURL, accessToken, organizationID, stateFile)
}

var _ BitwardenProvider = bitwardenMockClientProvider{}

type bitwardenMockClientProvider struct {
	expectedReturn map[string]string
	expectedErr    error
}

func NewBitwardenMockClientProvider(expectedReturn map[string]string, expectedErr error) bitwardenMockClientProvider {
	return bitwardenMockClientProvider{
		expectedReturn: expectedReturn,
		expectedErr:    expectedErr,
	}
}

func (p bitwardenMockClientProvider) NewBitwardenClient(apiURL string, identityURL string, accessToken string, organizationID string, stateFile string) (BitwardenClient, error) {
	return NewBitwardenMockClient(p.expectedReturn, p.expectedErr), nil
}
