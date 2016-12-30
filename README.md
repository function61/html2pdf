What?
-----

A small Docker-based microservice that takes HTML as input and renders it as PDF.


Starting the server
-------------------

```
$ docker run -d --name html2pdf -p 8080:80 fn61/html2pdf
```

Testing the conversion process:

```
$ curl --form document=@example_html_input/report.html http://localhost:8080/render > report.pdf
```


Metrics
-------

[Prometheus](https://prometheus.io/) metrics are available at /metrics, debug:

```
$ curl http://localhost:8080/metrics
```


Alternatives
------------

- https://github.com/shouldbee/docker-html2pdf (go, wkhtmltox)
- https://github.com/carlescliment/html2pdf-service (PHP, wkhtmltopdf)
- https://github.com/spipu/html2pdf (PHP, self-implemented HTML rendering)
