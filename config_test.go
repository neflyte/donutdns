package donutdns

import (
	"strings"
	"testing"

	"github.com/shoenig/extractors/env"
	"github.com/shoenig/test/must"
)

func TestCoreConfig_Generate(t *testing.T) {
	cc := CoreConfig{
		Port:       1053,
		NoDebug:    true,
		NoLog:      true,
		Allows:     []string{"example.com", "pets.com"},
		AllowFile:  "/etc/allow.list",
		AllowDir:   "/etc/allows",
		Blocks:     []string{"facebook.com", "instagram.com"},
		BlockFile:  "/etc/block.list",
		BlockDir:   "/etc/blocks",
		Suffix:     []string{"fb.com", "twitter.com"},
		SuffixFile: "/etc/suffix.list",
		SuffixDir:  "/etc/suffixes",
		Forward: &Forward{
			Addresses:  []string{"1.1.1.1", "1.0.0.1"},
			ServerName: "cloudflare-dns.com",
			MaxFails:   2,
		},
	}

	result := cc.Generate()
	must.Eq(t, noWhitespace(`
.:1053 {
  donutdns {
    defaults true
    allow_file /etc/allow.list
    block_file /etc/block.list
    suffix_file /etc/suffix.list

    allow_dir /etc/allows
    block_dir /etc/blocks
    suffix_dir /etc/suffixes

    allow example.com
    allow pets.com
    
    block facebook.com
    block instagram.com

    suffix fb.com
    suffix twitter.com

    upstream_1 1.1.1.1
    upstream_2 1.0.0.1
    forward_server_name cloudflare-dns.com

  }
  forward . 1.1.1.1 1.0.0.1 {
    tls_servername cloudflare-dns.com
    max_fails 2
  }
}
`), noWhitespace(result))
}

func TestCoreConfig_Generate_less(t *testing.T) {
	cc := CoreConfig{
		Port:       1054,
		NoDebug:    false,
		NoLog:      false,
		Allows:     nil,
		Blocks:     nil,
		NoDefaults: true,
		Forward: &Forward{
			Addresses:  []string{"8.8.8.8"},
			ServerName: "google.dns",
			MaxFails:   2,
		},
	}

	result := cc.Generate()
	must.Eq(t, noWhitespace(`
.:1054 {
  debug
  log
  donutdns {
    defaults false
    upstream_1 8.8.8.8
    forward_server_name google.dns
  }
  forward . 8.8.8.8 {
    tls_servername google.dns
    max_fails 2
  }
}
`), noWhitespace(result))
}

func noWhitespace(s string) string {
	a := strings.ReplaceAll(s, " ", "")
	b := strings.ReplaceAll(a, "\n", "")
	return b
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("DONUT_DNS_PORT", "1234")
	t.Setenv("DONUT_DNS_NO_DEBUG", "1")
	t.Setenv("DONUT_DNS_NO_LOG", "1")
	t.Setenv("DONUT_DNS_ALLOW", "example.com,pets.com")
	t.Setenv("DONUT_DNS_ALLOW_FILE", "/etc/allow.list")
	t.Setenv("DONUT_DNS_ALLOW_DIR", "/etc/allows")
	t.Setenv("DONUT_DNS_BLOCK", "facebook.com,reddit.com")
	t.Setenv("DONUT_DNS_BLOCK_FILE", "/etc/block.list")
	t.Setenv("DONUT_DNS_BLOCK_DIR", "/etc/blocks")
	t.Setenv("DONUT_DNS_SUFFIX", "fb.com,twitter.com")
	t.Setenv("DONUT_DNS_SUFFIX_FILE", "/etc/suffix.list")
	t.Setenv("DONUT_DNS_SUFFIX_DIR", "/etc/suffixes")
	t.Setenv("DONUT_DNS_NO_DEFAULTS", "")
	t.Setenv("DONUT_DNS_UPSTREAM_1", "8.8.8.8")
	t.Setenv("DONUT_DNS_UPSTREAM_2", "8.8.4.4")
	t.Setenv("DONUT_DNS_UPSTREAM_NAME", "dns.google")
	t.Setenv("DONUT_DNS_UPSTREAM_MAX_FAILS", "5")

	cc := ConfigFromEnv(env.OS)
	must.Eq(t, &CoreConfig{
		Port:       1234,
		NoDebug:    true,
		NoLog:      true,
		Allows:     []string{"example.com", "pets.com"},
		AllowFile:  "/etc/allow.list",
		AllowDir:   "/etc/allows",
		Blocks:     []string{"facebook.com", "reddit.com"},
		BlockFile:  "/etc/block.list",
		BlockDir:   "/etc/blocks",
		Suffix:     []string{"fb.com", "twitter.com"},
		SuffixFile: "/etc/suffix.list",
		SuffixDir:  "/etc/suffixes",
		NoDefaults: false,
		Forward: &Forward{
			Addresses:  []string{"8.8.8.8", "8.8.4.4"},
			ServerName: "dns.google",
			MaxFails:   5,
		},
	}, cc)
}

func TestConfigFromEnv_2(t *testing.T) {
	t.Setenv("DONUT_DNS_PORT", "1234")
	t.Setenv("DONUT_DNS_NO_DEBUG", "0")
	t.Setenv("DONUT_DNS_NO_LOG", "true")
	t.Setenv("DONUT_DNS_ALLOW", "")
	t.Setenv("DONUT_DNS_ALLOW_FILE", "")
	t.Setenv("DONUT_DNS_ALLOW_DIR", "")
	t.Setenv("DONUT_DNS_BLOCK", "facebook.com")
	t.Setenv("DONUT_DNS_BLOCK_FILE", "")
	t.Setenv("DONUT_DNS_BLOCK_DIR", "")
	t.Setenv("DONUT_DNS_SUFFIX", "")
	t.Setenv("DONUT_DNS_SUFFIX_FILE", "")
	t.Setenv("DONUT_DNS_SUFFIX_DIR", "")
	t.Setenv("DONUT_DNS_NO_DEFAULTS", "true")
	t.Setenv("DONUT_DNS_UPSTREAM_1", "8.8.8.8")
	t.Setenv("DONUT_DNS_UPSTREAM_2", "")
	t.Setenv("DONUT_DNS_UPSTREAM_NAME", "dns.google")
	t.Setenv("DONUT_DNS_UPSTREAM_MAX_FAILS", "4")

	cc := ConfigFromEnv(env.OS)
	must.Eq(t, &CoreConfig{
		Port:       1234,
		NoDebug:    false,
		NoLog:      true,
		Allows:     nil,
		Blocks:     []string{"facebook.com"},
		BlockFile:  "",
		NoDefaults: true,
		Forward: &Forward{
			Addresses:  []string{"8.8.8.8"},
			ServerName: "dns.google",
			MaxFails:   4,
		},
	}, cc)
}
