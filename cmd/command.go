package cmd

import (
	"fmt"
	"microsegement/serviceHandler"

	"github.com/spf13/cobra"
)



var rootCmd = &cobra.Command{
	Use:   "mircosegement",
	Short: "Mircosegement",
	Long:  `A tool to generate networkpolicy based on Java bytecode `,
	Run: func(cmd *cobra.Command, arg []string){
		myfolder := arg[0]
		fmt.Println("start to scan",myfolder)
		fmt.Println("------------------------------------------")

		k8sServiceList, err := serviceHandler.RegisterService(myfolder, "")
		if err != nil{
			fmt.Println("fail to service register%v\n", err)
		}
		err = serviceHandler.DiscoverService(k8sServiceList)
		if err != nil{
			fmt.Println("fail to service discoverr%v\n", err)
		}
		fmt.Println("------------------------------------------")
		fmt.Println("finish scan",myfolder)
	},
}

func Execute() error{
	return rootCmd.Execute()
} 
