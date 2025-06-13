package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/alexmullins/zip"
)

var (
	zipFile    string
	dictionary string
	minLength  int
	maxLength  int
	charset    string
	workers    int
	verbose    bool
)

func init() {
	flag.StringVar(&zipFile, "f", "", "ZIP file to crack (required)")
	flag.StringVar(&dictionary, "d", "", "Dictionary file for attack")
	flag.IntVar(&minLength, "min", 1, "Minimum password length (for brute force)")
	flag.IntVar(&maxLength, "max", 8, "Maximum password length (for brute force)")
	flag.StringVar(&charset, "c", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "Character set for brute force")
	flag.IntVar(&workers, "w", runtime.NumCPU(), "Number of workers")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  Dictionary attack: zipcracker -f secret.zip -d passwords.txt")
		fmt.Fprintln(os.Stderr, "  Brute force attack: zipcracker -f secret.zip -min 4 -max 6")
	}
}

func main() {
	flag.Parse()

	if zipFile == "" {
		fmt.Println("Error: ZIP file must be specified")
		flag.Usage()
		os.Exit(1)
	}

	start := time.Now()
	var password string
	var found bool

	if dictionary != "" {
		fmt.Printf("Starting dictionary attack on %s...\n", zipFile)
		password, found = dictionaryAttack(zipFile, dictionary)
	} else {
		fmt.Printf("Starting brute force attack on %s (length %d-%d)...\n", zipFile, minLength, maxLength)
		password, found = bruteForceAttack(zipFile, minLength, maxLength, charset)
	}

	if found {
		fmt.Printf("\nPassword found: %s\n", password)
		fmt.Printf("Time elapsed: %v\n", time.Since(start))
	} else {
		fmt.Println("\nPassword not found")
		os.Exit(1)
	}
}

func dictionaryAttack(zipPath, dictPath string) (string, bool) {
	file, err := os.Open(dictPath)
	if err != nil {
		fmt.Printf("Error opening dictionary file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	passwords := make(chan string, 1000)
	results := make(chan string, 1)
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go zipWorker(zipPath, passwords, results, done, &wg)
	}

	// Feed passwords to workers
	go func() {
		for scanner.Scan() {
			select {
			case <-done:
				file.Close()
				return
			default:
				passwords <- scanner.Text()
			}
		}
		close(passwords)
	}()

	// Wait for result or completion
	select {
	case password := <-results:
		close(done)
		return password, true
	case <-waitGroupDone(&wg):
		return "", false
	}
}

func bruteForceAttack(zipPath string, minLen, maxLen int, charset string) (string, bool) {
	passwords := make(chan string, 1000)
	results := make(chan string, 1)
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go zipWorker(zipPath, passwords, results, done, &wg)
	}

	// Generate passwords
	go func() {
		generatePasswords("", minLen, maxLen, charset, passwords, done)
		close(passwords)
	}()

	// Wait for result or completion
	select {
	case password := <-results:
		close(done)
		return password, true
	case <-waitGroupDone(&wg):
		return "", false
	}
}

func generatePasswords(prefix string, minLen, maxLen int, charset string, passwords chan<- string, done <-chan struct{}) {
	if len(prefix) >= minLen {
		select {
		case passwords <- prefix:
			if verbose && len(prefix) == minLen {
				fmt.Printf("\rTesting length %d...", len(prefix))
			}
		case <-done:
			return
		}
	}

	if len(prefix) >= maxLen {
		return
	}

	for _, c := range charset {
		select {
		case <-done:
			return
		default:
			generatePasswords(prefix+string(c), minLen, maxLen, charset, passwords, done)
		}
	}
}

func zipWorker(zipPath string, passwords <-chan string, results chan<- string, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// Open the ZIP file once per worker
	zf, err := zip.OpenReader(zipPath)
	if err != nil {
		if verbose {
			fmt.Printf("Error opening ZIP file: %v\n", err)
		}
		return
	}
	defer zf.Close()

	// We'll try to decrypt the first file in the archive
	if len(zf.File) == 0 {
		return
	}

	for password := range passwords {
		select {
		case <-done:
			return
		default:
			for _, file := range zf.File {
				file.SetPassword(password)
				r, err := file.Open()
				if err == nil {
					// Password worked!
					io.Copy(io.Discard, r) // We need to read some data to verify the password
					r.Close()
					select {
					case results <- password:
					case <-done:
					}
					return
				}
				if verbose {
					fmt.Printf("\rTesting: %s", password)
				}
			}
		}
	}
}

func waitGroupDone(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

