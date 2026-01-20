package client

var _ BitwardenClient = bitwardenMockClient{}

type bitwardenMockClient struct {
	expectedReturn map[string]string
	expectedErr    error
}

func NewBitwardenMockClient(expectedReturn map[string]string, expectedErr error) bitwardenMockClient {
	return bitwardenMockClient{
		expectedReturn: expectedReturn,
		expectedErr:    expectedErr,
	}
}

func (c bitwardenMockClient) ReadSecrets() (map[string]string, error) {
	return c.expectedReturn, c.expectedErr
}
