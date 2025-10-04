package donutdns

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-set"
	"github.com/shoenig/ignore"
)

// Sets enables efficient look-ups of whether a domain should be allowable or blocked.
type Sets struct {
	allow       *set.Set[string]
	allowsuffix *set.Set[string]
	block       *set.Set[string]
	suffix      *set.Set[string]
}

// NewSets returns a Sets pre-filled according to cc.
func NewSets(logger Logger, cc *CoreConfig) *Sets {
	allow := set.New[string](100)
	allowsuffix := set.New[string](100)
	block := set.New[string](100)
	suffix := set.New[string](100)

	// initialize defaults if enabled
	if !cc.NoDefaults {
		defaults(cc.Forward, block, logger)
	}

	// insert individual custom allowable domains
	allow.InsertSlice(cc.Allows)

	// insert file of custom allowable domains
	customFile(cc.AllowFile, allow)

	// insert each file of custom allowable domains
	customDir(cc.AllowDir, allow)

	// insert individual allowable domain suffixes
	allowsuffix.InsertSlice(cc.AllowSuffix)

	// insert file of custom allowable domain suffixes
	customFile(cc.AllowSuffixFile, allowsuffix)

	// insert each file of allowable custom domain suffixes
	customDir(cc.AllowSuffixDir, allowsuffix)

	// insert individual custom block domains
	block.InsertSlice(cc.Blocks)

	// insert file of custom block domains
	customFile(cc.BlockFile, block)

	// insert each file of custom block domains
	customDir(cc.BlockDir, block)

	// insert individual domain sufix block
	suffix.InsertSlice(cc.Suffix)

	// insert file of custom domain suffix blocks
	customFile(cc.SuffixFile, suffix)

	// insert each file of custom domain suffix blocks
	customDir(cc.SuffixDir, suffix)

	return &Sets{
		allow:       allow,
		allowsuffix: allowsuffix,
		block:       block,
		suffix:      suffix,
	}
}

// Size returns the number of items in the allow, block, suffix sets.
func (s *Sets) Size() (int, int, int, int) {
	allow := s.allow.Size()
	allowsuffix := s.allowsuffix.Size()
	block := s.block.Size()
	suffix := s.suffix.Size()
	return allow, allowsuffix, block, suffix
}

// Allow indicates whether domain is on the explicit allow-list.
func (s *Sets) Allow(domain string) bool {
	return s.allow.Contains(domain)
}

func (s *Sets) AllowBySuffix(domain string) bool {
	if s.allowsuffix.Size() == 0 {
		return false
	}

	domain = strings.Trim(domain, ".")
	if domain == "" {
		return false
	}

	if s.allowsuffix.Contains(domain) {
		return true
	}

	idx := strings.Index(domain, ".")
	if idx <= 0 {
		return false
	}

	return s.AllowBySuffix(domain[idx+1:])
}

// BlockByMatch indicates whether domain is on the explicit block-list.
func (s *Sets) BlockByMatch(domain string) bool {
	return s.block.Contains(domain)
}

// BlockBySuffix indicates whether domain is on the suffix block-list.
func (s *Sets) BlockBySuffix(domain string) bool {
	if s.suffix.Size() == 0 {
		return false
	}

	domain = strings.Trim(domain, ".")
	if domain == "" {
		return false
	}

	if s.suffix.Contains(domain) {
		return true
	}

	idx := strings.Index(domain, ".")
	if idx <= 0 {
		return false
	}

	return s.BlockBySuffix(domain[idx+1:])
}

func defaults(fwd *Forward, set *set.Set[string], logger Logger) {
	d := NewDownloader(fwd, logger)
	s, err := d.Download(Defaults())
	if err != nil {
		panic(err)
	}
	set.InsertSet(s)
}

func customFile(filename string, set *set.Set[string]) {
	if filename == "" {
		return // nothing to do
	}

	// for now, everything uses the generic domain extractor
	ex := NewExtractor(Generic)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer ignore.Close(f)

	s, err := ex.Extract(f)
	if err != nil {
		panic(err)
	}
	set.InsertSet(s)
}

func customDir(dirname string, set *set.Set[string]) {
	if dirname == "" {
		return // nothing to do
	}

	files, err := os.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := filepath.Join(dirname, file.Name())
		customFile(filename, set)
	}
}
