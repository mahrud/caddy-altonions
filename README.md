# Add PROXY protocol support to caddy

## Syntax

```
proxyprotocol *cidr* ... {
    timeout *val*
}
```

- **cidr** CIDR ranges to process PROXY headers from
- **val** duration value (e.g. 5s, 1m)

The default timeout is `5s`. Set to `0` or `none` to disable the timeout.

## Examples

```
# Enable from any source (probably don't want this in prod)
proxyprotocol

# Enable from local subnet and fixed IP
proxyprotocol 10.22.0.0/16 10.23.0.1/32

# Set header timeout
proxyprotocol 10.22.0.0/16 10.23.0.1/32 {
    timeout 5s
}

```
