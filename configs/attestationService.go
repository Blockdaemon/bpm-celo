package configs

// AttesetationServiceEnvs Get envvars as string for env file.
const (
	AttesetationServiceEnvs = `ATTESTATION_SIGNER_ADDRESS={{ .Node.StrParameters.signer }}
CELO_VALIDATOR_ADDRESS={{ .Node.StrParameters.validator }}
CELO_PROVIDER={{ .Node.StrParameters.node_url }}
DATABASE_URL={{ .Node.StrParameters.db_host }}
SMS_PROVIDERS=twilio
TWILIO_MESSAGING_SERVICE_SID={{ .Node.StrParameters.twilio_service_sid }}
TWILIO_ACCOUNT_SID={{ .Node.StrParameters.twilio_account_sid }}
TWILIO_BLACKLIST={{ .Node.StrParameters.twilio_blacklist }}
TWILIO_AUTH_TOKEN={{ .Node.StrParameters.twilio_auth_token }}
PORT={{ .Node.StrParameters.port }}`

	PostgresEnvs = `POSTGRES_PASSWORD={{ .Node.StrParameters.db_password }}
POSTGRES_USER={{ .Node.StrParameters.db_user }}
POSTGRES_DATABASE=attestation-service`
)
