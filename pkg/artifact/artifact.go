package artifact

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"code-intelligence.com/cifuzz/pkg/artifact/archive"
	"code-intelligence.com/cifuzz/util/fileutil"
)

func Create(targetPath string, metadata *Metadata) error {
	yaml, err := metadata.ToYaml()
	if err != nil {
		return err
	}

	// create tmp dir to place all files needed for the artifact
	bundleDir, err := os.MkdirTemp("", "cifuzz-bundle")
	if err != nil {
		return errors.WithStack(err)
	}
	defer fileutil.Cleanup(bundleDir)

	// add metadata file
	filename := filepath.Join(bundleDir, MetadataFileName)
	os.WriteFile(filename, yaml, 0644)

	// copy all fuzzer
	for _, fuzzer := range metadata.Fuzzers {
		baseDir := filepath.Join(bundleDir, filepath.Dir(fuzzer.Path))
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return errors.WithStack(err)
		}
		if err := fileutil.CopyFile(fuzzer.Path, filepath.Join(baseDir, fuzzer.Target), 0755); err != nil {
			return err
		}
	}

	if err := archive.Create(bundleDir, targetPath); err != nil {
		return err
	}

	return nil
}
