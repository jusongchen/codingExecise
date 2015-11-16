package main

import (
	"encoding/xml"
	"fmt"

	"io/ioutil"
)

type Result struct {
	Dataset Datasets
}

//Stat store statistics data for one dataset
//GenOutput create files to show result, XML file only right now
func (ds *Datasets) GenOutput(param *Param) error {

	//create a working directory which is a subdirectory with time stamp
	//func MkdirAll(path string, perm FileMode) error
	filename := fmt.Sprintf("%s/%s", param.workingDir, param.OutputFilename)

	// fmt.Printf("\n\nafter Init:\n%#v", ds)

	rs := Result{
		*ds,
	}
	output, err := xml.MarshalIndent(rs, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Printf("\nResult file:%s", filename)
	// func WriteFile(filename string, data []byte, perm os.FileMode) error
	return ioutil.WriteFile(filename, output, 0777)

}
