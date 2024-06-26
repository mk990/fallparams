# fallparams

This project is a Go application for extracting parameters from HTML pages and URLs.

## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

- Go installed on your local machine

### Installation

```sh
go install -v github.com/mk990/fallparams@latest
```
Here's a basic README.md file for your GitHub repository:

### Usage

```sh
fallparams [options] [url]
```

#### Options

- `-X`, `--method`: HTTP method (default is GET)
- `-i`, `--site-map`: URL list (optional)
- `-m`, `--chunk`: URL parameter count on output (default is 20)
- `-O`, `--output-format`: Output format (params or url, default is params)
- `-r`, `--response-file`: Response file (optional)
- `-s`, `--silent`: Silent mode (default is false)
- `-v`, `--verbose`: Verbose mode (default is false)
- `-u`, `--urlparams`: URL parameters only (default is false)
- `-c`, `--concurrency`: Concurrency (default is 20)
- `-l`, `--headless`: Headless mode (default is false)
- `-d`, `--debug`: Debug mode (default is false)

### Examples

- Extract parameters from a single URL:

```sh
fallparams http://example.com
```
- Extract parameters from a list of URLs:

```sh
cat urls.txt | fallparams
```

Extract parameters from a sitemap:

```sh
fallparams -i http://example.com/sitemap.xml
```

# Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.
- Fork the Project
- Create your Feature Branch (git checkout -b feature/AmazingFeature)
- Commit your Changes (git commit -m 'Add some AmazingFeature')
- Push to the Branch (git push origin feature/AmazingFeature)
- Open a Pull Request

## License

Distributed under the GNU General Public License v3.0. See `LICENSE` for more information.


