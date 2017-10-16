# How to document a Giant Swarm product

We provide product-specific documentation inside the repository of the product, using a top-level directory named `docs`.

The content within this directory is intended primarily for reading on the `giantswarm.io` website, while also being accessible directly on `github.com`. To enable useful rendering in both contexts, please adhere to the following guidelines.

## Markdown content

Documentation content is in [Markdown](https://daringfireball.net/projects/markdown/syntax) format. Github-flavored Markdown extensions such as fenced code blocks and tables might be applied.

## Index files

We make use of the fact that when entering a directory of a repository, GitHub automatically renders the `README.md` file content if available. To allow for convenient reading/browsing on github.com, it is recommended to provide a `README.md` file with a table of contents for everything within the directory.

On `giantswarm.io`, if no `README.md` file is available in the `docs` folder or any subfolder, a table of contents page is automatically generated in it's place for browsing that folder.

## Linking between documentation pages

Always use *relative links* to reference pages within the documentation of the same repository. This guarantees that the hyperlinks also work in the `giantswarm.io` website context.

To link from page `README.md` to page `configuration.md` in the same folder, write your link like this:

```markdown
... more about the [Configuration](configuration.md) of ...
```

In order to link to a page in the parent directory, accordingly use this relative syntax:

```markdown
... more about the [Configuration](../configuration.md) of ...
```

To link to the index page of a sub-directory of the current page, directly reference the `README.md` file:

```markdown
... more about the [Configuration](configuration/README.md) of ...
```

## Nested directories

You can nest documentation content in as many sub-directories as you consider useful. Here, the same rules apply for automatic index page generation and the use of `README.md` to manually create index page content.

## Images

You can store images used within the documentation inside the `docs` folder wherever you like. It is recommended to create a common folder `docs/images` to make all images available in a common place.

## Syntax Highlighting

When rendering fenced code blocks, syntax highlighting is applied automatically. Some languages can be detected rather well. However, some are rather problematic and need help.

Especially shell commands and their output regularly lead to bad syntax highlighting results. For that sort of code, syntax highlighting should therefore be disabled completely. There is a helper `nohighlight` that can be added to the according code block. Example:

```nohighlight
```nohighlight
$ some shell code
...
```

To indicate a specific language used, these helpers are at your disposal: `dockerfile`, `go`, `java`, `javascript`, `json`, `php`, `python`, `ruby`, `yaml`.
