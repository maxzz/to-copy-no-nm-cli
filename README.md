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

1. All files and folders inside `<destination>` are moved to the Recycle Bin.
2. `<source>` is copied recursively into `<destination>`.
3. Any directory named `node_modules` is skipped entirely.
4. File creation dates, timestamps, and attributes are preserved.
5. On error: the message is shown in red and the app waits for a key press before closing.
6. On success: a green Inspector Gadget ASCII image is shown for 1.5 seconds, then the app closes.

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
