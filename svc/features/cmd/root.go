/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alexanderjophus/kie/svc/features/pkg"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "features",
	Short: "A service for parsing features in csv's and encoding them ready for ML",
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var (
	inDir  = "./cleaning/"
	outDir = "./out/"
)

func run() error {
	infs := os.DirFS(inDir)
	if err := fs.WalkDir(infs, ".", func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking dir: %w", err)
		}
		if filepath.Ext(path) != ".csv" {
			return nil
		}
		f, err := infs.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer f.Close()

		outPath := filepath.Join(outDir, path)
		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
				return fmt.Errorf("error creating out dir: %w", err)
			}
		}
		fOut, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer fOut.Close()

		if err := pkg.FeatureExtract(f, fOut); err != nil {
			return fmt.Errorf("error parsing file: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error walking dir: %w", err)
	}

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
