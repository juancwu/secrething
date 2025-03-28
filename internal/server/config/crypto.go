package config

// CryptoConfig holds cryptographic settings used throughout the application.
// It stores encryption keys and other cryptography-related configuration.
type CryptoConfig struct {
	// AesKey is the main encryption key used for general data encryption with AES.
	// In production, this should be a strong, randomly generated key.
	AesKey string `env:"CRYPTO_AES_KEY" env-required:"" env-description:"The general AES key use to encrypt data."`
}

// cryptoCfg is the private singleton instance of cryptographic configuration.
// It's populated when Load() is called.
var cryptoCfg CryptoConfig

// Crypto returns a copy of the initialized cryptographic configuration.
// The configuration must be loaded with Load() before this function is called.
func Crypto() CryptoConfig {
	return cryptoCfg
}
