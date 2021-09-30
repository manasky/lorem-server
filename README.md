# Lorem.space server
The backend of https://lorem.space lives here.

## Features
* Scan the directory & categorize the image resources (only JPEG for now)
* Provide HTTP API
* Resize the image resources
* Cache option to store the resized image as a file

## How to run

`$ go run main.go`

`$ go run main.go --dir="/IMAGE/DIRECTORY/PATH"`

### Options

* `host`: host:port for the HTTP server (string)
* `dir`: the path of image resources directory (string)
* `cache`: enable cache (store the result as files in .cache directory) (bool)
* `cdn`: CDN address to redirect to the cached file path, leave empty to write the file in the HTTP response (string)
* `min-width`: minimum supported width (integer)
* `max-width`: maximum supported width (integer)
* `min-height`: minimum supported height (integer)
* `max-height`: maximum supported height (integer)

### Docker
You can pull & run the public docker image from github package registry:

`$ docker pull ghcr.io/manasky/lorem-server:latest`
