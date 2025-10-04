package donutdns

import (
	"context"
	"flag"
	"strings"

	"github.com/google/subcommands"
	"github.com/shoenig/extractors/env"
)

const (
	checkCmdName = "check"
)

type CheckCmd struct {
	quiet    bool
	defaults bool
}

func NewCheckCmd() subcommands.Command {
	return new(CheckCmd)
}
func (cc *CheckCmd) Name() string {
	return checkCmdName
}

func (cc *CheckCmd) Synopsis() string {
	return "Check whether a domain will be blocked."
}

func (cc *CheckCmd) Usage() string {
	return strings.TrimPrefix(`
check [-quiet] [-defaults] <domain>
Check whether domain will be blocked.
`, "\n")
}

func (cc *CheckCmd) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&cc.quiet, "quiet", false, "silence verbose debug output")
	fs.BoolVar(&cc.defaults, "defaults", false, "also check against default block lists")
}

func (cc *CheckCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	logger := new(CLI)

	args := f.Args()
	if len(args) == 0 {
		logger.Errorf("must specify domain to check command")
		return subcommands.ExitUsageError
	}

	if err := cc.execute(logger, args[0]); err != nil {
		logger.Errorf("failure: %v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (cc *CheckCmd) execute(output *CLI, domain string) error {
	cfg := ConfigFromEnv(env.OS)
	ApplyDefaults(cfg)
	cfg.NoDefaults = !cc.defaults

	if !cc.quiet {
		cfg.Log(output)
	}

	sets := NewSets(output, cfg)
	switch {
	case sets.Allow(domain):
		output.Infof("domain %q on explicit allow list", domain)
	case sets.AllowBySuffix(domain):
		output.Infof("domain %q on suffix allow list", domain)
	case sets.BlockByMatch(domain):
		output.Infof("domain %q on explicit block list", domain)
	case sets.BlockBySuffix(domain):
		output.Infof("domain %q on suffix block list", domain)
	default:
		output.Infof("domain %q is implicitly allowable", domain)
	}
	return nil
}
