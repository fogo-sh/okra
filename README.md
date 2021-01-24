# Okra

A lightweight, web-based podcast app.

## Development

### Running a development build
If you wish to work on this app, there are a number of quality-of-life improvements that are provided by running a dev build. At the moment, these include:
- Debug logs
- Automatic reloading templates, so you don't need to reboot the server anytime you edit HTML

To run a development build, do the following:
1. Ensure you don't have a `pkged.go` file in the root directory of the project. If this file is present, any changes to the static files you make won't be reflected until you rebuild the `pkged.go` file. You can safely delete this file.
2. Run `OKRA_DEBUG=true go run .`

### Performing a production build
A production build has all static assets (templates, css, etc.) embedded into the binary, as well as the symbol table omitted. This produces a smaller binary, that is entirely self-contained and capable of running the app on it's own. To perform a production build, run the following commands:

```bash
go get github.com/markbates/pkger/cmd/pkger
pkger && go build -ldflags '-s'
```
