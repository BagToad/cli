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

type iprompter interface {
	Input(prompt, defaultValue string) (string, error)
}

type ListOptions struct {
	HttpClient func() (*http.Client, error)
	Config     func() (gh.Config, error)
	IO         *iostreams.IOStreams
	User       string
	Prompter   iprompter
	Exporter   cmdutil.Exporter
}

func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		HttpClient: f.HttpClient,
		Config:     f.Config,
		IO:         f.IOStreams,
		Prompter:   f.Prompter,
	}

	cmd := &cobra.Command{
		Use:     "list <user>",
		Short:   "List user's sponsors",
		Long:    heredoc.Doc(`List user's sponsors`),
		Example: heredoc.Doc(`$ gh sponsors list`),
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.User = args[0]
			}

			if !opts.IO.CanPrompt() && opts.User == "" {
				return cmdutil.FlagErrorf("user required when not running interactively")
			}

			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmdutil.AddJSONFlags(cmd, &opts.Exporter, []string{"login"})

	return cmd
}

func listRun(opts *ListOptions) error {
	client, err := opts.HttpClient()
	if err != nil {
		return err
	}

	if opts.User == "" {
		opts.User, err = opts.Prompter.Input("Which user do you want to target?", "")
		if err != nil {
			return err
		}
	}

	data, err := listSponsors(client, opts)
	if err != nil {
		return err
	}

	if opts.Exporter != nil {
		return opts.Exporter.Write(opts.IO, data)
	}

	if len(data.Sponsors) <= 0 && opts.IO.IsStdoutTTY() {
		fmt.Printf("No sponsors found for %s\n", opts.User)
		return nil
	}

	t := tableprinter.New(opts.IO, tableprinter.WithHeader("SPONSOR"))

	for _, sponsor := range data.Sponsors {
		t.AddField(sponsor)
		t.EndRow()
	}

	err = t.Render()
	if err != nil {
		return err
	}

	return nil
}
