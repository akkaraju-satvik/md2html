# md2htm

## Description

A simple markdown to html converter.

## Usage

Use this command to convert a markdown file to html:

```bash
go run main.go -f <input.md>
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
pageTitle: TITLE
description: DESCRIPTION
version: VERSION
author:
  name: AUTHOR
  email: EMAIL
github: GITHUB
```
