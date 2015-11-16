package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

//genData create files ;  it populate each file by generating random float numbers to populate
func (ds *Datasets) Init(param *Param) error {

	//create a working directory which is a subdirectory with time stamp
	//func MkdirAll(path string, perm FileMode) error
	param.workingDir = fmt.Sprintf("./%s", strings.Replace(time.Now().Format(time.RFC3339), ":", "", -1)[:17])
	//the above generate a working dir liek this one: ./2015-11-15T102932

	// fmt.Printf("\n\n%s", param.workingDir)

	//check error
	if err := os.MkdirAll(param.workingDir, 0777); err != nil {
		log.Print(err)
		return err
	}

	//create data files and populate
	for i := 0; i < param.NFile; i++ {

		//create file
		filename := fmt.Sprintf("%s/%06d%s", param.workingDir, i+1, DATAFILE_EXTENSION)
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalln(err)
		}
		d := Dataset{
			File:     filename,
			Count:    param.NItem,
			binWidth: param.BinWidth,
			// TopN:     []float64{},
			// histo:    []stat.HistoBin{},
		}

		defer f.Close()

		//populate data items, items seperated by comma

		w := bufio.NewWriter(f)

		for j := 0; j < param.NItem; j++ {
			//_, errW := f.Write([]byte(fmt.Sprintf("%f%c", random(param.MinVal, param.MaxVal), DELIMITOR)))
			_, errW := fmt.Fprintf(w, "%f%c", random(param.MinVal, param.MaxVal), DELIMITOR)
			if errW != nil {
				log.Fatalln(errW)
			}
			if j%2048 == 0 {
				w.Flush() //write buffer to disk
			}
		}

		w.Flush() // make sure data in buffer write to destination file

		//add new dataset to the datasets
		*ds = append(*ds, d)
		// ds[i].File =
		fmt.Printf("\nFile %s created and populated with %d numbers", d.File, param.NItem)
	}
	// fmt.Printf("\n\nafter Init:\n%#v", ds)
	return nil
}

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
