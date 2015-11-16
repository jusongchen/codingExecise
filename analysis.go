package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jusongchen/stat"
)

//Analyze get all datasets and analyze each of them
func (ds *Datasets) Analyze() *Datasets {
	// fmt.Printf("\n\n%#v", ds)

	destDS := Datasets{}
	for _, d := range *ds {
		d.analyze()
		destDS = append(destDS, d)
	}
	return &destDS
}

//analyze one file
func (d *Dataset) analyze() error {

	//first open file
	// f,err := os.Open(d.File)
	// fmt.Printf("\n%v:Analyzing file %s", time.Now(), d.File)
	fmt.Printf("\nAnalyzing file %s", d.File)

	f, err := os.Open(d.File)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	//initialize the histo map

	scanner := bufio.NewScanner(f)
	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)

	// func HistoAnalyze(s Scanner, N int, binWidth float64) (int64, []float64, []HistoBin, error) {
	count, topN, histo, err := stat.HistoAnalyze(scanner, TOP_N, d.binWidth)

	//number
	if count != int64(d.Count) {
		log.Fatalf("\nExpecting %d items in file %s, but get %d items", d.Count, d.File, count)
	}

	if err != nil {
		log.Fatal(err)
	} else {
		//pass result to Dataset
		// copy(d.TopN, topN)
		d.TopN = topN
		d.Histo = histo

	}
	return nil

}
