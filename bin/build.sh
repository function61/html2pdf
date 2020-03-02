#!/bin/bash -eu

source /build-common.sh

BINARY_NAME="html2pdf"
COMPILE_IN_DIRECTORY="cmd/html2pdf"

function maybeDownloadWkhtmlToPdf {
	if [ -f wkhtmltopdf ]; then
		return # already downloaded
	fi

	heading "Downloading wkhtmltopdf"

	apt install -y xz-utils

	# correct sha-256 is 840e3b30668af203dc685986db7ace92a5495d70d0c76ba13045b2aad24e201b
	curl --fail --location https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.4/wkhtmltox-0.12.4_linux-generic-amd64.tar.xz \
		| xz -d \
		| tar --strip-components=2 -xf - wkhtmltox/bin/wkhtmltopdf
}

# TODO: one deployerspec is done, we can stop overriding this from base image
function packageLambdaFunction {
	if [ ! -z ${FASTBUILD+x} ]; then return; fi

	cd rel/
	cp "${BINARY_NAME}_linux-amd64" "${BINARY_NAME}"
	rm -f lambdafunc.zip
	zip -j lambdafunc.zip "${BINARY_NAME}" "../wkhtmltopdf"
	rm "${BINARY_NAME}"
}

maybeDownloadWkhtmlToPdf

standardBuildProcess

packageLambdaFunction
