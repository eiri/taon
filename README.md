# taon
[![Build Status](https://travis-ci.com/eiri/taon.svg?branch=master)](https://travis-ci.com/eiri/taon)

Transform JSON into ASCII table.

## Installation

`taon` is a stand-alone cli utility. You can just [download a binary](https://github.com/eiri/taon/releases) and run it.

Drop the binary in your `$PATH` (e.g. `~/bin`) to make it easy to use.

## Usage

Read JSON from a file:
```bash
$ taon -c seq,name testdata/data.json
+-----+----------+
| seq |   name   |
+-----+----------+
| 1   | Donovan  |
| 2   | Timothy  |
| 3   | Nici     |
| 4   | Nigel    |
| 5   | Saqib    |
| 6   | Turkeer  |
| 7   | Damone   |
| 8   | Mick     |
| 9   | Theodore |
| 10  | Hsuan    |
| 11  | Ramneek  |
| 12  | Roderick |
+-----+----------+
```

Pass JSON from cURL output:
```bash
$ curl -s https://raw.githubusercontent.com/eiri/taon/master/testdata/array.json | taon -c seq,word,bool
+-----+---------------------+-------+
| seq |        word         | bool  |
+-----+---------------------+-------+
| 1   | electrophototherapy | false |
| 2   | twatterlight        | true  |
| 3   | phlebograph         | false |
| 4   | Ervipiame           | false |
| 5   | annexational        | false |
| 6   | unjewel             | true  |
| 7   | Anglic              | true  |
| 8   | alliable            | true  |
| 9   | seraphism           | true  |
| 10  | congenialize        | true  |
| 11  | phu                 | false |
| 12  | vial                | false |
+-----+---------------------+-------+
```

Print table as markdown
```bash
$ taon -c number,string -m testdata/data.json
| number |      string      |
|--------|------------------|
| 779    | 7Xf6cUtJEOPjAbEc |
| 700    | M2s3HKnr5zoWxAdd |
| 310    | xZzrKV5XIL1P9y9H |
| 742    | mfEefyltzS1lbfje |
| 841    | X4bjUqiAUhYZvNvD |
| 352    | ixF1I79VqoFyKFPx |
| 852    | BYTHmkHRtI9e48K9 |
| 818    | 6K3YjMZ7bzUrJ6kt |
| 822    | 96M1TNPDN3WugPuZ |
| 822    | PwbXirV2qj2vlK6g |
| 549    | EKfWADxgJ7obe1w9 |
| 777    | ex9esRTklAKofF8B |
```

## Help
```bash
$ taon --help
usage: taon [<flags>] [<file>]

Transform JSON into ASCII table.

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and --help-man).
      --version            Show application version.
  -c, --columns=COL1,COL2  List of columns to display
  -m, --markdown           Print as markdown table

Args:
  [<file>]  File to read
```

## Licence

[MIT](https://github.com/eiri/taon/blob/master/LICENSE)
