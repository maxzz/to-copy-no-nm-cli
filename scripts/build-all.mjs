import { spawnSync } from 'node:child_process';
import { mkdirSync } from 'node:fs';
import { join } from 'node:path';

mkdirSync('dist', { recursive: true });

const targets = [
  { goos: 'windows', goarch: 'amd64', output: 'copy-no-nm.exe' },
  { goos: 'windows', goarch: 'arm64', output: 'copy-no-nm-arm64.exe' },
];

for (const target of targets) {
  const out = join('dist', target.output);
  const result = spawnSync(
    'go',
    ['build', '-ldflags=-s -w', '-o', out, './cmd/copy-no-nm'],
    {
      env: { ...process.env, GOOS: target.goos, GOARCH: target.goarch },
      stdio: 'inherit',
    },
  );

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
}
