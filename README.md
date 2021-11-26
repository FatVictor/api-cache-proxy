# API CACHE PROXY

To avoid rate limit threshold from API provider. This proxy will capture the request and cache the response for period of time.

## Requirements

- Golang 1.17+
- Docker
- Redis

## Run app

Environment variables

| Name        | Description |
| ----------- | ----------- |
| TARGET_HOST | 3rd party target. Eg. https://api.3rd-party.com       |
| REDIS_CONNECTION | Redis connection string. Format: *redis://user:password@localhost:6789/3?dial_timeout=3&db=1&read_timeout=6s&max_retries=2*|
| CACHE_TTL | Caching periods. Default is 60 seconds|
|SHOW_REQUEST_LOG|true/false|

Listen port: 3000

Docker image can be found https://hub.docker.com/r/fatvictor/api-cache-proxy

## Reporting bugs and contributing

- Want to report a bug or request a feature? Please open [an issue](https://github.com/FatVictor/api-cache-proxy/issues/new).
- Want to help us build project? Please submit your pull request
