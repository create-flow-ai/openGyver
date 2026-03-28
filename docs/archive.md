# archive

Create and extract archive files in various formats.

## Usage

```bash
openGyver archive [command] [flags]
```

## Global Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--help` | `-h` | bool | `false` | Show help for archive |

## Supported Formats

| Format | Extensions | Create | Extract | Requirements |
|--------|-----------|--------|---------|--------------|
| ZIP | `.zip` | Yes | Yes | Pure Go (built-in) |
| TAR | `.tar` | Yes | Yes | Pure Go (built-in) |
| TAR.GZ | `.tar.gz`, `.tgz` | Yes | Yes | Pure Go (built-in) |
| 7Z | `.7z` | Yes | Yes | Requires `p7zip` / `7z` |
| RAR | `.rar` | No | Yes | Requires `unrar` |
| BZ2 | `.tar.bz2`, `.tbz2` | Yes | Yes | Requires `bzip2` |

## Subcommands

### create

Create a ZIP, TAR, TAR.GZ, TAR.BZ2, or 7Z archive from the given files and directories. The archive format is determined by the `--output` file extension.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `file` | Yes | One or more files or directories to add to the archive |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `""` | Output archive path (required) |
| `--help` | `-h` | bool | `false` | Show help for create |

#### Examples

```bash
# Create a ZIP archive from individual files
openGyver archive create -o backup.zip file1.txt file2.txt

# Create a ZIP archive including a directory
openGyver archive create -o backup.zip file1.txt dir/

# Create a tar.gz archive from source code and readme
openGyver archive create -o project.tar.gz src/ README.md

# Create a plain TAR archive
openGyver archive create -o files.tar doc1.txt doc2.txt

# Create archive from glob patterns
openGyver archive create -o data.tar *.csv

# Create a bzip2-compressed tar archive (requires bzip2)
openGyver archive create -o archive.tar.bz2 logs/ config/

# Create a 7z archive (requires 7z/p7zip)
openGyver archive create -o compressed.7z large-file.bin images/

# Archive an entire project directory
openGyver archive create -o release.zip dist/ package.json LICENSE
```

---

### extract

Extract a ZIP, TAR, TAR.GZ, TAR.BZ2, 7Z, or RAR archive. Extracts to the current directory by default. Use `-o` to specify a target directory.

#### Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `archive` | Yes | Path to the archive file to extract |

#### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `.` (current directory) | Extraction destination directory |
| `--help` | `-h` | bool | `false` | Show help for extract |

#### Examples

```bash
# Extract a ZIP archive to the current directory
openGyver archive extract backup.zip

# Extract a ZIP to a specific directory
openGyver archive extract backup.zip -o ./restored/

# Extract a tar.gz archive
openGyver archive extract project.tar.gz -o ./project/

# Extract a plain TAR archive
openGyver archive extract files.tar -o ./output/

# Extract a bzip2-compressed tar archive
openGyver archive extract archive.tar.bz2 -o ./data/

# Extract a 7z archive (requires 7z)
openGyver archive extract compressed.7z -o ./unpacked/

# Extract a RAR archive (requires unrar)
openGyver archive extract files.rar -o ./extracted/

# Extract to a nested path (directories are created automatically)
openGyver archive extract release.zip -o ./builds/v2.0/
```

## Notes

- The archive format for both create and extract is auto-detected from the file extension.
- ZIP, TAR, and TAR.GZ formats are implemented in pure Go and require no external tools.
- 7Z support requires `p7zip` (`brew install p7zip` on macOS, `apt install p7zip-full` on Linux).
- RAR extraction requires `unrar` (`brew install unrar` on macOS, `apt install unrar` on Linux). RAR creation is not supported.
- BZ2 support requires the `bzip2` command-line tool.
- Path traversal attacks (zip slip) are guarded against during extraction.
