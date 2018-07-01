# taon
[![Build Status](https://travis-ci.com/eiri/taon.svg?branch=master)](https://travis-ci.com/eiri/taon)

Transform JSON into ASCII table.

## Installation

`taon` is a stand-alone cli utility. You can just [download a binary](https://github.com/eiri/taon/releases) and run it.

Drop the binary in your `$PATH` (e.g. `~/bin`) to make it easy to use.

## Usage

Read JSON from a file:
```
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
```
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

## Help
```
$ taon --help
usage: taon [<flags>] [<file>]

Transform JSON into ASCII table.

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and --help-man).
      --version            Show application version.
  -c, --columns=COL1,COL2  List of columns to display

Args:
  [<file>]  File to read

```

## Licence

[MIT](https://github.com/eiri/taon/blob/master/LICENSE)
