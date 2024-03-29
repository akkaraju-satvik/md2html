---
pageTitle: Home
---
# md2htm

## Description

A simple markdown to html converter.

## Installation

### Download the binary

You can download the binary for your system from the [releases](https://github.com/akkaraju-satvik/md2htm/releases) page.

### Build from source using Go

```bash
git clone
cd md2htm
go install
```

## Usage

Use this command to convert a markdown file to html:

```bash
md2htm -f <input.md>
```

Metadata can be added to the markdown file using the following format:

```markdown
---
pageTitle: TITLE
authorName: AUTHOR
description: DESCRIPTION
---
```

## Configuration

You can add your data to a `yaml` file and use it to provide metadata to your page.

```yaml
projectTitle: TITLE
description: DESCRIPTION
version: VERSION
author:
  name: AUTHOR
  email: EMAIL
github: GITHUB
assetsDir: ASSETS_DIR
outputDir: OUTPUT_DIR
favicon: FAVICON
```

You can then use the `--config-file` flag to specify the configuration file.

```bash
md2htm -f <input.md> --config-file <config-file>
# or
md2htm -f <input.md> -c <config-file>
```

## Customization

You can customize the template used for the conversion by specifying the `--template-file` flag and providing a custom template file.

```bash
md2htm -f <input.md> --template-file <template-file>
# or
md2htm -f <input.md> -t <template-file>
```

## Output

The output will be saved in the `outputDir` specified in the configuration file. If no output directory is specified, the output will be saved in the `dist` directory and the file name will be the same as the input file with the `.html` extension unless the `--output-file` flag is used to specify a custom output file.

```bash
md2htm -f <input.md> --output-file <output-file>
# or
md2htm -f <input.md> -o <output-file>
```

## Fun Fact

This page was created using this tool.
