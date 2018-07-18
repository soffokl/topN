package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/wangjohn/quickselect"
)

type topN struct {
	N          int
	Multiplier int
	buff       []int
}

func (t *topN) Add(n int) error {
	t.buff = append(t.buff, n)

	if len(t.buff) > t.Multiplier*t.N {
		return t.reduce()
	}

	return nil
}

func (t *topN) reduce() error {
	chunk := quickselect.IntSlice(t.buff)
	reverse := quickselect.Reverse(chunk)
	if err := quickselect.QuickSelect(reverse, t.N); err != nil {
		return err
	}

	t.buff = chunk[:t.N]
	return nil
}

func (t *topN) Result() ([]int, error) {
	if err := t.reduce(); err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.IntSlice(t.buff[:t.N])))

	return t.buff[:t.N], nil
}

func main() {
	N := flag.Int("n", 10, "number of biggest numbers to show")
	m := flag.Int("m", 1000, "multiplier for the n value. It sets how many memory could be the tool use as a buffer")
	file := flag.String("f", "", "name of the source data file")

	flag.Parse()

	input := os.Stdin
	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		input = f
	}

	res, err := search(input, *N, *m)
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range res {
		fmt.Println(n)
	}
}

func search(input io.Reader, n, m int) ([]int, error) {
	r := bufio.NewReader(input)
	top := topN{N: n, Multiplier: m}
	for {
		s, _, err := r.ReadLine() // assuming that a line will be less than a buffer, ignoring isPrefix here
		if err != nil {
			break
		}

		n, err := strconv.Atoi(string(s))
		if err != nil {
			return nil, err
		}

		if err := top.Add(n); err != nil {
			return nil, err
		}
	}

	return top.Result()
}
