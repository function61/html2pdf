![Build status](https://github.com/function61/html2pdf/workflows/Build/badge.svg)

What?
-----

A small microservice that turns HTML into a PDF file. You can run this:

- on AWS Lambda
- with Docker
  * I didn't bother making a `Dockerfile` though, since I didn't need it. PR welcome!
- as a standalone binary

There also exists [a small client library for Go](pkg/html2pdfclient/)


Testing
-------

You can start a local server process with:

```console
$ html2pdf server
```

Then call it from the client:

```console
$ export HTML2PDF_TOKEN="doesntMatter" # optionally you can put the service behind authentication
$ html2pdf client-localhost '<h1>hello world</h1>' > out.pdf
```

Usage from curl is also simple:

```console
$ curl -d '{"html_base64": "PGgxPmhlbGxvIHdvcmxkPC9oMT4="}' http://localhost/render > out.pdf
```


Prerequisites for dev/testing
-----------------------------

```console
$ apt install -y libxrender1 libxext6 libfontconfig1
```

(Fortunately, these exist in Lambda's AMI)


Alternatives
------------

- https://github.com/shouldbee/docker-html2pdf (go, wkhtmltox)
- https://github.com/carlescliment/html2pdf-service (PHP, wkhtmltopdf)
- https://github.com/spipu/html2pdf (PHP, self-implemented HTML rendering)


Support / contact
-----------------

Basic support (no guarantees) for issues / feature requests via GitHub issues.

Paid support is available via [function61.com/consulting](https://function61.com/consulting/)

Contact options (email, Twitter etc.) at [function61.com](https://function61.com/)
