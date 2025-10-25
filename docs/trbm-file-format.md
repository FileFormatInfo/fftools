# Trilobyte mapping file format

Yaml file with the following structure:

* `handle`: short unique ID for this mapping (usually the same as the filename)
* `title`: one-line title for this mapping
* `description`: detailed description (markdown)
* `default`: default for bytes that don't have a map.  Can be `identity` or a byte array.
* `map`: map of bytes -> byte array
