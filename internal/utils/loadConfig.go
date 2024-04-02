package utils

type Config struct {
	DBUri           string `mapstructure:"MONGODB_URI"`
	MONGODB_DB      string `mapstructure:"MONGODB_DB"`
	USER_COLLECTION string `mapstructure:"USER_COLLECTION"`
}

// func LoadConfig(path string) (config Config, err error) {
// 	viper.AddConfigPath(path)
// 	viper.SetConfigType("env")
// 	viper.SetConfigName("app")

// 	viper.AutomaticEnv()

// 	err = viper.ReadInConfig()
// 	if err != nil {
// 		return
// 	}

// 	err = viper.Unmarshal(&config)
// 	return
// }
