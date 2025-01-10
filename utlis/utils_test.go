package utlis

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestEnv(t *testing.T) {
	initConfig()

	//Accessing Variables
	serverPort := viper.GetString("SERVERPORT")
	dbUser := viper.GetString("DATABASEUSER")
	dbPass := viper.GetString("DATABASEPASS")
	dbName := viper.GetString("DATABASE")

	fmt.Println("serverPort :", serverPort)
	fmt.Println("DB User  :", dbUser)
	fmt.Println("DB Pass :", dbPass)
	fmt.Println("DB Pass  :", dbName)

}
