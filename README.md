# go-edge

`go-edge` provides the command `edge` to show the edge of logs with conditional grep

## Usage

    $ edge [options] FILE

## Options

    -g   --grep  KEYWORD  grep condition
    -c                    also show the total line count

## Examples

Here is the 'error.log' file.

    $ cat error.log
    2023/02/27 08:43:00,not found
    2023/02/27 08:43:01,not found
    2023/02/27 08:43:02,not found
    2023/02/27 08:43:03,no auth

Then you execute the `edge` command. Then shown the first line and last line.

    $ edge error.log
    1: 2014/06/27 08:43:00,not found
    4: 2014/06/27 08:43:03,no auth

The number of top of line is line number of the file.

And you can use `-c` option to know total line number.

    $ edge -c error.log
    1: 2014/06/27 08:43:00,not found
    4: 2014/06/27 08:43:03,no auth
    total 4 lines

The above is pretty much the same as below command executions:

    $ cat error.log | head -n1
    $ cat error.log | tail -n1
    $ cat error.log | wc -l

With `--grep` option, you can get filtered 1st line and last line.

    $ edge --grep found error.log
    1: 2023/02/27 08:43:00,not found
    3: 2023/02/27 08:43:02,not found

With `--grep` and `-c` option:

    $ edge -c --grep found error.log
    1: 2023/02/27 08:43:00,not found
    3: 2023/02/27 08:43:02,not found
    total 4 lines

## Installation

    go install github.com/bayashi/go-edge/cmd/edge@latest

## License

Apache License 2.0

## Author

Dai Okabayashi: @bayashi
