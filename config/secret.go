package config

type Secret struct {
	Postgres       PostgresSecret
	Authentication AuthenticationSecret
}

type PostgresSecret struct {
	Password string
	User     string
	DBName   string `envconfig:"db"`
}

type AuthenticationSecret struct {
	// Using RSA key to sign and verify the token. If both RSAKey and SecretKey
	// are provided, RSAKey will be used.
	TokenRSAPrivateKey string `envconfig:"token_rsa_private_key"`
	TokenRSAPublicKey  string `envconfig:"token_rsa_public_key"`

	// Use HMAC to sign and verify the token. Not support verifying at client.
	TokenHMACSecretKey string `envconfig:"token_hmac_secret_key"`
}
