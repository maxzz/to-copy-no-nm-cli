# copy-no-nm

Recursively copy a directory to another location while skipping every `node_modules` folder. The copy preserves file creation times, modification times, access times, and attributes. Before copying, the destination folder is cleared selectively via the Windows Recycle Bin.

Built with Go. Distributed on npm with a small Node.js launcher.

## Contents

- [Requirements](#requirements)
- [Usage](#usage)
- [Paths](#paths)
- [Options](#options)
- [Process overview](#process-overview)
  - [Step 1 — Clear destination](#step-1--clear-destination)
  - [Step 2 — Copy from source](#step-2--copy-from-source)
  - [Console output](#console-output)
- [Development](#development)
  - [Debugging](#debugging)
  - [Windows shortcut](#windows-shortcut)
- [npm scripts](#npm-scripts)
- [Publishing to npm](#publishing-to-npm)
- [License](#license)

## Requirements

- Windows (Recycle Bin integration and metadata preservation use Windows APIs)
- [Go 1.22+](https://go.dev/dl/) for building from source
- [pnpm](https://pnpm.io/) for npm scripts

## Usage

```bash
copy-no-nm [options] <source> <destination>
```

Example:

```bash
copy-no-nm "C:\projects\my-app" "D:\backups\my-app"
```

On start the program prints its name, a short description, and the version (for example `copy-no-nm — Copy a folder… (version 0.26.2)`).

## Paths

| Path | Rule |
|------|------|
| `<source>` | Must exist and be a directory |
| `<destination>` | Created automatically if it does not exist (including parent folders); if it already exists, it must be a directory |

Source and destination must be different paths and cannot contain each other.

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `-r`, `--remove-node-modules` | **off** | Delete `node_modules` folders (including nested) in the destination before copying |
| `-g`, `--copy-git` | **off** | Copy the `.git` folder from the **source root** and clear the destination `.git` folder |

`node_modules` is always skipped during the copy, regardless of flags. This is done intentionally since `node_modules` should be created manually especially when using pnpm which creates links inside this folder.

Examples:

```bash
# Defaults: keep destination node_modules and .git; skip copying source .git
copy-no-nm "C:\projects\my-app" "D:\backups\my-app"

# Destination folder does not exist yet — it will be created
copy-no-nm "C:\projects\my-app" "D:\backups\new-folder"

# Also remove node_modules from destination before copy
copy-no-nm -r "C:\projects\my-app" "D:\backups\my-app"

# Also sync the root .git folder
copy-no-nm -g "C:\projects\my-app" "D:\backups\my-app"

# Both flags
copy-no-nm -r -g "C:\projects\my-app" "D:\backups\my-app"
```

Run without arguments (or with `-h`) to see in-program help: usage syntax, options with default values, and wrapped descriptions (~80 columns).

## Process overview

Each run performs two steps: **clear destination**, then **copy from source**.

### Step 1 — Clear destination

| Item in destination | Default (`-r` off, `-g` off) | With `-r` | With `-g` |
|---------------------|--------------------------------|-----------|-----------|
| Files at any level | Recycle Bin | Recycle Bin | Recycle Bin |
| `node_modules` folders | **Kept** (contents not scanned) | Recycle Bin | **Kept** (contents not scanned) |
| Root `.git` folder | **Kept** | **Kept** | Recycle Bin |
| Other subfolders without nested `node_modules` | Recycle Bin (whole folder) | Recycle Bin | Recycle Bin |
| Subfolders containing nested `node_modules` | Cleared recursively with same rules | Cleared recursively with same rules | Cleared recursively with same rules |

When `-r` is off, `node_modules` folders are left untouched and their contents are not inspected.

### Step 2 — Copy from source

| Item in source | Default (`-g` off) | With `-g` |
|----------------|--------------------|-----------|
| All files and folders | Copied | Copied |
| `node_modules` folders (anywhere) | **Skipped** | **Skipped** |
| Root `.git` folder | **Skipped** | Copied |

File creation dates, timestamps, and attributes are preserved on copied items (including read-only files such as git objects).

### Console output

| Situation | Behaviour |
|-----------|-----------|
| Missing or invalid arguments | Yellow help text with syntax, options, and defaults; press any key to close |
| Runtime error | Inspector Gadget in red, error message, press any key to close |
| Success | Inspector Gadget in green for 1.5 seconds, then the app closes |

## Development

Install dependencies (no runtime npm packages; pnpm is used for scripts only):

```bash
pnpm install
```

Build the Windows executable. The build script:

- Generates `assets/icon.ico` from `assets/icon-source.png` (only if the `.ico` file is not already present)
- Bumps the patch version in `package.json` and syncs `versioninfo.json`
- Embeds the icon and version metadata via `goversioninfo`
- Compiles `dist/copy-no-nm.exe`

```bash
pnpm run build
```

Run directly:

```bash
.\dist\copy-no-nm.exe <source> <destination>
```

Or via the npm bin shim after linking:

```bash
pnpm link --global
copy-no-nm <source> <destination>
```

### Debugging

Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) and [Delve](https://github.com/go-delve/delve):

```powershell
go install github.com/go-delve/delve/cmd/dlv@latest
```

Use `.vscode/launch.json` for ready-made configurations (F5 in Cursor/VS Code). Debug builds use `version dev` unless you set `-ldflags` manually.

### Windows shortcut

Create a shortcut and set **Target** to run the executable with preset arguments, for example:

```
C:\y\w\2-web\0-dp\utils\to-copy-no-nm-cli\dist\copy-no-nm.exe "C:\y\w\2-web\0-dp\utils\pmac" "C:\Users\maxzz\Desktop\pmac"
```

Double-clicking the shortcut runs the copy with those paths.

## npm scripts

| Script | Description |
|--------|-------------|
| `pnpm run build` | Bump version, generate icon (if needed), embed Windows resources, build `dist/copy-no-nm.exe` |
| `pnpm run prepublishOnly` | Build before publish (runs automatically on `npm publish`) |
| `pnpm run publish:npm` | Publish to npm (`--access public`) |

## Publishing to npm

This package ships a prebuilt Windows binary in `dist/` plus a Node launcher in `bin/copy-no-nm.js`.

```bash
pnpm run build
pnpm run publish:npm
```

Replace the icon by editing `assets/icon-source.png`, then run `pnpm run build` again.

## License

MIT
