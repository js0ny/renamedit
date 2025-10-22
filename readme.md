# `renamedit`

Rename files with your favourite editor

## Installation

Build from source:

```bash
git clone https://github.com/js0ny/renamedit.git
cd renamedit
go build
```

## Usage

```bash
renamedit [options] <directory>
```

### Options

- `-ignore-ext`: Ignore file extensions when renaming files. This will preserve the original extensions.

### Examples

Basic usage:
```bash
renamedit /path/to/directory
```

Ignore file extensions (useful for batch renaming while preserving extensions):
```bash
renamedit -i /path/to/directory
renamedit --ignore-ext /path/to/directory
```

## How It Works

1. The program opens a temporary file in your preferred text editor
2. Edit the file names as needed
3. Save and exit the editor
4. The program will rename the files according to your edits

## Environment Variables

- `EDITOR`: Set this to your preferred text editor (defaults to vim)

## Supported Editors

The following editors are explicitly supported with wait flags:
- Visual Studio Code (code)
- Sublime Text (subl)
- Zed Editor (zeditor)
- Atom (atom)
- Gedit (gedit)
- And most terminal editors (vim, nano, etc.)
