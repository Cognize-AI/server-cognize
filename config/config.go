package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PORT                    string `mapstructure:"PORT"`
	DbString                string `mapstructure:"DB_STRING"`
	JwtSecret               string `mapstructure:"JWT_SECRET"`
	GoogleOAuthClientID     string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOAuthClientSecret string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectURL  string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
	Environment             string `mapstructure:"ENVIRONMENT"`
	SMTPHost                string `mapstructure:"SMTP_HOST"`
	SMTPPort                int    `mapstructure:"SMTP_PORT"`
	SMTPUsername            string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword            string `mapstructure:"SMTP_PASSWORD"`
	AxiomToken              string `mapstructure:"AXIOM_TOKEN"`
	AxiomOrg                string `mapstructure:"AXIOM_ORG"`
	AxiomDataset            string `mapstructure:"AXIOM_DATASET"`
	EncSecret               string `mapstructure:"ENC_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	//viper.SetConfigName("app")
	//viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
