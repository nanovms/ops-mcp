# ops-mcp
mcp server for ops

Build like normally.

Put this in your Claud config:

```
~/Library/Application Support/Claude/claude_desktop_config.json
```

Ensure your command is in the right path and more importantly the PATH
env is set to run.

```
{
"mcpServers": {
  "ops-mcp": {
      "command": "/Users/eyberg/go/src/github.com/nanovms/ops-mcp/ops-mcp",
      "args": [],
      "env": {
        "HOME":"/Users/eyberg",
        "LOGNAME":"eyberg",
        "PATH":"/bin:/Users/eyberg/.ops/bin",
        "SHELL":"/bin/zsh",
        "USER":"eyberg"
        }
    }
  }
}
```
