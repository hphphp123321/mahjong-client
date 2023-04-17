/*
Copyright © 2023 hphphp123321
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/hphphp123321/mahjong-client/client"
	pb "github.com/hphphp123321/mahjong-common/services/mahjong/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var Client *client.MahjongClient

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mahjong-client",
	Short: "A cli to mahjong server",
	Long:  `A client to mahjong server, you can use it to play mahjong game.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, Mahjong!")
		fmt.Println("set log format: ", logFormat)
		fmt.Println("set log level: ", logLevel)
		fmt.Println("set log output: ", logOutput)
		if logFile != "" {
			fmt.Println("set log file: ", logFile)
		}

		log.Println("Start to connect to server: ", address, ":", port)
		var kacp = keepalive.ClientParameters{
			Time:                time.Duration(timeTicker),
			Timeout:             time.Duration(timeout),
			PermitWithoutStream: true,
		}

		tcpAddr := fmt.Sprintf("%s:%d", address, port)
		log.Debug("Start dial tcpAddr: ", tcpAddr)
		conn, err := grpc.Dial(tcpAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(kacp),
		)
		//defer conn.Close()
		if err != nil {
			log.Fatalf("can not dial: %v", err)
		}

		MahjongClient := pb.NewMahjongClient(conn)
		ctx := context.Background()
		Client = client.NewMahjongClient(ctx, MahjongClient)
		setupLogger()
		fmt.Println()
		logMenu()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	address    string
	port       int
	timeout    int
	timeTicker int
	logFormat  string
	logLevel   string
	logOutput  string
	logFile    string
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mahjong-client.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVar(&address, "address", "127.0.0.1", "server address")
	rootCmd.PersistentFlags().IntVar(&port, "port", 16548, "server port")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 5, "seconds for timeout")
	rootCmd.PersistentFlags().IntVar(&timeTicker, "timeTicker", 10, "seconds for time ticker")
	rootCmd.PersistentFlags().StringVar(&logFormat, "logFormat", "text", "log format(json or text)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "debug", "log level(debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVar(&logOutput, "logOutput", "stdout", "log output(stdout or stderr)")
	rootCmd.PersistentFlags().StringVar(&logFile, "logFile", "", "log file path")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
