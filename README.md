# Cron Parser

## Build

    go build ./cmd/parser

## Run

    Usage: ./parser [cron line...]

## Output

    ./parser '*/15 0 1,15 * Mon /usr/bin/find'

    minute         0 15 30 45
    hour           0
    day of month   1 15
    month          1 2 3 4 5 6 7 8 9 10 11 12
    day of week    1
    command        /usr/bin/find

## License

MIT.
