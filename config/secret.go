package config

type Secret struct {
	Postgres       PostgresSecret       `envconfig:"postgres"`
	Authentication AuthenticationSecret `envconfig:"auth"`
	Redis          RedisSecret          `envconfig:"redis"`
}

type PostgresSecret struct {
	User     string `envconfig:"user"`
	Password string `envconfig:"password"`
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

type RedisSecret struct {
	Username string `envconfig:"username"`
	Password string `envconfig:"password"`
}
