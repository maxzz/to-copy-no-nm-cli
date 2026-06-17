# Renaming swap option to reserve
_Exported on 6/17/2026 at 14:25:27 PDT from Cursor (3.7.42)_

---

**User**

The `--swap` option needs to be renamed to `--reserve`, and all references to it must be updated; the logic remains the same as it was for `--swap`.

---

**Cursor**

I'll search the codebase for every `--swap` reference and related naming so we can rename it to `--reserve` consistently.


Renaming `--swap` to `--reserve` across code, tests, and docs; internal field `swapPaths` becomes `reservePaths`.


Updating the remaining historical doc reference for consistency.


Renamed `--swap` to `--reserve` everywhere. Behavior is unchanged: `-s` / `--reserve` still treats the first argument as destination and the second as source.

**Code**
- `parse_args.go` — flag `swap` → `reserve`, field `swapPaths` → `reservePaths`
- `validate_paths.go` — parameter `swap` → `reserve`
- `main.go` — help text and call site updated
- `parse_args_test.go` — assertions updated

**Docs**
- `README.md` — option table, path rules, and example comment
- `assets/docs/cursor_code_porting_to_go_with_output_m.md`
- `assets/docs/cursor_cli_utility_for_file_copying_in.md`
- `assets/TODOs.md` — marked the rename task done

Short flag `-s` is unchanged. `go test ./...` passes.

---

**User**

Oops, there was a typo: instead of the word "--reserve", you need to use the option name "--reverse". And one more thing: the short form for this option should be `-r`, not `-s` as it was before.

---

**Cursor**

Updating the option to `--reverse` with short form `-r` across the codebase.


Corrected: the option is now `-r` / `--reverse` (replacing `-s` / `--reserve`). Behavior is unchanged — first argument is destination, second is source when the flag is on.

**Updated**
- `parse_args.go` — `reversePaths`, flags `r` / `reverse`
- `validate_paths.go`, `main.go`, `parse_args_test.go`
- `README.md` and docs
- `assets/TODOs.md`

`go test ./...` passes.

---

**User**

The current file copying mode is essentially a full copy; a new mode needs to be implemented that copies only new and modified files. 
Additionally, if files have been deleted from the source, they must also be removed from the destination—meaning files present in the destination 
but absent from the source should be deleted. This mode will be called "synchronization" (or "light copying"). 
The current full-copy mode should be retained, but the synchronization mode will be used by default if 
no specific option '--full', or '-f' is selected. File comparison should be based on file size and modification time.
Also, organize the files that implement the synchronization and full-copy modes into separate modules.

---

**Cursor**

Exploring the copy implementation and CLI structure to add sync mode as the default and keep full copy behind `-f`/`--full`.


Implementing sync and full-copy modules and wiring the CLI so sync is the default.


Synchronization mode is implemented and is now the default. Summary of changes:

### New modules

| Module | Role |
|--------|------|
| `internal/4-syncdir` | **Sync mode (default)** — scans source and destination, copies new/changed files (size + mtime), removes extras via Recycle Bin, prunes empty orphan directories |
| `internal/5-fullcopy` | **Full copy mode** — wraps the existing clear-then-copy pipeline |
| `internal/2-copydir` | Shared copy engine; added `CopyRel` for single-file copies used by sync |

### CLI

- **Default:** sync mode (light copy)
- **`-f` / `--full`:** full copy (clear destination, then copy everything)

### Sync behavior

- Compares files by **size** and **modification time**
- Copies only missing or changed files
- Deletes destination files no longer in the source (Recycle Bin)
- Preserves `node_modules` in the destination (same as before)
- Respects root `.git` skip rules (unless `-g` is set)

### Examples

```bash
# Sync (default)
copy-no-nm "C:\projects\my-app" "D:\backups\my-app"

# Full copy
copy-no-nm -f "C:\projects\my-app" "D:\backups\my-app"
```

`go test ./...` passes, including new tests in `internal/4-syncdir`.

