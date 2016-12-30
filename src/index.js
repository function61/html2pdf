var child_process = require('child_process')
  , fs = require('fs')
  , express = require('express')
  , bodyParser = require('body-parser')
  , multipart = require('connect-multiparty')
  , promClient = require('prom-client');

// set up instrumentation
var html2pdf_requests_total = new promClient.Counter('html2pdf_requests_total', 'Total requests');
var html2pdf_render_duration_ms = new promClient.Histogram('html2pdf_render_duration_ms', 'wkhtml2pdf render duration', {
    buckets: [ 10, 100, 150, 200, 300, 400 ]
});

var app = express();

app.use(bodyParser.urlencoded({ // to support URL-encoded bodies
	extended: true
}));

var multipartMiddleware = multipart();

var uniqueCounter = 0;

// "1_17468917", "2_73679142", ...
function semiUniqueId() {
	return (++uniqueCounter) + '_' + Math.round(Math.random() * 99999999).toString();
}

function createFifoReadStream(next) {
	var fifoFilePath = '/tmp/html2pdf_fifo_' + semiUniqueId();
	var mkfifoProcess = child_process.spawn('mkfifo',  [ fifoFilePath ]);
	mkfifoProcess.on('exit', function (exit_code) {
		if (exit_code !== 0) {
			throw new Error('fail to create fifo with exit_code:  ' + exit_code);
		}

		var fifoReadStream = fs.createReadStream(fifoFilePath);

		next(fifoReadStream, fifoFilePath);
	});
}

function renderDocument(req, res) {
	html2pdf_requests_total.inc();

	if (!('document' in req.files)) {
		throw new Error('upload expected');
	}

	// TODO: rename file because wkhtmltopdf requires the filename to end in .html
	// (currently body-parser & connect-multiparty preserve uploaded filename but it's not ok to trust it)
	//
	// looks like "/tmp/kjWbRzvaECslQ5BbWKfXeBKZ.html"
	var document_path = req.files.document.path;

	/*	we need to create a FIFO because for some stupid reason wkhtmltopdf fails to write to stdout when
		invoked from node.js (shell invokation and stdout writing seems to be ok)

		file would also work but we use FIFO to bypass unnecessary disk I/O.

		if performance becomes an issue we could easily implement FIFO pool as well.
		
		Many people seem to be having the same issue:

		https://github.com/wkhtmltopdf/wkhtmltopdf/issues/2036#issuecomment-230863803 (my own report incl. strace output)
		http://www.perlmonks.org/index.pl?node_id=1085338
		https://github.com/nodejs/node-v0.x-archive/issues/3974 (might be related to node.js itself since with PHP we could use passthru() )
	*/
	createFifoReadStream(function (fifoReadStream, fifoFilePath){
		// for A4 this would be "page_size_mm=210x297"
		var page_size_mm = /^(\d+)x(\d+)$/.exec(req.body.page_size_mm);

		var args = [ ];

		if (page_size_mm) {
			args.push('--page-width', page_size_mm[1]);
			args.push('--page-height', page_size_mm[2]);
		}
		else {
			args.push('--page-size', 'A4');
		}

		var allArgs = args.concat([
			// '--quiet',
			document_path, // "-" could be given here to signify stdin
			fifoFilePath // see above bug report - otherwise we'd use "-" here to mean stdout
		]);

		console.log('command line:', 'wktmltopdf', allArgs.join(' '));

		var startedMs = new Date().getTime();

		var wkhtmltopdf = child_process.spawn('wkhtmltopdf', allArgs);

		var stderr = '';

		// need this to know if we'll dump either the STDERR or the PDF to HTTP response
		var waitingFirstPdfChunk = true;

		// if wkhtmltopdf supported writing to STDOUT:
		//     wkhtmltopdf.stdout.on('data', function (chunk){
		fifoReadStream.on('data', function (chunk){
			if (waitingFirstPdfChunk) {
				waitingFirstPdfChunk = false;

				res.setHeader('Content-Type', 'application/pdf');
			}

			res.write(chunk);
		});

		wkhtmltopdf.stderr.on('data', function (chunk){
			stderr += chunk.toString();
		});

		wkhtmltopdf.on('close', function (exit_code){
			var duration = new Date().getTime() - startedMs;

			// prom-client also has utility: var end = metric.startTimer(); ... end();
			html2pdf_render_duration_ms.observe(duration);

			console.log('wkhtmltopdf finished in ' + duration + 'ms exit_code=' + exit_code);

			fs.unlink(document_path, function (err){ if (err) { throw err; } });
			fs.unlink(fifoFilePath,  function (err){ if (err) { throw err; } });

			if (waitingFirstPdfChunk) { // error occurred => send stderr instead
				console.log('wkhtmltopdf stderr', stderr);

				res.write(stderr);
				res.end();
			}
			else {
				res.end();
			}
		});
	});
}

app.post('/render', multipartMiddleware, renderDocument);
app.get('/metrics', function (req, res){
	res.send(promClient.register.metrics());
});

var server = app.listen(80, function () {
  var host = server.address().address;
  var port = server.address().port;

  console.log('html2pdf server listening at port %s', port);
});
