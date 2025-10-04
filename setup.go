package donutdns

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
)

var pluginLogger = log.NewWithPlugin(PluginName)

func init() { plugin.Register(PluginName, Setup) }

// Setup will parse plugin config and register the donutdns plugin
// with the CoreDNS core server.
//
// todo: test with TestController
func Setup(c *caddy.Controller) error {
	// reconstruct the parts of CoreConfig for initializing the allow/block lists
	cc := new(CoreConfig)
	cc.Forward = new(Forward)

	for c.Next() {
		_ = c.RemainingArgs()
		for c.NextBlock() {
			switch c.Val() {
			case "defaults":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.NoDefaults = c.Val() == "false"

			case "allow_dir":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.AllowDir = c.Val()

			case "allowsuffix_dir":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.AllowSuffixDir = c.Val()

			case "block_dir":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.BlockDir = c.Val()

			case "suffix_dir":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.SuffixDir = c.Val()

			case "allow_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.AllowFile = c.Val()

			case "allowsuffix_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.AllowSuffixFile = c.Val()

			case "block_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.BlockFile = c.Val()

			case "suffix_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.SuffixFile = c.Val()

			case "allow":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Allows = append(cc.Allows, c.Val())

			case "allowsuffix":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.AllowSuffix = append(cc.AllowSuffix, c.Val())

			case "block":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Blocks = append(cc.Blocks, c.Val())

			case "suffix":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Suffix = append(cc.Suffix, c.Val())

			case "upstream_1":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Forward.Addresses = append(cc.Forward.Addresses, c.Val())

			case "upstream_2":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Forward.Addresses = append(cc.Forward.Addresses, c.Val())

			case "forward_server_name":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Forward.ServerName = c.Val()

			case "forward_max_fails":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cc.Forward.MaxFails, _ = strconv.Atoi(c.Val())
			}
		}
	}

	sets := NewSets(pluginLogger, cc)
	allow, allowsuffix, block, suffix := sets.Size()
	pluginLogger.Infof("domains on explicit allow-list(s): %d", allow)
	pluginLogger.Infof("domains on suffixes allow-list(s): %d", allowsuffix)
	pluginLogger.Infof("domains on explicit block-list(s): %d", block)
	pluginLogger.Infof("domains on suffixes block-list(s): %d", suffix)
	pluginLogger.Infof("forward upstreams: %v", cc.Forward.Addresses)
	pluginLogger.Infof("forward name: %s", cc.Forward.ServerName)

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dd := DonutDNS{sets: sets}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		dd.Next = next
		return dd
	})

	// Plugin loaded okay.
	return nil
}
