
package main

import (
	"strings"
	"fmt"
	"time"
)
const(
	numWordCounters = 2
	sleepInterval	= 5 * time.Second
)
var textStrings = []string{
	"this is the first string one two three four four three",
	"this is the second string one two three five one five one",
	"this is the third string six one two two two two",
	"the big brown fox jumped over the lazy dogs back",
	"the cow jumped over the moon",
}

type Text struct {
	str string
}

type Counts struct {
	wordcount map[string]int
}

func updateWordCount( wordcount map[string] int, count Counts){

	for word, v := range count.wordcount {
		    _, ok := wordcount[word]
			if ok {
				wordcount[word] += v
			} else {
				wordcount[word] = v
			}
			fmt.Println(word, " ", wordcount[word])
	}
}

func StateMonitor() chan <- Counts {
	counts := make(chan Counts)
	wordCount := make(map[string]int)
	go func() {
		for {
			select {
			case s := <-counts:
				updateWordCount(wordCount, s)
			}
		}
	}()

	return counts
}

func (t *Text) Sleep ( pending chan <- *Text){
	time.Sleep(sleepInterval)
	pending <- t
}

// WordCount returns a map of the counts of each “word” in the string s.
func WordCount(pending <-chan *Text, complete chan <- *Text, counts chan <-  Counts)  {
	for s := range pending {
		words := strings.Fields(s.str)
		countMap := make(map[string]int)
		for _, word := range words {
			_, ok := countMap[word]
			if ok {
				countMap[word]++
			} else {
				countMap[word] = 1
			}
		}
		counts <- Counts{countMap}
		complete <- s
	}
}

func main() {


	pending, complete := make(chan *Text,1), make(chan *Text,1)
	counts := StateMonitor()


	for i := 0; i < numWordCounters; i++ {
		go WordCount(pending, complete, counts)
	}


	go func() {
       for _, str := range textStrings{
       	pending <- &Text{str: str}
	   }
	}()

	for r := range complete {
		go r.Sleep(pending)
	}

}