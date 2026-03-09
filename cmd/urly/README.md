# urly

Get a URL into the specific form required by a program without regular expressions or custom parsers.

## Usage

The URL can be passed on the command line or in an environment value specified with `--url-env`.

The component (see table) can be set with the corresponding flag, or deleted with `-no-XXX` flag.

```
┌─────────┬───┬──────────┬──────────┬─┬────────────────┬─┬──────┬──────────┬─┬──────────────┬─┬──────────┐
│ scheme  │   │ username │ password │ │   hostname     │ │ port │  path    │ │    query     │ │ fragment │
│  https   ://    user   :   pass    @  sub.example.com :  8080   /p/a/t/h  ?  query=string  #  hash     │
└────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

Password input is different: for security reasons, it should not be passed on the command line.  You can pass
it in via stdin with `--password-stdin` or set it via the environment with `--password-env`.  Username can also
be passed via the environment with `--username-env`.

Query string parameters are a little different: they can be set with `--setparam` or added with `--addparam` or
deleted with `--delparam`.

The output can be the complete URL or a specific component with the `--output` flag.  You can also output the
complete parsed URL as formatted JSON with `--output=json` or a single line JSON with `--output=jsonl`.

## Exit Status

The program will exit with an errorlevel if it cannot parse the URL.

## Examples

Combine username and password from the environment to make a complete URL:
```
urly --username-env=DB_USER --password-env DB_PSWD "postgres://pgdbhost.example.com/db?sslMode=required"
```

## Standards

[RFC 3986](https://datatracker.ietf.org/doc/html/rfc3986): Uniform Resource Identifier (URI): Generic Syntax

## Security Considerations

While urly does not accept passwords on the command line, it will output URLs with a password if specified.  Be careful with this output!
