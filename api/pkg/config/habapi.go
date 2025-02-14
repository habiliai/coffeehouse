package config

type HabApiConfig struct {
	Address      string `env:"ADDRESS"`
	Port         int    `env:"PORT"`
	WebPort      int    `env:"WEB_PORT"`
	IncludeDebug bool   `env:"INCLUDE_DEBUG"`

	DB     DBConfig
	OpenAI OpenAIConfig
}

type OpenAIConfig struct {
	ApiKey string `env:"OPENAI_API_KEY"`
}
