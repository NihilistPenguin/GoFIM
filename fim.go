// TODO: implement file deletion check

package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

// reading file by chunks of 2048 bytes
const chunksize = 2048

// function to get the hash of a given file
func getHash(file string) string {
	// open file, if error exits, log fatal
	f, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}

	// close file after this function closes
	defer f.Close()

	// get file size
	info, err := f.Stat()
	filesize := info.Size()

	// find the number of blocks we need, devide filesize by chunksize
	numBlocks := uint64(math.Ceil(float64(filesize) / float64(chunksize)))

	// instantiate md5 hash
	hash := md5.New()

	// for the number of blocks calculated
	for i := uint64(0); i < numBlocks; i++ {
		// choose lowest of chunksize or remaining filesize
		blocksize := uint64(math.Min(float64(chunksize), float64(filesize-int64(i*chunksize))))

		// buf creates a byte slice of size blocksize
		buf := make([]byte, blocksize)

		// read blocksize amount of bytes into buf
		f.Read(buf)
		// add bytes to our hash data
		io.WriteString(hash, string(buf))
	}

	// hash the data and return the hex encoded version
	return hex.EncodeToString(hash.Sum(nil))

}

// func checkFileDeletion(lookup map[string]string, fileInfo ) {
// 	for _, file := range fileInfo {
// 		fmt.Println(file)
// 	}
// }

func main() {
	// create map (like python dictionary) for our files:hashes
	lookup := make(map[string]string)

	// infinite loop, babyyyy
	for {
		// get everything in directory
		// if error, log fatal
		fileInfo, err := ioutil.ReadDir(".")

		if err != nil {
			log.Fatal(err)
		}

		// create a WaitGroup
		// this allows us to wait for all goroutines to finish before moving on
		var wg sync.WaitGroup

		// tell the wg we need (# of files) amount of goroutines
		wg.Add(len(fileInfo))

		// the first var in range is the index, which we don't need
		for _, file := range fileInfo {
			filename := file.Name()

			// create goroutine func
			// get hash of file, compare to map to see if changed
			go func(fname string) {
				// wait until function finishes to tell wg that this goroutine is done
				defer wg.Done()

				//check if file is a directory, contiue if not
				if info, _ := os.Stat(fname); info.IsDir() == false {
					hash := getHash(fname)

					// h is value of fname in the lookup map
					// ok is bool result of if the fname exists in map
					// so, if fname exists and hash doesn't match with hash
					// 		in lookup, then we know it changed
					if h, ok := lookup[fname]; ok && h != hash {
						fmt.Printf("%s\t%s has been changed!\n", time.Now().Format("01-02-2006 15:04:05"), filename)
					} else if _, ok := lookup[fname]; !ok {
						fmt.Printf("%s\t%s has been added!\n", time.Now().Format("01-02-2006 15:04:05"), filename)
					}

					// update lookup map with fname:hash
					lookup[fname] = hash
				}
			}(filename) // this calls the goroutine with this parameter
		}

		// wait until all specified # of goroutines finish
		wg.Wait()

		time.Sleep(1 * time.Second)
	}
}
