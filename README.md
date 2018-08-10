# Opportunistic Onions for caddy

## Syntax

```
altonions *addr:port* ... {
    max-age *ma*
	persist *persist*
}
```

- **addr:port** the onion service address
- **ma**, integer, max-age value in seconds
- **persist**, integer

## Examples

```
perfectoid.space:8443 {
    tls perfectoid.pem perfectoid-key.pem
    altonions zkiefsz3zbkg4nnl5p7r64qxugfeb7g5agz2pqwci4w7hwzfgu2gobad.onion:8443 {
        ma 086400
        persist 1
    }
    root /var/www/html/
}
```
