/*
Copyright © 2021 Manish <itzmanish108@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"os/signal"

	"github.com/itzmanish/go-loganalyzer/handler"
	"github.com/itzmanish/go-loganalyzer/internal/logger"
	"github.com/itzmanish/go-loganalyzer/internal/server"
	"github.com/itzmanish/go-loganalyzer/internal/store"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Log analyzer server to collect logs from agent and process it.",

	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			logger.Fatal(err)
		}
		store, err := store.NewFileStore(store.WithDirectory("sample"))
		if err != nil {
			logger.Fatal(err)
		}
		hdl := handler.NewHandler(store)
		s := server.NewServer(server.WithPort(port), server.WithHandler(hdl))
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt)
		go func() {
			<-exit
			if err := s.Stop(); err != nil {
				logger.Fatal(err)
			}
		}()
		if err := s.Start(); err != nil {
			logger.Error(err)
		}
	},
}

func init() {
	appCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().StringP("port", "p", "33555", "log analyzer server port")
}
