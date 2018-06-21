# taon
[![Build Status](https://travis-ci.org/eiri/taon.svg?branch=master)](https://travis-ci.org/eiri/taon)

Transform JSON into ASCII table.

## Installation

`taon` is a stand-alone cli utility. You can just [download a binary](https://github.com/eiri/taon/releases) and run it.

Drop the binary in your `$PATH` (e.g. `~/bin`) to make it easy to use.

## Usage

Pass JSON from a file:
```
TBD
```

Pass JSON from cURL output:
```
TBD
```

## Help
```
$ taon --help
Transform JSON into ASCII table.

Usage:
  taon [OPTIONS] [FILE|-]

Options:
  -c, --colorize   Colorize output
      --version    Print version information

Exit Codes:
  0 OK
  1 Failed to parse JSON
  2 Failed to open file
  3 Failed to read input
  4 Failed to render table

Examples:
  taon /tmp/somedata.json
  curl -s http://api.example.com/data/1 | taon
```

## Licence

[MIT](https://github.com/eiri/taon/blob/master/LICENSE)
