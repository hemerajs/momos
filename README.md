<p align="center">
    <img src="logo.png" alt="Momos logo" /><br /><br />
</p>

[![License MIT](https://img.shields.io/badge/License-MIT-blue.svg)](http://opensource.org/licenses/MIT)
[![Build Status](https://travis-ci.org/hemerajs/momos.svg?branch=master)](http://travis-ci.org/hemerajs/momos)

Momos - Reverse proxy to define server-side-includes via HTML5 and attributes. No html comments or complicate configurations. This is proof-of-concept. 

- **Cache:** Requests are cached with RFC7234 with support for memory and file storage.
- **Lightweight:** Just ~300 lines of code. We trust only in well-tested packages.
- **Fast:** SSI Fragments are loaded in parallel.
- **No proxy configs**: Everything is configurable via HTML5 attributes.
- **Dev-friendly**: Frontend developer can create fragments easily.
- **Fallback**: Define default content or an error template with `<ssi-error>`.
- **Reliable**: Define a timeout message with `<ssi-timeout>`.
- **Just HTML**: Define SSI fragments with pure HTML `<ssi>`.
- **Templating**: Use Go Templates inside fragments and ssi tags.

## Why you don't use Nginx?
Good point. Nginx is a great proxy and although it already provides robust SSI directives I would like to see a solution which don't require a restart or reload of the proxy when parameters has to be changed. The transition between defining SSI fragments and configure them should be smooth for any kind of developer. Momos should provide a high performance proxy with advanced SSI functionality. Any developer should be able to place and configure SSI fragments with html knowledge. Momos is very easy to extend and is compiled to a single binary. It provides great debugging experience to understand how your page is build which is often difficult in proxys like Nginx or Apache.

## What are SSI?

> SSI (Server Side Includes) are directives that are placed in HTML pages, and evaluated on the server while the pages are being served. They let you add dynamically generated content to an existing HTML page, without having to serve the entire page via a CGI program, or other dynamic technology.
[Reference](https://httpd.apache.org/docs/current/howto/ssi.html#page-header)

## Advantages

- Easy integration of html fragments from external services
- Share all layout html fragments to keep those snippets DRY
- In highly distributed environments there is no other way to include dynamic components without using javascript


## Example
```html
<!doctype html>

<html lang="en">
<head>
  <meta charset="utf-8">
<body>

  <ssi
    name="basket"
    timeout="2000"
    template="true"
    src="http://starptech.de">

    <!-- Used when no error or timeout field was set default empty space -->
    Default content!
    
    <!-- Timeout errors based on the timeout duration-->
    <ssi-timeout>
    <span>Please try it again! {{.DateLocal}}</span>
    </ssi-timeout>
    
    <!-- None 2xx status code or any other error -->
    <ssi-error>
    <span>Please call the support!</span>
    </ssi-error>
  </ssi>
  
</body>
</html>
```

- `name`      : The name of the fragment (default `unique-id`)
- `timeout`   : The maximum request timeout (default `2000`)
- `no-scripts`: Filter javascript and css includes from the fragment (default `true`)
- `src`       : The url of the server-side-include
- `template`  : Enables template rendering via go templates (default `false`)

## Run it

```
$ go get ./...
$ go run examples/server.go
$ go run examples/client.go
```
### Expected output
Requests are cached for 10 seconds `max-age=10`
```
2017/08/20 16:34:41.162786 [INF] PreResponse url: localhost:9090/, cache: SKIP
2017/08/20 16:34:41.196786 [INF] PreResponse url: localhost:9090/favicon.ico, cache: SKIP
2017/08/20 16:34:41.198786 [INF] Fragment "basket" was cached
2017/08/20 16:34:41.199786 [TRC] Call fragment basket, url: http://localhost:8081, duration: 1.0001ms
2017/08/20 16:34:41.198786 [INF] Fragment "basket4" was cached
2017/08/20 16:34:41.198786 [INF] Fragment "basket5" was cached
2017/08/20 16:34:41.199786 [INF] Fragment "basket2" was cached
2017/08/20 16:34:41.198786 [INF] Fragment "basket3" was cached
2017/08/20 16:34:41.199786 [TRC] Call fragment basket4, url: http://localhost:8081/c, duration: 1.0001ms
2017/08/20 16:34:41.200785 [TRC] Call fragment basket5, url: http://localhost:8081/d, duration: 1.9995ms
2017/08/20 16:34:41.202787 [TRC] Call fragment basket2, url: http://localhost:8081/a, duration: 4.0014ms
2017/08/20 16:34:41.203787 [TRC] Call fragment basket3, url: http://localhost:8081/b, duration: 5.0014ms
2017/08/20 16:34:41.204786 [TRC] Processing complete "http://127.0.0.1:8080/favicon.ico" took "8ms"
```


## TODO
- [X] Use the net/http package to fetch the SSI Content
- [X] Use [goquery](https://github.com/PuerkitoBio/goquery) to parse a web component
- [X] Use [httpcache](https://github.com/lox/httpcache) to provides an rfc7234 compliant caching
- [X] Use [httpcache](https://github.com/gregjones/httpcache) to provides an rfc7234 compliant client caching for SSI requests
- [ ] Generate great debug informations about the structure of your page
- [X] Use single http client and create request with http.NewRequest
- [X] Start multiple requests concurrently with go channels
- [ ] Collect metrics
- [ ] Write tests

### References
- [Microservice-websites](https://gustafnk.github.io/microservice-websites/#integration-techniques)
- [Apache SSI](https://httpd.apache.org/docs/current/howto/ssi.html#page-header)
### Credits
Icon made by [author](https://www.flaticon.com/authors/dinosoftlabs) from www.flaticon.com
