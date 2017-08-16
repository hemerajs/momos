# Momos
Reverse proxy to handle server-side-includes as web components without static configurations.

```
+-----+    +---------+    +-------+
|     |    |         |    |       |
| API +---->  Proxy  +----> NGINX |
|     |    |         |    |       |
+-----+    +----+----+    +-------+
                |
                |
           +----v----+
           |         |
           |   SSI   |
           |         |
           +---------+

```

## Proposal
```html
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">
<body>

  <ssi-basket
    cache="5000"
    timeout="2000"
    fallback="<url> | <elementId>"
    url="https://domain.de/basket"
    query="userId=5"
    headers="x-request-id=123456">

   <!-- place here content on success -->

   <ssi-error>
    <span>Please try it again!</span>
   </ssi-error>
  </ssi-basket>
  
</body>
</html>
```

## Example

```
go run examples/example.go
```
