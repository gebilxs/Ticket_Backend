package cmd

import (
	"log"
	"net/http"
	"os"
	"ticket-backend/config"
	db "ticket-backend/database"
	"ticket-backend/model"
	"ticket-backend/router"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	logger  = &logrus.Logger{}
	rootCmd = &cobra.Command{}
)

func initConfig() {
	config.MustInit(os.Stdout, cfgFile) // 配置初始化
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config/dev.yaml", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().Bool("debug", true, "开启debug")
	viper.SetDefault("gin.mode", rootCmd.PersistentFlags().Lookup("debug"))
}

func Execute() error {
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, err := db.Mysql(
			viper.GetString("db.hostname"),
			viper.GetInt("db.port"),
			viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.dbname"),
		)
		if err != nil {
			return err
		}
		db.DB.AutoMigrate(&model.User{}, &model.Passenger{}, &model.Admin{}, &model.Company{}, &model.Plane{}, &model.Ticket{})
		//AutoMigrate 陆续读取相关的数据
		defer db.DB.Close()

		r := router.NewRouter()
		//加载路由

		port := viper.GetString("server.port")
		//viper 包读取相关配置
		r.Run(port)
		log.Println("port = *** =", port)
		return http.ListenAndServe(port, nil) // listen and serve
	}

	return rootCmd.Execute()

}
