package config

type HabApiConfig struct {
	Address      string
	Port         int
	WebPort      int
	IncludeDebug bool

	DB     DBConfig
	OpenAI struct {
		ApiKey string
	}
}
