# taon
[![CI Status](https://github.com/eiri/taon/actions/workflows/test.yaml/badge.svg)](https://github.com/eiri/taon/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/eiri/taon)](https://goreportcard.com/report/github.com/eiri/taon)

Transform JSON into ASCII table.

## Installation

`taon` is a stand-alone cli utility. You can just [download a binary](https://github.com/eiri/taon/releases) and run it.

Drop the binary in your `$PATH` (e.g. `~/bin`) to make it easy to use.

## Usage

### Reading

Read JSON from a file:

```bash
$ taon pkg/taon/testdata/example.json
+-------+-------+--------+-----+---------------------+
| bool  | count | number | seq | word                |
+-------+-------+--------+-----+---------------------+
| false | 0001  | 779    | 1   | electrophototherapy |
| true  | 0002  | 700    | 2   | twatterlight        |
| false | 0003  | 310    | 3   | phlebograph         |
| false | 0004  | 742    | 4   | Ervipiame           |
| false | 0005  | 841    | 5   | annexational        |
| true  | 0006  | 352    | 6   | unjewel             |
| true  | 0007  | 852    | 7   | Anglic              |
| true  | 0008  | 818    | 8   | alliable            |
| true  | 0009  | 822    | 9   | seraphism           |
| true  | 0010  | 822    | 10  | congenialize        |
| false | 0011  | 549    | 11  | phu                 |
| false | 0012  | 777    | 12  | vial                |
+-------+-------+--------+-----+---------------------+
```

or from the cURL output:

```bash
$ curl -s https://github.com/eiri/taon/blob/main/pkg/taon/testdata/example.json | taon
+-------+-------+--------+-----+---------------------+
| bool  | count | number | seq | word                |
+-------+-------+--------+-----+---------------------+
| false | 0001  | 779    | 1   | electrophototherapy |
| true  | 0002  | 700    | 2   | twatterlight        |
| false | 0003  | 310    | 3   | phlebograph         |
| false | 0004  | 742    | 4   | Ervipiame           |
| false | 0005  | 841    | 5   | annexational        |
| true  | 0006  | 352    | 6   | unjewel             |
| true  | 0007  | 852    | 7   | Anglic              |
| true  | 0008  | 818    | 8   | alliable            |
| true  | 0009  | 822    | 9   | seraphism           |
| true  | 0010  | 822    | 10  | congenialize        |
| false | 0011  | 549    | 11  | phu                 |
| false | 0012  | 777    | 12  | vial                |
+-------+-------+--------+-----+---------------------+
```

_Note: By default `taon` sorts columns alphabetically by a name to preserve order's stability. To explicitly define the order of the columns use `--columns` flag with comma separated list of the names_

### Filtering

Filter columns to a given list in the specified order:

```bash
taon -c seq,word,bool pkg/taon/testdata/example.json
+-----+---------------------+-------+
| seq | word                | bool  |
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

### Markdown

Print a table as markdown:

```bash
$ taon --columns seq,word,bool --markdown pkg/taon/testdata/example.json
| seq | word                | bool  |
|-----|---------------------|-------|
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
```

### Limitations

`taon` only works with objects or arrays of objects, passing in an array of arbitrary data will return an error.

## Help

```bash
$ taon --help
Transform JSON into ASCII table

Usage: taon [flags] [file]

Flags:
  -c, --columns=COL1,COL2  List of columns to display
  -m, --markdown           Print as markdown table
  -h, --help               Show help
      --version            Show application version

Args:
  <file>                   Path to file to read, stdin when missing
```

## Licence

[MIT](https://github.com/eiri/taon/blob/master/LICENSE)
