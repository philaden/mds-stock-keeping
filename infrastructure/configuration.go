package infrastructure

import "github.com/spf13/viper"

type AppConfiguration struct {
	DbDialect  string `mapstructure:"DB_DIALECT"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbUser     string `mapstructure:"DB_USER"`
	DbName     string `mapstructure:"DB_NAME"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbPort     uint   `mapstructure:"DB_PORT"`
	AppPort    uint   `mapstructure:"APP_PORT"`
}

func LoadConfiguration(path string) (config AppConfiguration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return AppConfiguration{}, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return AppConfiguration{}, err
	}
	return config, nil
}
