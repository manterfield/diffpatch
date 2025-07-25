# Diff/Patch

diffpatch is a diff and patch algorithm with Go and Javascript implementations.

The aim is a diff algorithm that is fast in practice, produces small patch sizes, with patches that can be applied across different systems, clients, etc.

## Features

**Interoperable**
The patches produced by Go and JS are interoperable and serializable.
That is to say you could generate a patch client side with JS, send it to your backend written in Go (or Node of course), and be able to arrive at the correct result.

**Fast**
In terms of big-O, Myers' has a better worst case than this implementation.
In real-world usage however, this algorithm is very efficient.

**Efficient**
This algorithm outputs small patches (without guaranteeing they are optimally small)

**No dependencies**
Both the Go and Javascript implementations have zero external dependencies

**Tiny**
I actually need to measure exactly how small the packages end up being. For now you'll have to take my word for it that they're itty bitty. Perhaps even teeny weeny.


## Notes

> diffpatch is primarily focused on 'natural' text, such as documents, blog posts, markdown files, etc.
> The algorithm doesn't particularly _care_ what the text is, but when optimising I test against those cases rather than (e.g.) code or structured data.

> diffpatch doesn't produce patches that are nice to read, unless you're a computer program - in which case you'll love reading them. If you're a human though, not so much. This helps keep the size down for sending over network/storing/etc.
