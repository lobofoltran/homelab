package cmd

import (
	"fmt"

	"github.com/lobofoltran/homelab/apps/cli/internal/api"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Testa a comunicação com o daemon",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := api.Ping()
		if err != nil {
			fmt.Println("Erro:", err)
			return
		}
		fmt.Println("Resposta do daemon:", resp)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
