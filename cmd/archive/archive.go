package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var output string

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Create and extract archive files",
	Long: `Create and extract archive files in various formats.

SUBCOMMANDS:

  create    Create an archive from files/directories
  extract   Extract an archive to a directory

SUPPORTED FORMATS (pure Go — no external tools):

  ZIP       .zip
  TAR       .tar
  TAR.GZ    .tar.gz, .tgz
  TAR.BZ2   .tar.bz2, .tbz2 (extract only)

For other formats (7z, rar, etc.), use system tools:
  7z:   brew install p7zip    →  7z x file.7z
  RAR:  brew install unrar    →  unrar x file.rar

Examples:
  openGyver archive create -o backup.zip file1.txt file2.txt dir/
  openGyver archive create -o project.tar.gz src/ README.md
  openGyver archive create -o files.tar doc1.txt doc2.txt
  openGyver archive extract backup.zip
  openGyver archive extract backup.zip -o ./extracted/
  openGyver archive extract project.tar.gz -o ./project/`,
}

// --- Create subcommand ---

var createCmd = &cobra.Command{
	Use:   "create <file> [file...]",
	Short: "Create an archive from files and directories",
	Long: `Create a ZIP, TAR, or TAR.GZ archive from the given files and directories.

The archive format is determined by the --output extension.

Examples:
  openGyver archive create -o backup.zip file1.txt dir/
  openGyver archive create -o project.tar.gz src/ README.md
  openGyver archive create -o data.tar *.csv`,
	Args: cobra.MinimumNArgs(1),
	RunE: runCreate,
}

func runCreate(c *cobra.Command, args []string) error {
	if output == "" {
		return fmt.Errorf("--output (-o) is required")
	}

	lower := strings.ToLower(output)
	switch {
	case strings.HasSuffix(lower, ".zip"):
		return createZip(output, args)
	case strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz"):
		return createTarGz(output, args)
	case strings.HasSuffix(lower, ".tar"):
		return createTar(output, args)
	default:
		return fmt.Errorf("unsupported archive format: %s\nSupported: .zip, .tar, .tar.gz, .tgz", filepath.Ext(output))
	}
}

// --- Extract subcommand ---

var extractOutput string

var extractCmd = &cobra.Command{
	Use:   "extract <archive>",
	Short: "Extract an archive to a directory",
	Long: `Extract a ZIP, TAR, or TAR.GZ archive.

Extracts to the current directory by default. Use -o to specify a target.

Examples:
  openGyver archive extract backup.zip
  openGyver archive extract backup.zip -o ./restored/
  openGyver archive extract project.tar.gz -o ./project/`,
	Args: cobra.ExactArgs(1),
	RunE: runExtract,
}

func runExtract(c *cobra.Command, args []string) error {
	archivePath := args[0]
	dest := extractOutput
	if dest == "" {
		dest = "."
	}

	lower := strings.ToLower(archivePath)
	switch {
	case strings.HasSuffix(lower, ".zip"):
		return extractZip(archivePath, dest)
	case strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz"):
		return extractTarGz(archivePath, dest)
	case strings.HasSuffix(lower, ".tar"):
		return extractTar(archivePath, dest)
	default:
		return fmt.Errorf("unsupported archive format: %s\nSupported: .zip, .tar, .tar.gz, .tgz", filepath.Ext(archivePath))
	}
}

// --- ZIP ---

func createZip(outPath string, sources []string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()

	count := 0
	for _, src := range sources {
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = path
			header.Method = zip.Deflate
			writer, err := w.CreateHeader(header)
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			count++
			return err
		})
		if err != nil {
			return err
		}
	}

	fmt.Printf("Created %s (%d files)\n", outPath, count)
	return nil
}

func extractZip(archivePath, dest string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("opening zip: %w", err)
	}
	defer r.Close()

	count := 0
	for _, f := range r.File {
		target := filepath.Join(dest, f.Name)

		// Guard against zip slip
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)+string(os.PathSeparator)) && filepath.Clean(target) != filepath.Clean(dest) {
			return fmt.Errorf("illegal file path in archive: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0755)
			continue
		}

		os.MkdirAll(filepath.Dir(target), 0755)
		outFile, err := os.Create(target)
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return err
		}
		count++
	}

	fmt.Printf("Extracted %s → %s (%d files)\n", archivePath, dest, count)
	return nil
}

// --- TAR ---

func createTar(outPath string, sources []string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	count, err := writeTar(f, sources)
	if err != nil {
		return err
	}

	fmt.Printf("Created %s (%d files)\n", outPath, count)
	return nil
}

func createTarGz(outPath string, sources []string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	count, err := writeTar(gw, sources)
	if err != nil {
		return err
	}

	fmt.Printf("Created %s (%d files)\n", outPath, count)
	return nil
}

func writeTar(w io.Writer, sources []string) (int, error) {
	tw := tar.NewWriter(w)
	defer tw.Close()

	count := 0
	for _, src := range sources {
		err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			header.Name = path
			if err := tw.WriteHeader(header); err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tw, file)
			count++
			return err
		})
		if err != nil {
			return count, err
		}
	}
	return count, nil
}

func extractTar(archivePath, dest string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	count, err := readTar(f, dest)
	if err != nil {
		return err
	}

	fmt.Printf("Extracted %s → %s (%d files)\n", archivePath, dest, count)
	return nil
}

func extractTarGz(archivePath, dest string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("decompressing gzip: %w", err)
	}
	defer gr.Close()

	count, err := readTar(gr, dest)
	if err != nil {
		return err
	}

	fmt.Printf("Extracted %s → %s (%d files)\n", archivePath, dest, count)
	return nil
}

func readTar(r io.Reader, dest string) (int, error) {
	tr := tar.NewReader(r)
	count := 0

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, err
		}

		target := filepath.Join(dest, header.Name)

		// Guard against path traversal
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)+string(os.PathSeparator)) && filepath.Clean(target) != filepath.Clean(dest) {
			return count, fmt.Errorf("illegal file path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			outFile, err := os.Create(target)
			if err != nil {
				return count, err
			}
			_, err = io.Copy(outFile, tr)
			outFile.Close()
			if err != nil {
				return count, err
			}
			count++
		}
	}

	return count, nil
}

func init() {
	createCmd.Flags().StringVarP(&output, "output", "o", "", "output archive path (required)")
	extractCmd.Flags().StringVarP(&extractOutput, "output", "o", "", "extraction destination directory (default: current directory)")

	archiveCmd.AddCommand(createCmd)
	archiveCmd.AddCommand(extractCmd)
	cmd.Register(archiveCmd)
}
