Steps to build and run the program
==================================
a) Install Go version 1.4.2 or above

b) After Go install, set GOPATH env
#check your GOPATH environment values
echo $GOPATH

#create these two directories

mkdir -p $GOPATH/src/github.com/jusongchen/codingExercise
mkdir -p $GOPATH/src/github.com/jusongchen/stat

#The source code is delivered as two packages: main and stat. see notes below for reasons

c) Copy the following files to $GOPATH/src/github.com/jusongchen/stat
histogram.go
histogram_test.go

d) Run the following cmd to test package "stat" (this step is optional)

cd $GOPATH/src/github.com/jusongchen/stat
go test

e) Copy the following files to $GOPATH/src/github.com/jusongchen/codingExercise
analysis.go
dataGen.go
main.go
output.go

f) Build the binary

cd $GOPATH/src/github.com/jusongchen/codingExercise
go build

The above build cmd shoud generate an executable 

g) run the program

This is an example
codingExecise -nFile=80 -nNumber=20000 -binWidth=2.5 -maxVal=50

Note: This program takes 6 parameters. As the above cmd does not pass in a value for parameter minVal, the default value (0.0) is used for the parameter.

Run the command without any parameter will give Usage inforation:

codingExecise
Usage Example:
codingExecise.exe -nFile=3 -nNumber=20 -binWidth=0.5 -maxVal=50
Default parameter settings:
  -DOP=16:  Max Degree of Parallelism (not used in this version)
  -binWidth=1: Histogram bin width,must be in the range of (0,1.000000e+09]
  -maxVal=10: Maximum value of numbers,must be in the range of [0,1.000000e+09]
  -minVal=0: Minimum value of numbers,must be in the range of [0,1.000000e+09]
  -nFile=1: Number of files to generate,must be in the range of [1,10000]
  -nNumber=20: Number of Numbers in each file,must be in the range of [1,100000]


Working directory and results
=============================
The program automatically creates a directory within the current directory to store data files and the result file. The result file is in XML format right now.

Design Notes:
============
1) Although this is only an exercise, any another project can reuse code in package stat. Here is the signiture of the key function in this package:


//HistoAnalyze calls a Scanner to get tokens and convert tokens to numbers
//then,it do analysis to get Top N and Histogram of those numbers
//it returns an error if an token read in cannot be converted to
//a float64 number or the scanner encounter an error
//
// When there is no error detected, it returns
// 		1) number of tokens processed
// 		2) top N numbers
// 		3) histogram data
//		4) nil
func HistoAnalyze(s Scanner, N int, binWidth float64) (int64, []float64, []HistoBin, error) {


//the first parameter of the function is an interface which is defined as below:


// Scanner is the interface that wraps basic Scan methods.
type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}

This design allows this function to get data feed from different data sources(e.g. file,network,database, etc.) 
as long as the data source implements the Scanner interface. This design also makes this function testable(code in file histogram_test.go demostrates this) .

2) histogram_test.go contains test code for function HistoAnalyze().The table driven test method is used to cover several basic cases. Testing is time consuming. Microsoft Excel is used to generate some of test data and their expected results. Given the limited amount I can spend on this project, not all test cases are covered.