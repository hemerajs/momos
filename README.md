<p align="center">
    <img src="logo.png" alt="Momos logo" /><br /><br />
</p>

Momos - Reverse proxy to handle server-side-includes as custom elements without static configuration.

- **Cache:** Requests to your downstream are cached with Rfc7234
- **Fast:** SSI are loaded in parallel
- **No configs**: Everything is configurable via HTML5 attributes
- **Fallback**: Define default content or an error template `ssi-error`
- **Reliable**: Define timeouts to avoid hanging requests.
- **Custom elements**: Handle SSI blocks as custom elememts `ssi-*`

## Example
```html
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">
<body>

  <ssi-basket
    name="basket"
    cache="5000"
    timeout="2000"
    fallback="<url> | <elementId>"
    url="https://domain.de/basket">

   <!-- place here content on success -->

   <ssi-error>
    <span>Please try it again!</span>
   </ssi-error>
  </ssi-basket>
  
</body>
</html>
```

## Description

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

## TODO
- [ ] Use the net/http package to fetch the SSI Content
- [ ] Use the golang.org/x/net/html or [goquery](https://github.com/PuerkitoBio/goquery) to parse a web component
- [X] Use [httpcache](https://github.com/lox/httpcache) to provides an rfc7234 compliant caching

### Credit
Icon made by [author](https://www.flaticon.com/authors/dinosoftlabs) from www.flaticon.com
