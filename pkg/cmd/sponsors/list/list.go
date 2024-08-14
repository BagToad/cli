package list

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/internal/gh"
	"github.com/cli/cli/v2/internal/tableprinter"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	HttpClient func() (*http.Client, error)
	Config     func() (gh.Config, error)
	IO         *iostreams.IOStreams
	User       string
}

func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		HttpClient: f.HttpClient,
		Config:     f.Config,
		IO:         f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:     "list <user>",
		Short:   "List user's sponsors",
		Long:    heredoc.Doc(`List user's sponsors`),
		Example: heredoc.Doc(`$ gh sponsors list`),
		Aliases: []string{"ls"},
		Args:    cmdutil.ExactArgs(1, "expected exactly one argument: the user to list sponsors for"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.User = args[0]
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	return cmd
}

func listRun(opts *ListOptions) error {
	client, err := opts.HttpClient()
	if err != nil {
		return err
	}

	data, err := listSponsors(client, opts)
	if err != nil {
		return err
	}
	if len(data.User.Sponsors.Edges) <= 0 && opts.IO.IsStdoutTTY() {
		fmt.Printf("No sponsors found for %s\n", opts.User)
		return nil
	}

	t := tableprinter.New(opts.IO, tableprinter.WithHeader("SPONSOR"))

	for _, sponsor := range data.User.Sponsors.Edges {
		t.AddField(sponsor.Node.Login)
		t.EndRow()
	}

	err = t.Render()
	if err != nil {
		return err
	}

	return nil
}
