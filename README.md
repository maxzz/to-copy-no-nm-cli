# copy-no-nm

Recursively copy a directory to another location while skipping every `node_modules` folder. The copy preserves file creation times, modification times, access times, and attributes. Before copying, everything already in the destination folder is sent to the Windows Recycle Bin.

Built with Go. Distributed on npm with a small Node.js launcher.

## Requirements

- Windows (Recycle Bin integration and metadata preservation use Windows APIs)
- [Go 1.22+](https://go.dev/dl/) for building from source
- [pnpm](https://pnpm.io/) for npm scripts

## Usage

```bash
copy-no-nm <source> <destination>
```

Example:

```bash
copy-no-nm "C:\projects\my-app" "D:\backups\my-app"
```

### Behavior

1. Existing destination contents are cleared to the Recycle Bin using selective rules (`node_modules` folders are kept by default).
2. `<source>` is copied recursively into `<destination>`.
3. Any directory named `node_modules` is skipped during the copy.
4. File creation dates, timestamps, and attributes are preserved.
5. On error: Inspector Gadget is shown in red, then the app waits for a key press before closing.
6. On success: Inspector Gadget is shown in green for 1.5 seconds, then the app closes.

### Clearing the destination

By default, `node_modules` folders in the destination are **not** deleted. Other subfolders are removed as a whole unless they contain a nested `node_modules` folder; in that case the same rules are applied inside that subfolder.

Use `--remove-node-modules` to also send `node_modules` folders (including nested ones) to the Recycle Bin.

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

For cross-platform npm installs, consider splitting platform-specific binaries into separate packages listed under `optionalDependencies`.

Replace the temporary icon by editing `assets/icon-source.png`, then run `pnpm run build` again.

## License

MIT
