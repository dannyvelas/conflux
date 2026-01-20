package conflux

import (
	"errors"
	"fmt"

	"github.com/dannyvelas/conflux/internal/client"
)

var _ Reader = (*bitwardenSecretReader)(nil)

type bitwardenSecretReader struct {
	mapReader               mapReader
	bitwardenClientProvider client.BitwardenProvider
}

// NewBitwardenSecretReader creates a new Bitwarden secret reader using the provided config map to authenticate to Bitwarden.
func NewBitwardenSecretReader(configMap map[string]string, opts ...func(*bitwardenSecretReader)) *bitwardenSecretReader {
	bitwardenSecretReader := &bitwardenSecretReader{
		mapReader:               newMapReader(configMap),
		bitwardenClientProvider: client.NewBitwardenAPIClientProvider(),
	}

	for _, opt := range opts {
		opt(bitwardenSecretReader)
	}

	return bitwardenSecretReader
}

func (r *bitwardenSecretReader) Read() (ReadResult, error) {
	config := newBitwardenConfig()

	diagnostics, err := Unmarshal(r.mapReader, &config)
	if errors.Is(err, ErrInvalidFields) {
		return NewDiagnosticReadResult(nil, diagnostics), ErrInvalidFields
	} else if err != nil {
		return nil, fmt.Errorf("error unmarshalling bitwarden creds: %v", err)
	}

	bitwardenClient, err := r.bitwardenClientProvider.NewBitwardenClient(
		config.APIURL,
		config.IdentityURL,
		config.AccessToken,
		config.OrganizationID,
		config.StateFilePath,
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing bitwarden client: %v", err)
	}

	bitwardenSecrets, err := bitwardenClient.ReadSecrets()
	if err != nil {
		return nil, fmt.Errorf("error reading bitwarden secrets: %v", err)
	}

	return NewDiagnosticReadResult(bitwardenSecrets, diagnostics), nil
}

func WithProvider(provider client.BitwardenProvider) func(*bitwardenSecretReader) {
	return func(bitwardenSecretReader *bitwardenSecretReader) {
		bitwardenSecretReader.bitwardenClientProvider = provider
	}
}
