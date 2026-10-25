# Contributing to Weshnet

Thank you for helping improve Weshnet. Contributions may include bug reports,
design proposals, documentation, tests, and code.

## Before you begin

- Search the [issue tracker](https://github.com/core-laboratories/weshnet/issues)
  for related work.
- Use a focused issue or pull request for each independent change.
- Discuss significant protocol or compatibility changes before implementation.
- Do not disclose security vulnerabilities in public issues. Follow the
  [security policy](SECURITY.md) instead.

## Development setup

Follow the [installation guide](docs/installation.md), then verify the checkout:

```sh
make test
```

The project uses the Go version and development tools declared by the
repository. Avoid unrelated dependency or generated-code changes.

## Making changes

1. Create a branch from the current default branch.
2. Add or update tests for behavior changes.
3. Keep exported Go APIs documented.
4. Update relevant guides and protocol documentation.
5. Run formatting, linting, and tests before opening a pull request.

```sh
make lint.fix
make lint
make unittest
```

When protocol definitions change, regenerate checked-in outputs with the
repository generation target:

```sh
make generate
```

Review the generated diff before committing it.

## Pull requests

A pull request should explain:

- the problem being solved;
- the chosen approach and important tradeoffs;
- compatibility or migration considerations;
- how the change was tested.

Keep commits reviewable and avoid rewriting unrelated files. All required
checks must pass before merge.

## Documentation style

- Use concise, descriptive headings.
- Use relative links for repository files.
- Add a language identifier to fenced code blocks.
- Wrap prose at a readable width where practical.
- Mark historical documents clearly so they are not mistaken for current APIs.

## License

By contributing, you agree that your contribution is provided under the
repository's [MIT license](LICENSE-MIT).
