/*copyright 2015 Jusong Chen.
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/jusongchen/stat"
)

//Stat store statistics data for one dataset
type Dataset struct {
	File     string    //the filename
	Count    int       //total count of numbers  in this dataset
	TopN     []float64 //top 10 values
	Histo    []stat.HistoBin
	binWidth float64 //Bin width, unexposed
}

//Datesets keep a list of Dataset
type Datasets []Dataset

const (
	TOP_N = 10
	// DELIMITOR          = ';'
	DELIMITOR                       = ' ' //Do not change this as this version support space seperated only
	DATAFILE_EXTENSION              = ".DAT"
	MAX_FILES                       = 10000
	MAX_NUMBERS                     = 100000
	MAX_VALUE                       = 1.0e+9
	MIN_VALUE                       = 0.0
	TOO_MANY_BINS_WARNING_THRESHOLD = 10000
)

type Param struct {
	NFile          int     //int is good enough as the max is 10,000
	NItem          int     //int is good enough as the max is 100,000
	MinVal         float64 //default to 0.0, constraint >=0.0
	MaxVal         float64 //no default,  constraits : <=1.0e+9, >minVal
	BinWidth       float64 //no default, constraits : >0.0, <= (maxVal - minVal)
	OutputFormat   string  //default to xml
	OutputFilename string
	workingDir     string //the directory to create data files
}

var paramLists = map[string]string{
	"nFile":    fmt.Sprintf("Number of files to generate,must be in the range of [1,%d]", MAX_FILES),
	"nNumber":  fmt.Sprintf("Number of Numbers in each file,must be in the range of [1,%d]", MAX_NUMBERS),
	"minVal":   fmt.Sprintf("Minimum value of numbers,must be in the range of [0,%e]", MAX_VALUE),
	"maxVal":   fmt.Sprintf("Maximum value of numbers,must be in the range of [0,%e]", MAX_VALUE),
	"binWidth": fmt.Sprintf("Histogram bin width,must be in the range of (0,%e]", MAX_VALUE),
}

var (
	// The program will take the following parameters
	nFile    = flag.Int("nFile", 1, paramLists["nFile"])
	nNumber  = flag.Int("nNumber", 20, paramLists["nNumber"])
	minVal   = flag.Float64("minVal", 0, paramLists["minVal"])
	maxVal   = flag.Float64("maxVal", 10, paramLists["maxVal"])
	binWidth = flag.Float64("binWidth", 1.0, paramLists["binWidth"])

	dop = flag.Int("DOP", 4*runtime.NumCPU(), " Max Degree of Parallelism")
)

func getParams() (Param, bool) {

	valid := true
	// check if passed-in parameters are valid or not
	switch {
	case *nFile < 1 || *nFile > MAX_FILES:
		fmt.Fprintf(os.Stderr, "Invalid Parameter:%s", paramLists["nFile"])
		valid = false

	case *nNumber < 1 || *nNumber > MAX_NUMBERS:
		fmt.Fprintf(os.Stderr, "Invalid Parameter:%s", paramLists["nNumber"])
		valid = false

	case *minVal < MIN_VALUE || *minVal > MAX_VALUE:
		fmt.Fprintf(os.Stderr, paramLists["minVal"])
		valid = false
	case *maxVal < MIN_VALUE || *maxVal > MAX_VALUE:
		fmt.Fprintf(os.Stderr, paramLists["maxVal"])
		valid = false
	case *minVal >= *maxVal:
		fmt.Fprintf(os.Stderr, "minVal must be less than maxVal")
		valid = false
	case *binWidth < MIN_VALUE || *binWidth > MAX_VALUE:
		fmt.Fprintf(os.Stderr, paramLists["binWidth"])
		valid = false
	}
	if !valid {
		return Param{}, valid
	}

	param := Param{
		NFile:          *nFile,
		NItem:          *nNumber,
		MinVal:         *minVal,
		MaxVal:         *maxVal,
		BinWidth:       *binWidth,
		OutputFormat:   "xml",
		OutputFilename: "result.xml",
	}
	return param, true

}

// var dbConn *mssql.Conn

func main() {
	flag.Parse() //   SetupDB()

	if len(os.Args) == 1 {
		//expecting at least one parameter to be set
		// fmt.Printf("For help run:\n%s --help ", os.Args[0])
		fmt.Printf(`Usage Example:
%s -nFile=3 -nNumber=20 -binWidth=0.5 -maxVal=50`, os.Args[0])

		fmt.Println("\nDefault parameter settings:")
		flag.PrintDefaults()
		return
	}

	//start timing
	now := time.Now()
	param, valid := getParams()
	if !valid {
		return
	}
	/*
		param := Param{
			NFile:          3,
			NItem:          2000,
			MinVal:         20.780,
			MaxVal:         1.0e+2 + 0.4,
			BinWidth:       2.5,
			OutputFormat:   "xml",
			OutputFilename: "result.xml",
		}
	*/

	//warning if binWidth is too small
	if (param.MaxVal-param.MinVal)/param.BinWidth > TOO_MANY_BINS_WARNING_THRESHOLD {
		fmt.Println(`
The bin width is relative small and this may generate large amount of histogram data
Press Y to continue, anther other to abort ...`)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		c := []byte(input)[0]
		if !(c == 'y' || c == 'Y') {
			return
		}

	}

	fmt.Printf("%v:Start working now . . .\n", now)
	//seed random number generater
	rand.Seed(time.Now().Unix())

	datasets := Datasets{}

	//Init create files and populate numbers to each file
	if err := datasets.Init(&param); err != nil {
		//OS error ecountered
		return
	}

	ds := datasets.Analyze()

	ds.GenOutput(&param)
	fmt.Printf("\n\n%v: Time consumed to complete all jobs:%v\n", time.Now(), time.Since(now))

}
