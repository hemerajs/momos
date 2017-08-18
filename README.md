<p align="center">
    <img src="logo.png" alt="Momos logo" /><br /><br />
</p>

Momos - Reverse proxy to define server-side-includes via HTML5 and attributes. This is proof-of-concept. 

- **Cache:** Requests to your downstream are cached with Rfc7234
- **Fast:** SSI are loaded in parallel
- **No proxy configs**: Everything is configurable via HTML5 attributes
- **Dev-friendly**: Frontend developer can easily create fragments
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

  <ssi
    name="basket"
    cache="5000"
    timeout="2000"
    fallback="<url> | <elementId>"
    url="https://domain.de/basket">
    
    <!-- Default content -->
    
    <ssi-timeout>
    <span>Please try it again!</span>
   </ssi-timeout>
   
   <ssi-error>
    <span>Please call the support!</span>
   </ssi-error>
  </ssi>
  
</body>
</html>
```

## TODO
- [ ] Use the net/http package to fetch the SSI Content
- [ ] Use the golang.org/x/net/html or [goquery](https://github.com/PuerkitoBio/goquery) to parse a web component
- [X] Use [httpcache](https://github.com/lox/httpcache) to provides an rfc7234 compliant caching
- [ ] Use [httpcache](https://github.com/gregjones/httpcache) to provides an rfc7234 compliant client caching for SSI requests

### Credits
Icon made by [author](https://www.flaticon.com/authors/dinosoftlabs) from www.flaticon.com
