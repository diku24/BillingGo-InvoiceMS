package utlis

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func initConfig() {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error Loading Current working directory", err)
	}
	filePath := fmt.Sprintf("%v\\config\\", cwd)

	fmt.Println(filePath)

	viper.AddConfigPath(filePath)
	viper.SetConfigFile("config")
	viper.SetConfigType(".env")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading configuration variables", err)
	}
	viper.AutomaticEnv()
	serverPort := viper.GetString("SERVERPORT")
	fmt.Println(serverPort)
}
