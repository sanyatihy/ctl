package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	KronCmd.AddCommand(execCmd)
}

var execCmd = &cobra.Command{
	Use:   "exec cronjob [flags]",
	Short: "Executes a job now",
	Long: `Executes the specified job or the selected job if none.
Namespace and context flags can be set to help find the right cron job.
If multiple cron job are found, only the first one will be executed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")

		job, err := client.GetDefaultConfigClient().RunCronJob(ctxs, namespace, args[0])

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Successfully started job \"%s\" in context \"%s\" and namespace \"%s\"\n", job.Name, job.Context, job.Namespace)
	},
}