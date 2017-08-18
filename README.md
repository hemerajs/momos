<p align="center">
    <img src="logo.png" alt="Momos logo" /><br /><br />
</p>

Momos - Reverse proxy to define server-side-includes via HTML5 and attributes. This is proof-of-concept. 

- **Cache:** Requests are cached with RFC7234
- **Fast:** SSI Fragments are loaded in parallel
- **No proxy configs**: Everything is configurable via HTML5 attributes
- **Dev-friendly**: Frontend developer can create fragments easily
- **Fallback**: Define default content or an error template with `<ssi-error>`
- **Reliable**: Define timeout message with `<ssi-timeout>`
- **Just HTML**: Define SSI fragments with HTML element `<ssi>`

## What are SSI?

> SSI (Server Side Includes) are directives that are placed in HTML pages, and evaluated on the server while the pages are being served. They let you add dynamically generated content to an existing HTML page, without having to serve the entire page via a CGI program, or other dynamic technology.
[Reference](https://httpd.apache.org/docs/current/howto/ssi.html#page-header)


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
    fallback=""
    src="http://starptech.de">

    <!-- a) Used when no error or timeout field was set -->
    Default content!
    
    <!-- b) Timeout errors based on the timeout duration-->
    <ssi-timeout>
    <span>Please try it again!</span>
    </ssi-timeout>
    
    <!-- c) statusCode > 199 && statusCode < 300 -->
    <!-- d) Any other error -->
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
- [ ] Generate great debug informations about the structure of your page.
- [ ] Use single http client and create request with http.NewReques
- [ ] Start multiple requests concurrently with go channels

### References
- [Microservice-websites](https://gustafnk.github.io/microservice-websites/#integration-techniques)
- [Apache SSI](https://httpd.apache.org/docs/current/howto/ssi.html#page-header)
### Credits
Icon made by [author](https://www.flaticon.com/authors/dinosoftlabs) from www.flaticon.com
