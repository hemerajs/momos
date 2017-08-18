<p align="center">
    <img src="logo.png" alt="Momos logo" /><br /><br />
</p>

Momos - Reverse proxy to define server-side-includes via HTML5 and attributes. No html comments or complicate configurations. This is proof-of-concept. 

- **Cache:** Requests are cached with RFC7234
- **Fast:** SSI Fragments are loaded in parallel
- **No proxy configs**: Everything is configurable via HTML5 attributes
- **Dev-friendly**: Frontend developer can create fragments easily
- **Fallback**: Define default content or an error template with `<ssi-error>`
- **Reliable**: Define timeout message with `<ssi-timeout>`
- **Just HTML**: Define SSI fragments with HTML element `<ssi>`

## Why you don't use Nginx?
Good point. Nginx is a great proxy and although it already provides robust SSI directives I would like to see a solution which don't require a restart or reload of the proxy when parameters has to be changed. The transitions between defining SSI fragments and configure them should be smooth for any kind of developer. Momos should provide a high performance proxy to provide advanced SSI functionality. Any developer should be able to place SSI fragments and configure them with only html knowledge. Momos is very easy to extend and is compiled to a single binary. It also provide great debugging functionality to understand how your page is build.

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
- [ ] Generate great debug informations about the structure of your page
- [ ] Use single http client and create request with http.NewRequest
- [ ] Start multiple requests concurrently with go channels

### References
- [Microservice-websites](https://gustafnk.github.io/microservice-websites/#integration-techniques)
- [Apache SSI](https://httpd.apache.org/docs/current/howto/ssi.html#page-header)
### Credits
Icon made by [author](https://www.flaticon.com/authors/dinosoftlabs) from www.flaticon.com
