## nana dockerfile args

List the ARG declarations of a Dockerfile with their defaults

```
nana dockerfile args [flags]
```

### Examples

```
  nana dockerfile args
  nana dockerfile args --file build/Dockerfile -o json
```

### Options

```
  -f, --file string     path to the Dockerfile (default "Dockerfile")
  -h, --help            help for args
  -o, --output string   output format: table or json (default "table")
```

### SEE ALSO

* [nana dockerfile](nana_dockerfile.md)	 - Read facts from a Dockerfile

