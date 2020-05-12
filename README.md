# Akmey client 🔑

## Installation 🚀

![Does master builds?](https://github.com/akmey/akmey-client-new/workflows/Go/badge.svg)

On Archlinux, you can either use the `akmey` or the `akmey-bin` package from AUR:

```
yay -S akmey
```

On any other platform, you can simple build it with Go:

```
go get github.com/akmey/akmey-client
```

Pre-compiled binary are also available for a selection of platform (Linux, MacOS and Windows).

## Why ❓

As the "old" [Akmey client](https://github.com/akmey/akmey-client) was a bit hard to maintain and CLI wasn't very well documented, I decided to migrate to the [Cobra](https://github.com/spf13/cobra) framework.

For now, some feature are missing, but every basic functionaliy are working great.

Functions to add:

- [x] install
- [x] remove
- [x] reset
- [ ] team-install

## Configuration 📝

The main configuration file is config.json, located at ~/.akmey/config.json.
Example of config.json:

```json
{
  "server": "https://akmey.leonekmi.fr",
  "keyfile": "/home/yourusername/.ssh/authorized_keys"
}
```

## License

This client is licensed under the [Unlicense](https://github.com/akmey/akmey-client/blob/master/LICENSE).

