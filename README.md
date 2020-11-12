# Genatree

Genatree is a small go tool that generates very fastly a huge directory and file tree.

## Installation

On \*nix platform (`/` path separator):
```console
go install github.com/Jorropo/genatree
```
For performance reasons it's cheaper to just concatenate with `/` the path, that mean it might not works on windows.

## What is is usefull for ?

- Testing stuff, that the reason I made it, I was benchmarking [IPFS](https://jorropo.ovh/ipns/ipfs.io/) and I was needing a huge tree.
- Eating all the inodes of your fs, yes on linux system you need some inodes to describe anything, the number of available inodes is set at creating of the partition (I guess some tools can change it) and without some left you can't create any thing on your FS even if you have some free space left.

That pretty much all.

## Performance

Theses tests are done on a tmpfs (ram disk) with a bandwith of aproximately 10Gb/s (just to say, bottleneck is not there).

In shell (I guess you can do better with forking) :
```sh
$ time for i in $(echo t/{a,b,c}/{a,b,c}/{a,b,c}/{a,b,c}/{a,b,c}/{a,b,c}/{a,b,c}/{a,b,c}/{a,b,c}/{a,b,c}); do mkdir -p $i && for j in $(echo {a,b,c}); do echo "$i/$j" > $i/$j ; done; done

real	1m40,137s
```
Genatree :
```console
$ time genatree -count=10
Creating 177147 files.
Done !

real	0m0,416s
```
Globaly genatree is very heavely paralelised and scales great accross lots of cores, that mean he is not twice as long to create twice many files.

It can use all my inodes in only 11s while `rm` takes 40s to delete this same structure.

## Usage

```console
$ genatree --help
Usage of genatree:
  -count uint
    	depth of folders to create (default 3)
  -fdcap uint
    	maximum number of concurrently opened files (default 1024)
  -root string
    	first directory to be created (default "t")
  -set string
    	csv set to use at each step (default "a,b,c")
```

## Output

```console
$ genatree -count 2 -set "a,b" && tree
Creating 8 files.
Done !
.
└── t
    ├── a
    │   ├── a
    │   │   ├── a
    │   │   └── b
    │   └── b
    │       ├── a
    │       └── b
    └── b
        ├── a
        │   ├── a
        │   └── b
        └── b
            ├── a
            └── b

7 directories, 8 files
```

Each file contains the path to itself from the point genatree started :
```console
$ cat t/a/b/a 
t/a/b/a
```
