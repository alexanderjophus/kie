/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alexanderjophus/kie/gen/ingestion/v1/ingestionv1connect"
	"github.com/alexanderjophus/kie/svc/ingestion/pkg"
	"github.com/pachyderm/pachyderm/v2/src/client"
	"github.com/pachyderm/pachyderm/v2/src/pfs"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ingestion",
	Short: "An ingestion service for hockey data",
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func run() error {
	client, err := client.NewInCluster()
	if err != nil {
		return fmt.Errorf("unable to get client: %w", err)
	}

	if err := client.UpdateProjectRepo(pfs.DefaultProjectName, pkg.RepoName); err != nil {
		return fmt.Errorf("cannot create repo %s: %w", pkg.RepoName, err)
	}

	s := pkg.NewServer(&pkg.PachdRepo{Client: client})

	mux := http.NewServeMux()
	path, handler := ingestionv1connect.NewIngestionServiceHandler(s)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
