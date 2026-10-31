package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var (
	nlink_width, owner_width, group_width, fsize_width, mtime_width int

	modes = map[uint32]rune{
		syscall.S_IFDIR:  'd',
		syscall.S_IFLNK:  'l',
		syscall.S_IFCHR:  'c',
		syscall.S_IFBLK:  'b',
		syscall.S_IFIFO:  'p',
		syscall.S_IFSOCK: 's',
	}
)

func main() {
	path := "."
	show_all, long_format := false, false

	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			path = arg
			continue
		}
		// -l, -la, -a
		for _, char := range arg[1:] {
			switch char {
			case 'a':
				show_all = true
			case 'l':
				long_format = true
			default:
				fmt.Fprintf(os.Stderr, "usage: %s [-al] [path]\n", os.Args[0])
				os.Exit(1)
			}

		}
	}

	entries, error := os.ReadDir(path)
	if error != nil {
		fmt.Fprintln(os.Stderr, error)
		os.Exit(1)
	}

	var names []string
	if show_all {
		names = append(names, ".", "..")
	}

	for _, entry := range entries {
		if !show_all && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		names = append(names, entry.Name())
	}

	if long_format {
		set_widths(path, names)
		for _, name := range names {
			print_long(path, name)
		}
	} else {
		for _, name := range names {
			fmt.Println(name)
		}
	}
}

func mode_string(mode uint32) string {
	var s [10]byte
	if r, ok := modes[mode&syscall.S_IFMT]; ok {
		s[0] = byte(r)
	} else {
		s[0] = '-'
	}
	return string(s[:1]) + os.FileMode(mode).Perm().String()[1:]

}

func set_widths(dir string, names []string) {
	for _, name := range names {
		_, nlink, owner, group, size, mtime := get_long_info(dir, &name)
		nlink_width = max(nlink_width, len(strconv.Itoa(int(nlink))))
		owner_width = max(owner_width, len(owner))
		group_width = max(group_width, len(group))
		fsize_width = max(fsize_width, len(strconv.Itoa(int(size))))
		mtime_width = max(mtime_width, len(mtime))
	}
}

func print_long(dir, name string) {
	modestr, nlink, owner, group, size, mtime := get_long_info(dir, &name)

	fmt.Printf("%s %*d %*s %*s %*d %*s %s\n",
		modestr,
		nlink_width, nlink,
		owner_width, owner,
		group_width, group,
		fsize_width, size,
		mtime_width, mtime,
		name,
	)
}

func get_long_info(dir string, name *string) (modestr string, nlink uint64, owner, group string, size int64, mtime string) {
	fullpath := filepath.Join(dir, *name)

	file_info, error := os.Lstat(fullpath)
	if error != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", *name, error)
		return
	}
	file_stats := file_info.Sys().(*syscall.Stat_t)
	owner, group = "?", "?"
	if u, error := user.LookupId(strconv.Itoa(int(file_stats.Uid))); error == nil {
		owner = u.Username
	}
	if g, error := user.LookupGroupId(strconv.Itoa(int(file_stats.Gid))); error == nil {
		group = g.Name
	}

	if file_stats.Mode&syscall.S_IFMT == syscall.S_IFLNK {
		target, err := os.Readlink(fullpath)
		if err == nil {
			*name += " -> " + target
		}
	}

	return mode_string(file_stats.Mode), file_stats.Nlink, owner, group, file_info.Size(), file_info.ModTime().Format("Jan _2 15:04")
}
