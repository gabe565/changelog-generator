# Changelog Generator
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/gabe565/changelog-generator)](https://github.com/gabe565/changelog-generator/releases)
[![Build](https://github.com/gabe565/changelog-generator/actions/workflows/build.yaml/badge.svg)](https://github.com/gabe565/changelog-generator/actions/workflows/build.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabe565/changelog-generator)](https://goreportcard.com/report/github.com/gabe565/changelog-generator)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gabe565_changelog-generator&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gabe565_changelog-generator)

A configurable commit-based changelog generator. It finds commits since the previous release, filters them, then groups them.

## Usage

### GitHub Action

#### Inputs

| Name    | Description                                | Default             |
|---------|--------------------------------------------|---------------------|
| `token` | GitHub token used to fetch release assets. | `${{ github.token }}` |

#### Outputs

| Name        | Description                       |
|-------------|-----------------------------------|
| `changelog` | The generated changelog markdown. |


#### Example
```yaml
name: Release

on:
  push:
    tags:
      - "v*.*.*"

jobs:
  release:
    name: Release
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Generate Changelog
        uses: gabe565/changelog-generator-action@v1
        id: changelog
      - name: Release
        uses: softprops/action-gh-release@v2
        with:
          body: ${{ steps.changelog.outputs.changelog }}
```

## Configuration

Configuration is loaded from `.changelog-generator.yaml` in the git repo root. See the [config example](config_example.yaml) for more details.
