package cmd

import (
	"fmt"
	"log"
	"microsegement/serviceHandler"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configScan string
)

var rootCmd = &cobra.Command{
	Use:   "mircosegement",
	Short: "Mircosegement",
	Long:  `A tool to generate networkpolicy based on Java bytecode `,
	// Run: func(cmd *cobra.Command, arg []string) {
	// 	myfolder := arg[0]
	// 	fmt.Println("start to scan", myfolder)
	// 	fmt.Println("--------------------------------------------------------------------")

	// 	//扫描jar包，deployment文件，配置文件 获得微服务的基础信息
	// 	k8sServiceList, err := serviceHandler.RegisterService(myfolder, viper.GetString("configScan"))
	// 	// k8sServiceList, err := serviceHandler.RegisterService(myfolder, "nacos_yaml")
	// 	fmt.Println("scan config:___", viper.GetString("configScan"), "___")
	// 	if err != nil {
	// 		log.Fatalf("fail to service register %v\n", err)
	// 	}
	// 	//通过soot分析 获得完整的微服务调用信息
	// 	k8sServiceList, err = serviceHandler.DiscoverService(k8sServiceList)
	// 	if err != nil {
	// 		log.Fatalf("fail to service discovery %v\n", err)
	// 	}
	// 	//分析配置文件 获得完整的微服务调用外部ip的策略并生成策略文件
	// 	serviceHandler.GenerateIpPolicy(k8sServiceList)

	// 	fmt.Println("--------------------------------------------------------------------")
	// 	fmt.Println("finish scan", myfolder)
	// },
}

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "generate network policy",
	Long:  `generate network policy based on java bytecode`,
	Run: func(cmd *cobra.Command, arg []string) {
		myfolder := arg[0]
		fmt.Println("start to scan", myfolder)
		fmt.Println("--------------------------------------------------------------------")

		//扫描jar包，deployment文件，配置文件 获得微服务的基础信息
		k8sServiceList, err := serviceHandler.RegisterService(myfolder, viper.GetString("configScan"))
		// k8sServiceList, err := serviceHandler.RegisterService(myfolder, "nacos_yaml")
		fmt.Println("scan config:___", viper.GetString("configScan"), "___")
		if err != nil {
			log.Fatalf("fail to service register %v\n", err)
		}
		//通过soot分析 获得完整的微服务调用信息
		k8sServiceList, err = serviceHandler.DiscoverService(k8sServiceList)
		if err != nil {
			log.Fatalf("fail to service discovery %v\n", err)
		}
		//分析配置文件 获得完整的微服务调用外部ip的策略并生成策略文件
		serviceHandler.GenerateIpPolicy(k8sServiceList)

		fmt.Println("--------------------------------------------------------------------")
		fmt.Println("finish scan", myfolder)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "get nacos config information",
	Long:  `Get nocas config information, especially namespace`,
	Run: func(cmd *cobra.Command, arg []string) {
		myfolder := arg[0]
		serviceHandler.GetNacosNameSpace(myfolder)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	//如果未设置 则configScan为""
	//带p则5个参数，不带p则4个参数

	rootCmd.AddCommand(policyCmd)
	rootCmd.AddCommand(configCmd)
	policyCmd.PersistentFlags().StringVarP(&configScan, "configScan", "c", "", "if need to scan extra config files, input their path")
}

func initConfig() {
	viper.BindPFlag("configScan", policyCmd.PersistentFlags().Lookup("configScan"))
}

func Execute() error {
	return rootCmd.Execute()
}
