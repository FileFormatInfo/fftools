# luhncheck

Checks if a number passes the [Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm) check

## Options

* `--no-error-level`: don't return an errorlevel if it doesn't pass (note: you can still get an errorlevel if the number is missing or cannot be parsed)
* `--verbose`: if it does not pass, print a message with the last digit fixed to pass
* `--quiet`: do not print `PASS`/`FAIL`
