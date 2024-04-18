# checkcidr

A CLI to check that IP addresses are included in CIDR.

## Usage

```bash
checkcidr <CIDR list file> <IP addresses list file>...

# ex:
checkcidr testdata/cidr_1.txt testdata/ip_1.txt
```

`checkcidr` prints result as free text style.
You can choose a format from `free_text`, `json` or `json_stream` (ex: `-style json`).

```bash
checkcidr -style json <CIDR list file> <IP addresses list file>...
```

`checkcidr` prints progress to stderr by default. Please set `-noprogress` if you don't need it.

```bash
checkcidr -noprogress <CIDR list file> <IP addresses list file>...
```

## License

MIT
