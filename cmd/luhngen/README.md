# luhngen

Generate numbers that pass the [Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm) check.

## Options

* `--length`: number of digits
* `--prefix`: first few digits (if you need a specific BIN or card type)
* `--cardtype`: if you want a specific card type.  Options are `V` (Visa), `M` (Mastercard), `D` (Discover) and `A` (American Express).  Will set the length and prefix.
* `--seed`: seed for deterministic pseudo-random output
* `--trailing-newline`: if a trailing newline should be emitted
