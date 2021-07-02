package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func MD5Sum(root string) {
	m, err := MD5All(root)
	if err != nil {
		fmt.Printf("MD5All err:%v\n", err)
		return
	}

	var paths []string
	for path := range m {
		paths = append(paths, path)
	}

	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x %s\n", m[path], path)
	}
}

func MD5All(root string) (map[string][md5.Size]byte, error) {
	out := make(map[string][md5.Size]byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stage1: walk the file tree
	paths, errc := walkFiles(ctx, root)

	// Stage2: read and sum md5
	const numDigester = 20
	c := boundedDigester(ctx, paths, numDigester)

	// Stage3: collect file md5
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		out[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return out, nil
}

func walkFiles(ctx context.Context, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}

			select {
			case paths <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	}()
	return paths, errc
}

func digester(ctx context.Context, paths <-chan string, c chan<- result) {
	for path := range paths {
		data, err := ioutil.ReadFile(path)
		select {
		case c <- result{path, md5.Sum(data), err}:
		case <-ctx.Done():
			return
		}
	}
}

func boundedDigester(ctx context.Context, paths <-chan string, numDigester int) <-chan result {
	c := make(chan result)
	var wg sync.WaitGroup
	wg.Add(numDigester)
	for i := 0; i < numDigester; i++ {
		go func() {
			defer wg.Done()
			digester(ctx, paths, c)
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	return c
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	MD5Sum(root)
}
