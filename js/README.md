# Diff/Patch

>! This is work in progress! Currently the code is ugly, hacky, and probably too loud at parties.
>! It's the kind of code where the commit message would be an apology for writing it.
>! Any part of the API including the patch format may change on minor or patch versions pre 1.0.
>! Check back later if you want to use this in production.


@mantty/diffpatch is a fast, efficient diff and patch algorithm.

The aim is a diff algorithm that is fast in practice, produces small patch sizes, with patches that can be applied across different systems, clients, etc.

It has a corresponding go implementation, which produce interoperable patches.

[READ MORE](https://github.com/manterfield/diffpatch)
