/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alexanderjophus/kie/svc/cleaning/pkg"
	"github.com/spf13/cobra"
)

var (
	inDir  string
	outDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cleaning",
	Short: "Small service for transforming json data into csv data",
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(inDir, outDir); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func run(inDir, outDir string) error {
	infs := os.DirFS(inDir)
	if err := fs.WalkDir(infs, "api/v1/people", func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".raw" {
			return nil
		}
		fmt.Printf("processing %s\n", path)

		f, err := infs.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer f.Close()

		outPath := filepath.Join(outDir, path)
		outPath = changeFileExtension(outPath, "csv")
		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
				return fmt.Errorf("error creating out dir: %w", err)
			}
		}
		fOut, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer func() {
			if err := fOut.Close(); err != nil {
				fmt.Printf("error closing file: %s\n", err)
			}
		}()

		if err := pkg.JSONToCSV(f, fOut); err != nil {
			return fmt.Errorf("error converting json to csv: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error walking dir: %w", err)
	}

	return nil
}

func changeFileExtension(filePath string, newExtension string) string {
	ext := filepath.Ext(filePath)
	return filePath[:len(filePath)-len(ext)] + "." + newExtension
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
	rootCmd.Flags().StringVar(&inDir, "in-dir", "/pfs/statsapi", "Input directory")
	rootCmd.Flags().StringVar(&outDir, "out-dir", "/pfs/out", "Output directory")
}
