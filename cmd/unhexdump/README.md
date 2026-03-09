# unhexdump

Convert `hexdump` output back into a binary file

## Description

If you don't have a full-blown binary file editor, you can use this to change bytes in a file.

It takes the output of `hexdump -c` (or the [`hexdumpc` utility](../hexdumpc/)) and creates a binary file.

It only looks at the middle section (the hex bytes).
