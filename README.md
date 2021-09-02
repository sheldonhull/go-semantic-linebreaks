# Go Semantic Line Breaks For Markdown

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-21%25-brightgreen.svg?longCache=true&style=flat)</a>

[Semantic Line Breaks](https://sembr.org/)

🚧 This is a work in progress.

- Used to test creating a simple linter cli using TDD.
- Doesn't have any intelligence on parsing a markdown, just raw text, so might break snippet code blocks in some language I haven't seen yet, though my initial tests for my blog worked well.

## Install

Go 1.16+

```shell
go install https://github.com/sheldonhull/go-semantic-linebreaks/cmd/go-semantic-linebreaks@latest
```

## Use

```shell
go-semantic-linebreaks -source ./markdowndirectory
```
