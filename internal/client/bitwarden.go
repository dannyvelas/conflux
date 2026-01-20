package client

import (
	"fmt"

	"github.com/bitwarden/sdk-go"
)

var _ BitwardenClient = bitwardenAPIClient{}

type bitwardenAPIClient struct {
	organizationID string
	client         sdk.BitwardenClientInterface
}

func newBitwardenAPIClient(apiURL, identityURL, accessToken, organizationID, stateFile string) (bitwardenAPIClient, error) {
	bitwardenClient, err := sdk.NewBitwardenClient(&apiURL, &identityURL)
	if err != nil {
		return bitwardenAPIClient{}, fmt.Errorf("error initializing bitwarden client: %v", err)
	}

	if err := bitwardenClient.AccessTokenLogin(accessToken, &stateFile); err != nil {
		return bitwardenAPIClient{}, fmt.Errorf("error logging in to bitwarden client: %v", err)
	}

	return bitwardenAPIClient{
		organizationID: organizationID,
		client:         bitwardenClient,
	}, nil
}

func (c bitwardenAPIClient) ReadSecrets() (map[string]string, error) {
	m := make(map[string]string)

	secrets := c.client.Secrets()
	listResponse, err := secrets.List(c.organizationID)
	if err != nil {
		return nil, fmt.Errorf("error listing secrets: %v", err)
	}

	for _, secret := range listResponse.Data {
		secretData, err := secrets.Get(secret.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting secret: %v", err)
		}

		m[secret.Key] = secretData.Value
	}

	return m, nil
}
