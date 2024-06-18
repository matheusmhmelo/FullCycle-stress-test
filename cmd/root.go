/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/matheusmhmelo/FullCycle-stress-test/internal/stresstest"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Full Cycle Project Stress Test",
	Short: "Stress Test system with concurrency",
	Long: `The project consists in a stress test application that allow
users to do tests to HTTP endpoints using concurrency.`,
	Run: stresstest.ExecuteTests,
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
	rootCmd.Flags().StringP("url", "u", "", "URL do serviço a ser testado.")
	rootCmd.MarkFlagRequired("url")
	rootCmd.Flags().IntP("requests", "r", 0, "Número total de requests.")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.Flags().IntP("concurrency", "c", 0, "Número de chamadas simultâneas.")
	rootCmd.MarkFlagRequired("concurrency")
}
