# version_bumper

`version_bumper` is a simple Go tool for incrementing version numbers in a VERSION file. It supports bumping major, minor, and patch versions.

## Installation

To install `version_bumper`, make sure you have Go installed on your system, then run:

```bash
go install github.com/lavabyrd/version_bumper@latest
```

## Usage with exist Makefiles

```bash
VERSION=$(version_bumper --dry-run) make bump
```

By default it will bump the patch version. You can pass `--minor` or `--major` to bump those instead.

```bash
VERSION=$(version_bumper --dry-run --major) make bump
```
