package config

import "github.com/pkg/errors"

type TwitterConfig struct {
	ConsumerKey       string `env:"TWITTER_CONSUMER_KEY"`
	ConsumerSecret    string `env:"TWITTER_CONSUMER_SECRET"`
	AccessToken       string `env:"TWITTER_ACCESS_TOKEN"`
	AccessTokenSecret string `env:"TWITTER_ACCESS_TOKEN_SECRET"`
}

func (c TwitterConfig) Validate() error {
	if c.ConsumerKey == "" {
		return errors.Errorf("missing consumer key")
	}
	if c.ConsumerSecret == "" {
		return errors.Errorf("missing consumer secret")
	}
	if c.AccessToken == "" {
		return errors.Errorf("missing access token")
	}
	if c.AccessTokenSecret == "" {
		return errors.Errorf("missing access token secret")
	}

	return nil
}
