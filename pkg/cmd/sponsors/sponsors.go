package sponsors

import (
	"github.com/MakeNowJust/heredoc"
	cmdList "github.com/cli/cli/v2/pkg/cmd/sponsors/list"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSponsors(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sponsors <command>",
		Short: "View GitHub sponsors",
		Long:  `Work with GitHub sponsors.`,
		Example: heredoc.Doc(`
			$ gh sponsors list
		`),
	}

	cmd.AddCommand(cmdList.NewCmdList(f, nil))

	return cmd
}
