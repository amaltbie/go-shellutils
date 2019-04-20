package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

var (
	aFlagVar bool
	lFlagVar bool
	oneFlagVar bool
	fileList []string
)

func usage() {
	fmt.Printf("usage: ls [-1al] [file ...]")
	flag.VisitAll(func (f *flag.Flag){
		fmt.Printf("  -%s\t%s\n", f.Name, f.Usage)
	})
}

func init() {
	flag.Usage = usage
	flag.BoolVar(&aFlagVar, "a", false, "Include directory entries whose names begin with a dot (.).")
	flag.BoolVar(&lFlagVar, "l", false, "(The lowercase letter ``ell''.)  List in long format.  (See below.)  If the output is to a terminal, a total sum for all the file sizes is output on a line before the long listing.")
	flag.BoolVar(&oneFlagVar, "1", false, "(The numeric digit ``one''.)  Force output to be one entry per line.  This is the default when output is not to a terminal.")
	flag.Parse()
	if flag.NArg() == 0 {
		fileList = []string{"."}
	} else {
		fileList = flag.Args()
	}
}

func printFiles(files []os.FileInfo) {
	fileSep := "\t"
	// Print files line by line if -l or -1 are used
	if oneFlagVar || lFlagVar {
		fileSep = "\n"
	}

	for _, file := range files {
		// Only print files that begin with '.' if -a is present
		if !aFlagVar && file.Name()[0] == '.' {
			continue
		}

		// By default only print the file name
		// If -l is present then print extended info
		fileInfo := file.Name()
		if lFlagVar {
			// Get system-specific stat info
			stat := file.Sys().(*syscall.Stat_t)

			// Get username and groupname but use id if no name is found
			uidStr := strconv.Itoa(int(stat.Uid))
			owner, err := user.LookupId(uidStr)
			ownerName := ""
			if err == nil {
				ownerName = owner.Username
			} else {
				ownerName = uidStr
			}
			gidStr := strconv.Itoa(int(stat.Gid))
			group, err := user.LookupGroupId(gidStr)
			groupName := ""
			if err == nil {
				groupName = group.Name
			} else {
				groupName = gidStr
			}

			fileInfo = fmt.Sprintf("%10s%5d%6s%6s%10d%15s %s",file.Mode(),stat.Nlink,ownerName,groupName,file.Size(),file.ModTime().Format("Jan 2 15:05"),file.Name())
		}
		fmt.Printf("%s%s",fileInfo,fileSep)
	}
	// Print a newline if all files are on oneline
	if fileSep == "\t" {
		fmt.Printf("\n")
	}
}

func main() {
	normalFiles := []os.FileInfo{}
	dirs := map[string][]os.FileInfo{}
	for _, file := range fileList {
		fileinfo, err := os.Stat(file)
		if err != nil {
			pErr := err.(*os.PathError)
			fmt.Printf("ls: %s: %s\n",pErr.Path,pErr.Err)
			continue
		}
		if fileinfo.IsDir() {
			dirFiles, _ := ioutil.ReadDir(file)
			dirs[file] = dirFiles

			// Prepend dir and parent dir if -a is used
			if aFlagVar {
				wd, _ := os.Getwd()
				os.Chdir(file)
				parDirStat, _ := os.Stat("..")
				dirs[file] = append([]os.FileInfo{parDirStat}, dirs[file]...)
				dirStat, _ := os.Stat(".")
				dirs[file] = append([]os.FileInfo{dirStat}, dirs[file]...)
				os.Chdir(wd)
			}
		} else {
			fileInfo, _ := os.Stat(file)
			normalFiles = append(normalFiles, fileInfo)
		}
	}
	if len(normalFiles) > 0 {
		printFiles(normalFiles)
	}
	for k, v := range dirs {
		if len(dirs) > 1 || len(normalFiles) > 0 {
			fmt.Printf("\n%s:\n", k)
		}
		printFiles(v)
	}
}
