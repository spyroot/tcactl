package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/app/main/cmds/templates"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/client/request"
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
		Aliases: []string{"repo", "rp"},
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

// CmdCreateExtension - Command register new extension
func (ctl *TcaCtl) CmdCreateExtension() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:   "extension [spec file]",
		Short: "Command create a new extension in TCA.",
		Long: templates.LongDesc(`

Command attaches cloud provider to TCA.`),

		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			spec, err := request.ExtensionSpecsFromFile(args[0])
			CheckErrLogError(err)

			eid, err := ctl.tca.CreateExtension(spec)
			CheckErrLogError(err)

			fmt.Printf("Extention type %s registered extention id %s\n", spec.Type, eid)
		},
	}

	return _cmd
}

// CmdDeleteExtension - Command delete extension
func (ctl *TcaCtl) CmdDeleteExtension() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle
		_outputFilter   string
	)

	var _cmd = &cobra.Command{
		Use:   "extension [name or id]",
		Short: "Command deletes an extension in TCA.",
		Long: templates.LongDesc(`

Command delete extension from TCA.`),

		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup("output").Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			_, err := ctl.tca.DeleteExtension(args[0])
			CheckErrLogError(err)

			fmt.Printf("Extention %s deleted\n", args[0])
		},
	}

	return _cmd
}
