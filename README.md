# SSSH

## Build for macOS with Docker
```sh
docker run -v ${PWD}:/src -w /src golang:latest env GOOS=darwin GOARCH=amd64 go build
```

## Themes
```sh
export SSSH_THEME=<theme>
```
Possible values:
- `base`
- `base16`
- `dracula`
- `charm`
- `catppuccin`

(default: `dracula`)