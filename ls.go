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

	for _, name := range names {
		if long_format {
			print_long(path, name)
		} else {
			fmt.Println(name)
		}
	}
}

func mode_string(mode uint32) string {
	var s [10]byte
	switch mode & syscall.S_IFMT {
	case syscall.S_IFDIR:
		s[0] = 'd'
	case syscall.S_IFLNK:
		s[0] = 'l'
	case syscall.S_IFCHR:
		s[0] = 'c'
	case syscall.S_IFBLK:
		s[0] = 'b'
	case syscall.S_IFIFO:
		s[0] = 'p'
	case syscall.S_IFSOCK:
		s[0] = 's'
	default:
		s[0] = '-'
	}
	return string(s[:1]) + os.FileMode(mode).Perm().String()[1:]

}

func print_long(dir, name string) {
	path := filepath.Join(dir, name)
	file_info, error := os.Lstat(path)
	if error != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, error)
		return
	}
	file_stats := file_info.Sys().(*syscall.Stat_t)
	owner, group := "?", "?"
	if u, error := user.LookupId(strconv.Itoa(int(file_stats.Uid))); error == nil {
		owner = u.Username
	}
	if g, error := user.LookupGroupId(strconv.Itoa(int(file_stats.Gid))); error == nil {
		group = g.Name
	}
	if file_stats.Mode & syscall.S_IFLNK != 0 {
		if real_path, error := os.Readlink(path); error == nil {
			name = name + " -> " + real_path
		}
	}
	fmt.Printf("%s %2d %s %s %10d %s %s\n",
		mode_string(file_stats.Mode),
		file_stats.Nlink,
		owner,
		group,
		file_info.Size(),
		file_info.ModTime().Format("Jan _2 15:04"),
		name,
	)
}
