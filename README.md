# go-lsbtw

a from-scratch reimplementation of `ls` in Go, built step by step in the video. the
goal isn't a production tool, it's to understand what `ls` actually does: how it
reads a directory and how a single `lstat` call gives you everything in an `ls -l` row.

you should be able to take this video and implement more features to `ls`, and other coreutils.

## build / run

```sh
go run ls.go # run without compiling 
go build -o go-lsbtw ls.go # compile
 
./go-lsbtw            # list the current directory
./go-lsbtw -a         # include dotfiles
./go-lsbtw -l         # long format
./go-lsbtw -la /etc   # combine flags, list a path
```

## what's implemented

- **`os.ReadDir`** — walk a directory and print its entries. the error is the second
  return value, and `if error != nil` is the block you'll write most in go.
- **`-a`** — show hidden entries (the ones starting with `.`), plus `.` and `..`,
  which `ReadDir` filters out for you.
- **`-l`** — the long listing, where the real work is. For each entry we `os.Lstat`
  it and reach past the portable `FileInfo` into the raw inode with
  `file_info.Sys().(*syscall.Stat_t)`:
  - `Mode` -> file type + permission string (`-rw-r--r--`) via `mode_string`,
    `syscall.S_IFMT` and `os.FileMode.Perm().String()`
  - `Nlink` -> hard-link count
  - `Uid` / `Gid` -> owner and group names via `user.LookupId` / `user.LookupGroupId`
  - `file_info.Size()` -> size in bytes
  - `file_info.ModTime()` -> modification time, formatted with `time.Format("Jan _2 15:04")`
  - `os.Readlink()` -> real file path for links, when `syscall.S_IFLNK`

## where we left off

the output is correct field-for-field, but it differs from real `ls` in how
   padding is determined

**padding / alignment.** the link count and size columns are padded naively.
   real `ls` right-aligns them into nice columns, based off the largest number in each column.

sorting you get for free: `os.ReadDir` already returns entries sorted by filename,
which is the part we had to write by hand in the C video.

if you followed the video, you've got almost everything you need to add.

## contributing

feel free to contribute any hand written ls args, or specifically the padding and
   that was displayed at the end of the video.

cheers.

tony, btw
