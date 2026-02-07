# go-analyzer

A small Go command-line tool that scans a directory tree and reports files larger than a specified size.

## Usage

- `-d` : directory to start scanning from (default `.`)
- `-s` : minimum file size in GB to report (default `1`)

## Examples

Run against the current directory with default thresholds:

```powershell
go run .
```

Run against a specific directory and a 1.5 GB threshold:

```powershell
go run . -d path\to\scan -s 1.5
```

Run in the current directory and report files above ~100KB (for testing):

```powershell
go run . -s 0.0001
```

## Notes

- The `-s` value is in gigabytes (GB). Use fractional values for sizes smaller than 1GB.
- The program prints lines like `Large file: <path> | Size: <x> GB` for matches.

See the `main.go` source for implementation details.
