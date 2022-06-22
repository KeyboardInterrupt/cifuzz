package artifact

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"code-intelligence.com/cifuzz/internal/config"
)

// MetadataFileName is the name of the meta information yaml file within an artifacts archive
const MetadataFileName = "cifuzz.yaml"

// Metadata defines meta information for artifacts contained within a fuzzing archive
type Metadata struct {
	*RunEnvironment `yaml:"run_environment"`
	CodeRevision    *CodeRevision `yaml:"code_revision,omitempty"`
	Fuzzers         []*Fuzzer     `yaml:"fuzzers"`
}

// Fuzzer defines run time information related to the fuzzing artifact
type Fuzzer struct {
	Target string `yarn:"target"`
	Path   string `yarn:"path"`
	Engine string `yarn:"engine"`
}

// RunEnvironment defines run environment specific information
type RunEnvironment struct {
	// The docker image and tag to be used: eg. debian:stable
	Docker string
}

type CodeRevision struct {
	Git *GitRevision `yaml:"git,omitempty"`
}

type GitRevision struct {
	Commit string `yaml:"commit,omitempty"`
	Branch string `yaml:"branch,omitempty"`
}

// ToYaml converts a meta information artifacts struct to a human readable yaml representation
func (a *Metadata) ToYaml() ([]byte, error) {
	out, err := yaml.Marshal(a)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal configuration to yaml")
	}

	return out, nil
}

// FromYaml  fills an artifacts meta information struct from the bytes of a yaml representation
func (a *Metadata) FromYaml(data []byte) error {
	err := yaml.Unmarshal(data, a)
	if err != nil {
		return errors.Wrap(err, "failed to marshal configuration to yaml")
	}

	return nil
}

func NewMetadata(cfg *config.Config) (*Metadata, error) {
	codeRevision, err := getCodeRevision()
	if err != nil {
		return nil, err
	}

	var fuzzers []*Fuzzer
	if cfg.BuildSystem == config.BuildSystemCMake {
		fuzzers, err = getCMakeFuzzers()
		if err != nil {
			return nil, err
		}
	}

	return &Metadata{
		// TODO replace Docker image with a "real" value
		RunEnvironment: &RunEnvironment{Docker: "foo/bar"},
		CodeRevision:   codeRevision,
		Fuzzers:        fuzzers,
	}, nil
}
