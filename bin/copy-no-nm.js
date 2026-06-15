#!/usr/bin/env node

import { spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const root = join(__dirname, '..');
const binaryName = process.platform === 'win32' ? 'copy-no-nm.exe' : 'copy-no-nm';
const binaryPath = join(root, 'dist', binaryName);

if (!existsSync(binaryPath)) {
  console.error(`copy-no-nm: binary not found at ${binaryPath}. Run "pnpm run build" first.`);
  process.exit(1);
}

const result = spawnSync(binaryPath, process.argv.slice(2), { stdio: 'inherit' });
process.exit(result.status ?? 1);
