package config

type HabApiConfig struct {
	Address string
	Port    int
	WebPort int

	DB     DBConfig
	OpenAI struct {
		ApiKey string
	}
}
