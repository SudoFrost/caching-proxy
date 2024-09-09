# Caching Proxy

This is a caching proxy that can proxy to an origin server and cache its responses.

## installation

1. Make sure you have Go installed: https://golang.org/doc/install
2. Clone this repository: `git clone https://github.com/sudoforst/caching-proxy.git`
3. Navigate to the `caching-proxy` directory
4. Build the binary: `go build`
5. Run the binary: `./caching-proxy`

## usage

Proxying to an origin server:
```bash
./caching-proxy -p 8080 -o http://example.com
```

Clearing the cache:
```bash
./caching-proxy -c
```

##Project Page

Visit the project page for more information: [Project URL](https://roadmap.sh/projects/caching-server)
