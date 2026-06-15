import { spawnSync } from 'node:child_process';
import { existsSync, mkdirSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = join(dirname(fileURLToPath(import.meta.url)), '..');
const iconPath = join(root, 'assets', 'icon.ico');
const goversioninfo = 'github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.7.0';

mkdirSync(join(root, 'dist'), { recursive: true });

function run(command, args, options = {}) {
    const result = spawnSync(command, args, {
        cwd: root,
        stdio: 'inherit',
        ...options,
    });

    if (result.status !== 0) {
        process.exit(result.status ?? 1);
    }
}

if (!existsSync(iconPath)) {
    run('powershell', [
        '-NoProfile',
        '-ExecutionPolicy',
        'Bypass',
        '-File',
        join(root, 'scripts', 'generate-icon.ps1'),
    ]);
}

run('go', [
    'run',
    goversioninfo,
    '-o',
    'cmd/copy-no-nm/resource.syso',
    'cmd/copy-no-nm/versioninfo.json',
]);

run('go', [
    'build',
    '-ldflags=-s -w',
    '-o',
    'dist/copy-no-nm.exe',
    './cmd/copy-no-nm',
]);
