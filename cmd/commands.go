package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	funcName string
	cfgFile string
	env string
	sourceZip string
	invokePayloadFile string
	jfConfig *Config
)

var uploadCmd = &cobra.Command{
    Use:   "upload",
	Short: "upload function with funcName, zipFile and env",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := upload(jfConfig.ManagerAddr, funcName, sourceZip, env)
		if err != nil {
			cmd.PrintErr(err)
		} else {
			cmd.Println(res)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:	"delete",
	Short: "delete function",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var invokeCmd = &cobra.Command{
	Use: "invoke",
	Short: "invoke a function with name and payload",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var rootCmd = &cobra.Command{
	Use:   "jf",
	Short: "jf is a simple cli for jointFaaS management",
	PreRun: func(cmd *cobra.Command, args []string) {
		
	},
}

func rootInit() {
	cobra.OnInitialize(initConfig)

	uploadCmd.Flags().StringVarP(&env, "env", "e", "python3", "function env.")
	uploadCmd.Flags().StringVarP(&sourceZip, "sourceZip", "z", "", "source code file path.")
	uploadCmd.MarkFlagRequired("sourceZip")

	invokeCmd.Flags().StringVarP(&invokePayloadFile, "payload", "p", "", "payload file path.")
	invokeCmd.MarkFlagRequired("payload")

	rootCmd.PersistentFlags().StringVarP(&funcName, "funcName", "n", "", "function name.")
	rootCmd.MarkPersistentFlagRequired("funcName")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "cfgFile", "c", "~/.jf/config.yml", "config file path")
	
	rootCmd.AddCommand(uploadCmd, deleteCmd, invokeCmd)
}

func initConfig() {
	var err error
	jfConfig, err = readConfigFromFile(cfgFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

// Execute root cmd entry
func Execute() {
	rootInit()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}