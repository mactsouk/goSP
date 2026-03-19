// gotree: a tree(1)-like utility in Go
//
// Features (subset of tree(1)):
//
//	-a              : All files, include entries beginning with '.'
//	-d              : List directories only
//	-L <n>          : Max display depth (1=just the named directories' entries)
//	--dirsfirst     : List directories before files
//	-f              : Print full path for each file
//	-p              : Print permissions (mode)
//	-s              : Print file size
//	-h              : Human-readable sizes (use with -s)
//	--noreport      : Do not print the summary report at the end
//	--version       : Print version and exit
//
// Notes:
//
//	  Symlinks are not followed; they are displayed as “name -> target”.
//	  Hidden entries are skipped unless -a is provided.
//	  Errors reading directories are reported to stderr and traversal continues.
//	  On Windows, box-drawing characters may require a Unicode-capable console.
//
//		If no path is provided, '.' is used.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

const version = "0.1.0"

// options holds CLI options.
type options struct {
	all       bool
	dirsOnly  bool
	level     int
	dirsFirst bool
	fullPath  bool
	showPerm  bool
	showSize  bool
	human     bool
	noReport  bool
}

type stats struct {
	dirs  int
	files int
	bytes int64
}

func main() {
	var opt options

	flag.BoolVar(&opt.all, "a", false, "all files, include entries beginning with '.'")
	flag.BoolVar(&opt.dirsOnly, "d", false, "list directories only")
	flag.IntVar(&opt.level, "L", -1, "max display depth")
	flag.BoolVar(&opt.dirsFirst, "dirsfirst", false, "list directories before files")
	flag.BoolVar(&opt.fullPath, "f", false, "print the full path prefix for each file")
	flag.BoolVar(&opt.showPerm, "p", false, "print permissions (mode)")
	flag.BoolVar(&opt.showSize, "s", false, "print the size of each file in bytes")
	flag.BoolVar(&opt.human, "h", false, "with -s, print human readable sizes")
	flag.BoolVar(&opt.noReport, "noreport", false, "turn off file/directory count at end of tree listing")

	ver := flag.Bool("version", false, "print version and exit")

	// Allow both -h (human) and --help.
	help := flag.Bool("help", false, "show help")

	// Capture the default Usage to provide extended help.
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "gotree %s — a tree(1)-like utility in Go\n\n", version)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] [path ...]\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "Examples:")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -a -p -s -h --dirsfirst -L 2\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -d -L 1 /usr/local\n", filepath.Base(os.Args[0]))
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}
	if *ver {
		fmt.Println(version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var total stats
	exitCode := 0

	for i, root := range args {
		st, err := doTree(root, opt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gotree: %v\n", err)
			exitCode = 1
		}
		if i > 0 {
			fmt.Println()
		}
		if !opt.noReport {
			total.dirs += st.dirs
			total.files += st.files
			total.bytes += st.bytes
		}
	}

	if !opt.noReport && len(args) == 1 {
		fmt.Printf("\n%d directories, %d files", total.dirs, total.files)
		if opt.showSize {
			if opt.human {
				fmt.Printf(", %s total\n", humanizeBytes(total.bytes))
			} else {
				fmt.Printf(", %d bytes total\n", total.bytes)
			}
		} else {
			fmt.Println()
		}
	}

	os.Exit(exitCode)
}

func doTree(root string, opt options) (stats, error) {
	var st stats

	// Resolve root printing name
	displayRoot := root
	if !opt.fullPath {
		if base, err := filepath.Abs(root); err == nil {
			// Show just the base name for the root when not full path
			displayRoot = filepath.Base(base)
		} else {
			displayRoot = filepath.Base(root)
		}
	} else {
		if abs, err := filepath.Abs(root); err == nil {
			displayRoot = abs
		}
	}

	fmt.Println(displayRoot)

	// If level is 0, do not list children.
	if opt.level == 0 {
		return st, nil
	}

	// We Lstat the root to decide if it's a file or directory.
	info, err := os.Lstat(root)
	if err != nil {
		return st, err
	}

	if !info.IsDir() {
		// Single file: print details if requested and return.
		line := formatEntry(root, info, "", true, opt)
		if line != "" {
			fmt.Println(line)
		}

		st.files++
		st.bytes += info.Size()
		return st, nil
	}

	err = walkDir(root, "", 1, opt, &st)
	return st, err
}

func walkDir(root, prefix string, depth int, opt options, st *stats) error {
	entries, err := readDir(root, opt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gotree: %s: %v\n", root, err)
		return nil // continue
	}

	for i, entry := range entries {
		last := i == len(entries)-1
		// Build the line prefix (├──/└──) and the next recursive prefix.
		branch := "├── "
		nextPrefix := prefix + "│   "
		if last {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		// Fetch FileInfo via Lstat to avoid following symlinks.
		full := filepath.Join(root, entry.Name())
		info, err := os.Lstat(full)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gotree: %s: %v\n", full, err)
			continue
		}

		// Skip hidden entries unless -a
		if !opt.all && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Print the line for this entry
		line := prefix + branch + formatEntry(full, info, root, false, opt)
		if line != "" {
			fmt.Println(line)
		}

		// Update stats
		if info.IsDir() {
			st.dirs++
		} else {
			st.files++
			st.bytes += info.Size()
		}

		// Recurse into directories if depth allows
		if info.IsDir() {
			if opt.level < 0 || depth < opt.level {
				if err := walkDir(full, nextPrefix, depth+1, opt, st); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// readDir returns directory entries according to options, sorted as requested.
func readDir(path string, opt options) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Filter hidden unless -a; filter files if -d
	filtered := entries[:0]
	for _, e := range entries {
		name := e.Name()
		if !opt.all && strings.HasPrefix(name, ".") {
			continue
		}
		if opt.dirsOnly {
			if !e.IsDir() {
				continue
			}
		}
		filtered = append(filtered, e)
	}

	// Sort
	sort.Slice(filtered, func(i, j int) bool {
		ei, ej := filtered[i], filtered[j]
		// Optionally list directories first
		if opt.dirsFirst {
			if ei.IsDir() != ej.IsDir() {
				return ei.IsDir() && !ej.IsDir()
			}
		}
		// Case-insensitive name order
		return strings.ToLower(ei.Name()) < strings.ToLower(ej.Name())
	})

	return filtered, nil
}

// formatEntry returns the printable line for a file or directory, without the tree branch prefix.
// If isRootFile is true, we are printing a single file root (rare case).
func formatEntry(fullPath string, info fs.FileInfo, root string, isRootFile bool, opt options) string {
	name := info.Name()

	// Determine display name per -f
	display := name
	if opt.fullPath {
		var base string
		if isRootFile {
			base = ""
		} else {
			base = root
		}
		p := fullPath
		if base != "" {
			// Normalize absolute
			if abs, err := filepath.Abs(fullPath); err == nil {
				p = abs
			}
		}
		display = p
	}

	// Identify symlink and potential target
	if info.Mode()&os.ModeSymlink != 0 {
		if target, err := os.Readlink(fullPath); err == nil {
			display = fmt.Sprintf("%s -> %s", display, target)
		}
	}

	var meta []string
	if opt.showPerm {
		meta = append(meta, modeString(info.Mode()))
	}
	if opt.showSize {
		if opt.human {
			meta = append(meta, humanizeBytes(info.Size()))
		} else {
			meta = append(meta, fmt.Sprintf("%d", info.Size()))
		}
	}

	if len(meta) > 0 {
		return fmt.Sprintf("%s [%s]", display, strings.Join(meta, " "))
	}
	return display
}

// modeString renders a file mode similar to ls(1), e.g., drwxr-xr-x.
func modeString(m fs.FileMode) string {
	var b strings.Builder

	// File type
	switch {
	case m.IsDir():
		b.WriteByte('d')
	case m&os.ModeSymlink != 0:
		b.WriteByte('l')
	case m&os.ModeSocket != 0:
		b.WriteByte('s')
	case m&os.ModeNamedPipe != 0:
		b.WriteByte('p')
	case m&os.ModeDevice != 0:
		if m&os.ModeCharDevice != 0 {
			b.WriteByte('c')
		} else {
			b.WriteByte('b')
		}
	default:
		b.WriteByte('-')
	}

	// Permissions
	rwx := []fs.FileMode{0400, 0200, 0100, 0040, 0020, 0010, 0004, 0002, 0001}
	for i, bit := range rwx {
		if m&bit != 0 {
			switch i % 3 {
			case 0:
				b.WriteByte('r')
			case 1:
				b.WriteByte('w')
			case 2:
				b.WriteByte('x')
			}
		} else {
			b.WriteByte('-')
		}
	}

	// Setuid, setgid, sticky adjustments
	if m&os.ModeSetuid != 0 {
		// user execute position (index 3)
		replaceAt(&b, 3, 's', 'S')
	}
	if m&os.ModeSetgid != 0 {
		// group execute position (index 6)
		replaceAt(&b, 6, 's', 'S')
	}
	if m&os.ModeSticky != 0 {
		// others execute position (index 9)
		replaceAt(&b, 9, 't', 'T')
	}

	return b.String()
}

func replaceAt(b *strings.Builder, idx int, withExec, withoutExec rune) {
	s := b.String()
	if idx >= len(s) {
		return
	}
	r := []rune(s)
	if r[idx] == 'x' {
		r[idx] = withExec
	} else if r[idx] == '-' {
		r[idx] = withoutExec
	}
	b.Reset()
	b.WriteString(string(r))
}

func humanizeBytes(n int64) string {
	if n < 0 {
		return fmt.Sprintf("%d", n)
	}
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := int64(unit), 0
	for m := n / unit; m >= unit; m /= unit {
		div *= unit
		exp++
	}
	prefixes := []string{"K", "M", "G", "T", "P", "E"}
	if exp >= len(prefixes) {
		exp = len(prefixes) - 1
	}
	return fmt.Sprintf("%.1f%siB", float64(n)/float64(div), prefixes[exp])
}

// Platform-specific helper: try to detect same file (device+inode) to avoid potential cycles if symlink following is added later.
// Not used for traversal since we do not follow symlinks, but kept for completeness.
func sameFile(fi1, fi2 fs.FileInfo) bool {
	if fi1 == nil || fi2 == nil {
		return false
	}
	s1, ok1 := fi1.Sys().(*syscall.Stat_t)
	s2, ok2 := fi2.Sys().(*syscall.Stat_t)
	if ok1 && ok2 {
		return s1.Dev == s2.Dev && s1.Ino == s2.Ino
	}
	// Fallback: compare names & sizes if Stat_t unavailable
	return fi1.Name() == fi2.Name() && fi1.Size() == fi2.Size() && fi1.Mode() == fi2.Mode()
}

// Guard to avoid unused import error on syscall on platforms where Stat_t isn't available.
var _ = errors.New
