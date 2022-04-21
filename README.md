# leindex

Easily find a given little endian integer in a string based on a range.

This can be used to locate timestamps.

## Example

	res := leindex.IndexLE32(buf, 1546300800, 1672531200)

As for `bytes.Index` this method returns -1 if the value was not found.
