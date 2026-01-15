# conflux

conflux is a simple Go library to read configuration values from multiple sources such as YAML files, environment variables, and Bitwarden secrets.

## Usage

```go
// Define a struct with required fields
type proxmox struct {
	SSHPublicKeyPath     string `json:"ssh_public_key_path" required:"true"`
	NodeCIDRAddress      string `json:"node_cidr_address" required:"true"`
	GatewayAddress       string `json:"gateway_address" required:"true"`
}

// Create a config mux that reads from two YAML files, the environment, and Bitwarden
configMux := conflux.NewConfigMux(
  conflux.WithYAMLFileReader("config/all.yml", conflux.WithPath("config/other.yml")),
  conflux.WithEnvReader(),
  conflux.WithBitwardenSecretReader(),
)

// Read configuration values into a struct
var proxmoxConfig proxmox
diagnostics, err := conflux.Unmarshal(configMux, &proxmoxConfig)
if errors.Is(err, conflux.ErrInvalidFields) {
  return fmt.Errorf("invalid or missing config fields:\n%s", conflux.DiagnosticsToTable(diagnostics))
} else if err != nil {
  return fmt.Errorf("failed to unmarshal config: %w", err)
}

// at this point, `proxmoxConfig` is guaranteed to have all required fields set to some non-zero value
```

## Why use this library instead of Viper, Koanf, etc?

Use this library if you want something that:
- Is light (751 SLOC)
- Supports reading from Bitwarden Secrets
- Built-in validation. The `required` tag allows `conflux` to give you an exact report of the configurations that were found and missing. This report can be printed as a table for a user-friendly experience.
- Is easily extensible. You can easily your own `Reader`s that read from any source you wish.
- Has flexible initialization logic. `Reader`s can be initialized lazily. This allows us to initialize a `Reader` with a map of configs that have been read so far. The `BitwardenSecretReader` actually uses the configs found by `YAMLFileReader` and `EnvReader` to authenticate to Bitwarden.
- Has flexible validation logic. If your struct has some specific validation rules, `conflux` will run them if you define a receiver with the following signature: `Validate(map[string]string) bool`.
- Has flexible struct-filling logic. If your struct needs to fill-in additional fields after the required fields have been filled, `conflux` will fill those fields for you if you define a receiver with the following signature: `FillInKeys() error`.

## Conflicts

- Sometimes there are two sources that define different values for the same configuration key. Conflux has two rules for resolving conflicts:
  - When initializing `ConfigMux`, `Reader`s should be passed as arguments from lowest to highest priority. For example, in the usage above, conflicting values in `BitwardenSecretReader` will override env values, which will in turn override file values.
  - When initializing `YAMLFileReader`, paths should be passed as arguments from lowest to highest priority. The path with the lowest priority is the one that is passed as the first argument. The path with the next-lowest priority is the one that is passed-in as the first functional option: `WithPath("some-file-here.yml")`, and so on.
  - If you pass-in a directory to `YAMLFileReader`, `YAMLFileReader` will read the files in that directory recursively in lexicographic order. That is to say that a file `A` that falls lexicographically before another file `B`, will have lower priority that `B`.

## Name

I chose `conflux` for this library because it's Latin for "flowing together".

This word kind-of describes the logic this library has of merging multiple configuration sources into a single source of truth.

## Installation

```sh
go get github.com/dannyvelas/conflux
```

## Testing

```sh
go test ./...
```

## License

MIT License. See [LICENSE](LICENSE) for details.
