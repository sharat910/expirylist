# Expiry List
A (small) golang package that essentially implements a linked list that is sorted on time and 
provides apis to use in conjunction with a standard Golang `map` to add an expiry feature i.e. entries in 
the map can be expired if not accessed (updated) within a `timeout`.

For example usage, refer to `cmd/stringcache`.