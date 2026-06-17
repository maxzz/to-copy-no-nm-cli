# copy-no-nm

Recursively copy a directory to another location while skipping every `node_modules` folder. The copy preserves file creation times, modification times, access times, and attributes. Before copying, the destination folder is cleared selectively via the Windows Recycle Bin.

Built with Go. Distributed on npm with a small Node.js launcher.

## Contents

- [Requirements](#requirements)
- [Usage](#usage)
- [Options](#options)
- [Process overview](#process-overview)
  - [Step 1 — Clear destination](#step-1--clear-destination)
  - [Step 2 — Copy from source](#step-2--copy-from-source)
  - [Exit behaviour](#exit-behaviour)
- [Development](#development)
  - [Debugging](#debugging)
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

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `-r`, `--remove-node-modules` | **off** | Delete `node_modules` folders (including nested) in the destination before copying |
| `-g`, `--copy-git` | **off** | Copy the `.git` folder from the **source root** and clear the destination `.git` folder |

Examples:

```bash
# Defaults: keep destination node_modules and .git; skip copying source .git
copy-no-nm "C:\projects\my-app" "D:\backups\my-app"

# Also remove node_modules from destination before copy
copy-no-nm -r "C:\projects\my-app" "D:\backups\my-app"

# Also sync the root .git folder
copy-no-nm -g "C:\projects\my-app" "D:\backups\my-app"

# Both flags
copy-no-nm -r -g "C:\projects\my-app" "D:\backups\my-app"
```

## Process overview

Each run performs two steps: **clear destination**, then **copy from source**.

### Step 1 — Clear destination

| Item in destination | Default (`-r` off, `-g` off) | With `-r` | With `-g` |
|---------------------|--------------------------------|-----------|-----------|
| Files at any level | Recycle Bin | Recycle Bin | Recycle Bin |
| `node_modules` folders | **Kept** | Recycle Bin | **Kept** |
| Root `.git` folder | **Kept** | **Kept** | Recycle Bin |
| Other subfolders without nested `node_modules` | Recycle Bin (whole folder) | Recycle Bin | Recycle Bin |
| Subfolders containing nested `node_modules` | Cleared recursively with same rules | Cleared recursively with same rules | Cleared recursively with same rules |

### Step 2 — Copy from source

| Item in source | Default (`-g` off) | With `-g` |
|----------------|--------------------|-----------|
| All files and folders | Copied | Copied |
| `node_modules` folders (anywhere) | **Skipped** | **Skipped** |
| Root `.git` folder | **Skipped** | Copied |

File creation dates, timestamps, and attributes are preserved on copied items.

### Exit behaviour

- On error: Inspector Gadget is shown in red, then the app waits for a key press before closing.
- On success: Inspector Gadget is shown in green for 1.5 seconds, then the app closes.

## Development

Install dependencies (no runtime npm packages; pnpm is used for scripts only):

```bash
pnpm install
```

Build the Windows executable (generates `assets/icon.ico` from `assets/icon-source.png`, embeds it via `goversioninfo`, then compiles):

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

See `.vscode/launch.json` for ready-made debug configurations (F5 in Cursor/VS Code with the Go extension).

## npm scripts

| Script | Description |
|--------|-------------|
| `pnpm run build` | Generate icon, embed Windows resources, build `dist/copy-no-nm.exe` |
| `pnpm run prepublishOnly` | Build before publish (runs automatically on `npm publish`) |
| `pnpm run publish:npm` | Publish to npm (`--access public`) |

## Publishing to npm

This package ships a prebuilt Windows binary in `dist/` plus a Node launcher in `bin/copy-no-nm.js`.

```bash
pnpm run build
pnpm run publish:npm
```

Replace the temporary icon by editing `assets/icon-source.png`, then run `pnpm run build` again.

## License

MIT
