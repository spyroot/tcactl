package app

import (
	"github.com/spf13/cobra"
	"github.com/spyroot/hestia/cmd/client/main/app/ui"
	"strings"
)

// CmdGetRepos get cnf or vnf instances
func (ctl *TcaCtl) CmdGetRepos() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	// cnf instances
	var _cmd = &cobra.Command{
		Use:     "repos",
		Short:   "Command Returns repositories list.",
		Long:    `Command repositories list.`,
		Example: "tcactl get repos",

		Run: func(cmd *cobra.Command, args []string) {

			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			repos, err := ctl.tca.GetRepos()
			CheckErrLogError(err)

			if repos != nil && len(repos.Items) > 0 {
				if printer, ok := ctl.RepoPrinter[ctl.Printer]; ok {
					printer(repos, ctl.DefaultStyle)
				}
			}
		},
	}

	// output filter , filter specific value from data structure
	_cmd.Flags().StringVar(&_outputFilter, "ofilter", "",
		"Output filter.")

	return _cmd
}
