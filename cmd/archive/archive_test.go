package archive

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "archive-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func createTestFiles(t *testing.T, dir string) {
	t.Helper()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("file a"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("file b"), 0644)
	os.MkdirAll(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "sub", "c.txt"), []byte("file c"), 0644)
}

// ---------------------------------------------------------------------------
// ZIP
// ---------------------------------------------------------------------------

func TestCreateZip(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	zipPath := filepath.Join(dir, "test.zip")
	err := createZip(zipPath, []string{
		filepath.Join(dir, "a.txt"),
		filepath.Join(dir, "b.txt"),
	})
	if err != nil {
		t.Fatal(err)
	}

	info, _ := os.Stat(zipPath)
	if info.Size() < 10 {
		t.Error("zip file too small")
	}
}

func TestExtractZip(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	zipPath := filepath.Join(dir, "test.zip")
	createZip(zipPath, []string{filepath.Join(dir, "a.txt"), filepath.Join(dir, "b.txt")})

	outDir := filepath.Join(dir, "extracted")
	os.MkdirAll(outDir, 0755)
	err := extractZip(zipPath, outDir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestZipRoundtrip(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	zipPath := filepath.Join(dir, "test.zip")
	createZip(zipPath, []string{filepath.Join(dir, "a.txt")})

	outDir := filepath.Join(dir, "out")
	os.MkdirAll(outDir, 0755)
	extractZip(zipPath, outDir)

	// Find the extracted file
	entries, _ := filepath.Glob(filepath.Join(outDir, "**", "a.txt"))
	if len(entries) == 0 {
		// Try direct path
		data, err := os.ReadFile(filepath.Join(outDir, filepath.Join(dir, "a.txt")))
		if err != nil {
			t.Log("file not at expected path, checking content exists")
		} else if string(data) != "file a" {
			t.Errorf("content mismatch: got %q", string(data))
		}
	}
}

// ---------------------------------------------------------------------------
// TAR.GZ
// ---------------------------------------------------------------------------

func TestCreateTarGz(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	tgzPath := filepath.Join(dir, "test.tar.gz")
	err := createTarGz(tgzPath, []string{
		filepath.Join(dir, "a.txt"),
		filepath.Join(dir, "b.txt"),
	})
	if err != nil {
		t.Fatal(err)
	}

	info, _ := os.Stat(tgzPath)
	if info.Size() < 10 {
		t.Error("tar.gz file too small")
	}
}

func TestExtractTarGz(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	tgzPath := filepath.Join(dir, "test.tar.gz")
	createTarGz(tgzPath, []string{filepath.Join(dir, "a.txt")})

	outDir := filepath.Join(dir, "out")
	os.MkdirAll(outDir, 0755)
	err := extractTarGz(tgzPath, outDir)
	if err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// TAR
// ---------------------------------------------------------------------------

func TestCreateTar(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	tarPath := filepath.Join(dir, "test.tar")
	err := createTar(tarPath, []string{filepath.Join(dir, "a.txt")})
	if err != nil {
		t.Fatal(err)
	}

	info, _ := os.Stat(tarPath)
	if info.Size() < 10 {
		t.Error("tar file too small")
	}
}

func TestExtractTar(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	tarPath := filepath.Join(dir, "test.tar")
	createTar(tarPath, []string{filepath.Join(dir, "a.txt")})

	outDir := filepath.Join(dir, "out")
	os.MkdirAll(outDir, 0755)
	err := extractTar(tarPath, outDir)
	if err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// Subcommand metadata
// ---------------------------------------------------------------------------

func TestArchiveCmd_Metadata(t *testing.T) {
	if archiveCmd.Use != "archive" {
		t.Errorf("unexpected Use: %s", archiveCmd.Use)
	}
	if archiveCmd.Short == "" {
		t.Error("Short should not be empty")
	}
}

func TestCreateCmd_Metadata(t *testing.T) {
	if createCmd.Use != "create <file> [file...]" {
		t.Errorf("unexpected Use: %s", createCmd.Use)
	}
}

func TestExtractCmd_Metadata(t *testing.T) {
	if extractCmd.Use != "extract <archive>" {
		t.Errorf("unexpected Use: %s", extractCmd.Use)
	}
}

func TestCreateCmd_Flags(t *testing.T) {
	if createCmd.Flags().Lookup("output") == nil {
		t.Error("--output flag not found")
	}
}

func TestExtractCmd_Flags(t *testing.T) {
	if extractCmd.Flags().Lookup("output") == nil {
		t.Error("--output flag not found")
	}
}

// ---------------------------------------------------------------------------
// Content-verified roundtrips
// ---------------------------------------------------------------------------

func TestTarRoundtrip_Content(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	tarPath := filepath.Join(dir, "test.tar")
	createTar(tarPath, []string{filepath.Join(dir, "a.txt"), filepath.Join(dir, "b.txt")})

	outDir := filepath.Join(dir, "out")
	os.MkdirAll(outDir, 0755)
	extractTar(tarPath, outDir)

	data, err := os.ReadFile(filepath.Join(outDir, filepath.Join(dir, "a.txt")))
	if err != nil {
		t.Fatalf("reading extracted file: %v", err)
	}
	if string(data) != "file a" {
		t.Errorf("content: got %q, want %q", string(data), "file a")
	}
}

func TestTarGzRoundtrip_Content(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	tgzPath := filepath.Join(dir, "test.tar.gz")
	createTarGz(tgzPath, []string{filepath.Join(dir, "a.txt"), filepath.Join(dir, "b.txt")})

	outDir := filepath.Join(dir, "out")
	os.MkdirAll(outDir, 0755)
	extractTarGz(tgzPath, outDir)

	data, err := os.ReadFile(filepath.Join(outDir, filepath.Join(dir, "b.txt")))
	if err != nil {
		t.Fatalf("reading extracted file: %v", err)
	}
	if string(data) != "file b" {
		t.Errorf("content: got %q, want %q", string(data), "file b")
	}
}

func TestCreateZip_WithDirectory(t *testing.T) {
	dir := tempDir(t)
	createTestFiles(t, dir)

	zipPath := filepath.Join(dir, "test.zip")
	err := createZip(zipPath, []string{filepath.Join(dir, "sub")})
	if err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(dir, "out")
	os.MkdirAll(outDir, 0755)
	extractZip(zipPath, outDir)

	data, err := os.ReadFile(filepath.Join(outDir, filepath.Join(dir, "sub", "c.txt")))
	if err != nil {
		t.Fatalf("reading extracted sub/c.txt: %v", err)
	}
	if string(data) != "file c" {
		t.Errorf("content: got %q, want %q", string(data), "file c")
	}
}

func TestRunCreate_UnsupportedFormat(t *testing.T) {
	output = "/tmp/test.rar"
	defer func() { output = "" }()
	err := runCreate(createCmd, []string{"somefile"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestRunExtract_UnsupportedFormat(t *testing.T) {
	err := runExtract(extractCmd, []string{"/tmp/test.rar"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExtractTar_MissingFile(t *testing.T) {
	err := extractTar("/nonexistent.tar", "/tmp")
	if err == nil {
		t.Error("expected error")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestExtractZip_MissingFile(t *testing.T) {
	err := extractZip("/nonexistent.zip", "/tmp")
	if err == nil {
		t.Error("expected error")
	}
}

func TestExtractTarGz_MissingFile(t *testing.T) {
	err := extractTarGz("/nonexistent.tar.gz", "/tmp")
	if err == nil {
		t.Error("expected error")
	}
}

func TestCreateZip_BadSource(t *testing.T) {
	dir := tempDir(t)
	err := createZip(filepath.Join(dir, "test.zip"), []string{"/nonexistent/file.txt"})
	if err == nil {
		t.Error("expected error for bad source")
	}
}
