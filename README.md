# Diff/Patch

>! This is work in progress! Currently the code is ugly, hacky, and probably too loud at parties.
>! It's the kind of code where the commit message would be an apology for writing it.
>! Any part of the API including the patch format may change on minor or patch versions pre 1.0.
>! Check back later if you want to use this in production.


diffpatch is a diff and patch algorithm with Go and Javascript implementations.

The aim is a diff algorithm that is fast in practice, produces small patch sizes, with patches that can be applied across different systems, clients, etc.

## Features

**Interoperable**

The patches produced by Go and JS are interoperable and serializable.
That is to say you could generate a patch client side with JS, send it to your backend written in Go (or Node of course), and be able to arrive at the correct result.

**Fast**

In real-world usage, this algorithm is very efficient.
In terms of big-O, Myers' has a better worst case than this implementation.

**Efficient**

This algorithm outputs small patches (without guaranteeing they are optimally small)

**No dependencies**

Both the Go and Javascript implementations have zero external dependencies

**Tiny**

I actually need to measure exactly how small the packages end up being. For now you'll have to take my word for it that they're itty bitty. Perhaps even teeny weeny.

## Input format

The diff functions take two arrays of strings as input. These are just the texts you are diffing split by something that's easy to join them on after.
I find for natural language text in latin scripts, splitting on '.' works really well. This gives a nice balance between speed and resulting patch size.
Splitting on new lines `\n` is a reliable option too. The way you split your text will impact speed and patch size, so it's worth checking what works for you.

## Patch format

diffpatch doesn't produce patches that are nice to read, unless you're a computer program - in which case you'll love reading them. If you're a human though, not so much. This helps keep the size down for sending over network/storing/hanging in a frame above your fireplace.

Patches are an array of operations. Each operation is an index, the number of items to delete, and the items to insert.

An example patch:
```json
[
  {"i":[23, 10], "a":[]},
  {"i":[48, 1], "a":["\n\nSuspendisse ullamcorper molestie"]},
  {
    "i":[54, 1],
    "a":[
      " Duis faucibus, quam vel tristique mattis, leo odio volutpat tortoise, ac molestie purus erat ac felis"
    ]
  },
  {
    "i":[61, 1],
    "a":[
      "\nDonec bibendum, erat eu consequat vehicula, elit justo lacinia tortor, at vehicula leo felis et mi"
    ]
  },
  {"i":[64, 0], "a":[" It was the best of days, it was the blurst of days"]}
]
```

## Focus

diffpatch is primarily focused on 'natural' text, such as documents, blog posts, markdown files, etc.
The algorithm doesn't particularly _care_ what the text is, but when optimising I test against those cases rather than (e.g.) code or structured data.

