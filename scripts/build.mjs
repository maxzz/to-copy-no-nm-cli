import { spawnSync } from 'node:child_process';
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const root = join(dirname(fileURLToPath(import.meta.url)), '..');
const iconPath = join(root, 'assets', 'icon.ico');
const pkgPath = join(root, 'package.json');
const versionInfoPath = join(root, 'cmd', 'copy-no-nm', 'versioninfo.json');
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

function bumpPatchVersion(version) {
    const parts = version.split('.').map((part) => Number.parseInt(part, 10));
    if (parts.length !== 3 || parts.some(Number.isNaN)) {
        throw new Error(`Invalid semver in package.json: ${version}`);
    }
    parts[2] += 1;
    return parts.join('.');
}

function updatePackageVersion() {
    const pkg = JSON.parse(readFileSync(pkgPath, 'utf8'));
    const newVersion = bumpPatchVersion(pkg.version);
    pkg.version = newVersion;
    writeFileSync(pkgPath, `${JSON.stringify(pkg, null, 4)}\n`, 'utf8');
    return newVersion;
}

function updateVersionInfo(version) {
    const [major, minor, patch] = version.split('.').map((part) => Number.parseInt(part, 10));
    const versionInfo = JSON.parse(readFileSync(versionInfoPath, 'utf8'));

    versionInfo.FixedFileInfo.FileVersion = { Major: major, Minor: minor, Patch: patch, Build: 0 };
    versionInfo.FixedFileInfo.ProductVersion = { Major: major, Minor: minor, Patch: patch, Build: 0 };
    versionInfo.StringFileInfo.FileVersion = version;
    versionInfo.StringFileInfo.ProductVersion = version;

    writeFileSync(versionInfoPath, `${JSON.stringify(versionInfo, null, 4)}\n`, 'utf8');
}

const version = updatePackageVersion();
updateVersionInfo(version);
console.log(`Building copy-no-nm v${version}`);

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
    `-ldflags=-s -w -X main.version=${version}`,
    '-o',
    'dist/copy-no-nm.exe',
    './cmd/copy-no-nm',
]);
