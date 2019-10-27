# Akmey client (new)

As the "old" [Akmey client](https://github.com/akmey/akmey-client) was a bit hard to maintain and CLI wasn't very well documented, I decided to migrate to the [Cobra](https://github.com/spf13/cobra) framework.

This version is not complete, and is not meant to replace the actual Akmey client before this refactor is completed, so not in the short future.

Functions to add:
- [x] install
- [x] remove (kinda, it ereases everything for now)
- [x] reset (basically reset)
- [ ] team-install

Example of .akmey.json:
```json
{
	"server": "https://akmey.leonekmi.fr",
	"keyfile": "/home/yourusername/.ssh/authorized_keys"
}
```
