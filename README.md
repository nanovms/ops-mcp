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

Available tools:

```
List instances
```

```
List images
```

```
Instance create <image_name>
```

```
Instance create redis-server
```

Note: Very open to suggestions on how this all should work as this initial cut was done not having
ever used Claude or MCP.
