package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"code-intelligence.com/cifuzz/pkg/log"
)

// creates an archive out of a given source directory
func Create(sourceDir, targetFile string) (err error) {
	artifactFile, err := os.Create(targetFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer artifactFile.Close()

	gzipWriter := gzip.NewWriter(artifactFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	err = filepath.Walk(sourceDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// if the archive is created in the sourceDir, just ignore it
			if path == targetFile {
				return nil
			}

			// skip everything that is not a regular file
			if !info.Mode().IsRegular() {
				return nil
			}

			if err := addHeaderForFile(sourceDir, path, tarWriter, info); err != nil {
				return err
			}

			if err := addFile(path, tarWriter); err != nil {
				return err
			}

			log.Debugf("Archive (%s): added file %s", targetFile, path)
			return nil
		})

	return
}

func Extract(sourceFile, targetDir string) (err error) {
	archiveFile, err := os.Open(sourceFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer archiveFile.Close()

	gzipReader, err := gzip.NewReader(archiveFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()

		// no more files
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return errors.WithStack(err)
		}

		// skip empty header
		if header == nil {
			continue
		}

		path := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, err := os.Stat(path); os.IsNotExist(err) {
				if err := os.MkdirAll(path, 0755); err != nil {
					return errors.WithStack(err)
				}
			}

		// regular file
		case tar.TypeReg:
			if err := extractFile(path, tarReader, header); err != nil {
				return err
			}
		}
	}
}

func extractFile(path string, tarReader *tar.Reader, header *tar.Header) error {

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.WithStack(err)
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if _, err := io.Copy(file, tarReader); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func addFile(path string, tarWriter *tar.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if _, err := io.Copy(tarWriter, file); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func addHeaderForFile(sourceDir, path string, tarWriter *tar.Writer, info os.FileInfo) error {
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return errors.WithStack(err)
	}

	// get relative file path for tar
	header.Name, err = filepath.Rel(sourceDir, path)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
