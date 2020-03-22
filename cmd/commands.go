package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/JointFaaS/Client-go/client"
	"github.com/spf13/cobra"
)

var (
	funcName string
	cfgFile string
	env string
	sourceZip string
	invokePayloadFile string
	enableNative string
	timeout string
	memorySize string
	jfConfig *Config
	hclient *client.Client
)

var uploadCmd = &cobra.Command{
    Use:   "upload",
	Short: "upload function with funcName, zipFile(or sourceDir) and env",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("start upload function", funcName, sourceZip, env)
		file, err := os.Open(sourceZip)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
	
		fileInfo, err := file.Stat()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		var code []byte
		if fileInfo.IsDir() {
			buf := bytes.NewBuffer(nil)
			compressDir(sourceZip, buf)
			code, err = ioutil.ReadAll(buf)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
		} else {
			code, err = ioutil.ReadAll(file)
			if err != nil {
				cmd.PrintErr(err)
				return
			}
		}
		res, err := hclient.FcCreate(&client.FcCreateInput{
			FuncName: funcName,
			Code: code,
			Env: env,
			Timeout: timeout,
			MemorySize: memorySize,
		})
		if err != nil {
			cmd.PrintErr(err)
		} else {
			cmd.Print(string(res.RespBody))
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:	"delete",
	Short: "delete function",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("start delete function", funcName)
		res, err := hclient.FcDelete(&client.FcDeleteInput{
			FuncName: funcName,
		})
		if err != nil {
			cmd.PrintErr(err)
		} else {
			cmd.Print(string(res.RespBody))
		}
	},
}

var invokeCmd = &cobra.Command{
	Use: "invoke",
	Short: "invoke a function with name and payload",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("start invoke function", funcName)
		payloadFile, err := os.Open(invokePayloadFile)
		if err != nil {
			cmd.PrintErr(err.Error())
			return
		}
		payload, err := ioutil.ReadAll(payloadFile)
		if err != nil {
			cmd.PrintErr(err.Error())
			return
		}
		res, err := hclient.FcInvoke(&client.FcInvokeInput{
			FuncName: funcName,
			Args: payload,
			EnableNative: enableNative,
		})
		if err != nil {
			cmd.PrintErr(err)
		} else {
			cmd.Print(string(res.RespBody))
		}
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
	cobra.OnInitialize(clientInit)

	uploadCmd.Flags().StringVarP(&env, "env", "e", "python3", "function env")
	uploadCmd.Flags().StringVarP(&sourceZip, "sourceZip", "z", "", "source code file path")
	uploadCmd.Flags().StringVarP(&timeout, "timeout", "t", "3", "the limitation of function excution time(s)")
	uploadCmd.Flags().StringVarP(&memorySize, "memorySize", "m", "128", "the limitation of function memory(MB)")
	uploadCmd.MarkFlagRequired("sourceZip")

	invokeCmd.Flags().StringVarP(&invokePayloadFile, "payload", "p", "", "payload file path")
	invokeCmd.Flags().StringVarP(&enableNative, "enableNative", "u", "true", "eanble native serverless")
	invokeCmd.MarkFlagRequired("payload")

	rootCmd.PersistentFlags().StringVarP(&funcName, "funcName", "n", "", "function name")
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

func clientInit() {
	var err error
	hclient, err = client.NewClient(client.Config{
		Host: jfConfig.ManagerAddr,
	})
	if err != nil {
		panic(err)
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