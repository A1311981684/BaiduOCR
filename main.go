package main

import (
	"fmt"
	"log"
	"ocr/ocr"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("no picture assigned")
		return
	}
	fmt.Println("Processing, please wait...")
	//Get file name from command line argument
	path := os.Args[1]
	//DO OCR
	res, err := ocr.OCR(path)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//Prepare for write result
	outF, err := os.Create("result.txt")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		err := outF.Close()
		if err != nil {
			log.Println(err.Error())
			return
		}
	}()
	//Write result
	_, err = outF.Write([]byte("《----OCR Results----》\r\n\r\n"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, v := range res{
		for k1, v1 := range v {
			if k1 == "words" {
				_, err = outF.Write([]byte(fmt.Sprintf("%s"+"\r\n", v1)))
				if err != nil {
					log.Println(err.Error())
					return
				}
			}
		}
	}
	_, err = outF.Write([]byte("\r\n(----置信度(Probability)----)\r\n"))
	for _, v := range res{
		for k1, v1 := range v {
			if k1 == "probability" {
				_, err = outF.Write([]byte(fmt.Sprintf("%f"+"\r\n", v1.(map[string]interface{})["average"].(float64))))
				if err != nil {
					log.Println(err.Error())
					return
				}
			}
		}
	}
	fmt.Println("Result written to result.txt.")
}
