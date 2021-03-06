package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// find represents the find command
var find = &cobra.Command{
	Use:   "find",
	Short: "find the occurence of any word or phrase throughout the entire project",
	Long: `given a word or phrase find where in the project is located what you're 
lookig for. You'll get a list of files in which matches were found including their 
line number.

For example:

term find --path /home/user/todo --word "debugger" --ext ".jsx" `,
	Run: func(cmd *cobra.Command, args []string) {

		path, _ := cmd.Flags().GetString("path")

		for a := range args {
			fmt.Println(a)
		}

		word, _ := cmd.Flags().GetString("word")

		// If icase flag is set to true transform word to lower case
		icase, _ := cmd.Flags().GetBool("icase")
		if icase {
			word = strings.ToLower(word)
		}

		// In case the ext flag is set
		ext, _ := cmd.Flags().GetString("ext")

		// Walk walks the file tree rooted at root, calling walkFn for each file or directory in the tree,
		// including root. All errors that arise visiting files and directories are filtered by walkFn.
		// The files are walked in lexical order, which makes the output deterministic but means that for
		// very large directories Walk can be inefficient. Walk does not follow symbolic links.

		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if len(ext) > 0 && filepath.Ext(path) != ext {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				panic(err)
			}

			// Scanner provides a convenient interface for reading data such as a file of newline-delimited
			// lines of text. Successive calls to the Scan method will step through the 'tokens' of a file,
			// skipping the bytes between the tokens. The specification of a token is defined by a split
			// function of type SplitFunc; the default split function breaks the input into lines with line
			// termination stripped. Split functions are defined in this package for scanning a file into
			// lines, bytes, UTF-8-encoded runes, and space-delimited words. The client may instead provide
			// a custom split function.

			// Scanning stops unrecoverably at EOF, the first I/O error, or a token too large to fit in the
			// buffer. When a scan stops, the reader may have advanced arbitrarily far past the last token.
			// Programs that need more control over error handling or large tokens, or must run sequential
			// scans on a reader, should use bufio.Reader instead.

			scanner := bufio.NewScanner(file)

			var line int = 0
			var b []byte

			for scanner.Scan() {
				// Bytes returns the most recent token generated by a call to Scan.
				if icase {
					b = bytes.ToLower(scanner.Bytes())
				} else {
					b = scanner.Bytes()
				}

				line++
				if bytes.Contains(b, []byte(word)) {
					fmt.Println(path, " ln: "+strconv.Itoa(line))
				}
			}

			return nil
		})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(find)
	find.Flags().String("path", "", "--path [directory]")
	find.Flags().String("word", "", "--word [word]")
	find.Flags().BoolP("icase", "i", false, "all lower and upper case words")
	find.Flags().String("ext", "", "--ext [.json]")
}
