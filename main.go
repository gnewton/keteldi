/*
 * Basic example for text searching: Retrieving position of a signature line in PDF where the signature line is given by
 * "__________________" text. And positioned with a Tm operation above.
 *
 * Run as: go run pdf_detect_signature.go input.pdf
 */

package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	//"strings"
	//"io/ioutil"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

var allOps map[string]string

func main() {
	//log.SetOutput(ioutil.Discard)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) < 2 {
		log.Printf("Usage: keteldi input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	allOps = make(map[string]string)

	err := extract(inputPath)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	//fmt.Println(allOps)
}

func extract(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			log.Printf("Error: %v\n", err)
			return err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}

	for i := 0; i < numPages; i++ {
		log.Println("\n\nPage", i, " ############################################################################")
		log.Println("############################################################################")
		pageNum := i + 1

		fmt.Println("\nPage", i, "  %%%%%%%%%%%%%%%%%%%%%%%")
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return err
		}
		dict := page.GetPageDict()
		//fmt.Println("PageDict", dict)
		if page.Resources.Font != nil {
			log.Printf("PageFontDict%v \n", page.Resources.Font.String())
		}

		//keys := dict.Keys()
		//for key := range keys
		anotts := dict.Get("Annots")
		if anotts != nil {
			fmt.Println("ANNOTS", reflect.TypeOf(anotts))
			if ann, ok := anotts.(*pdfcore.PdfObjectArray); ok {

				fmt.Println("ANNOTS x", reflect.TypeOf(ann))
				for _, k := range *ann {
					fmt.Println("ANNOTS x z", reflect.TypeOf(k))
					fmt.Println("ANNOTS x z", k.String(), "==", k.DefaultWriteString())
					z := pdfcore.TraceToDirectObject(k)
					fmt.Println("ANNOTS x z g", reflect.TypeOf(z))
					fmt.Println("ANNOTS x z g", z)

				}

			}
		}

		//if fontDict, ok := page.Resources.Font.(*pdfcore.PdfObjectDictionary); ok {
		//log.Println("PageResourcesFontDict", fontDict)
		//font := fontDict.Get("TT2")
		//log.Println("@@@@@@@@@@", font)
		//log.Println("&&&", reflect.TypeOf(font))
		//if indo, ok := font.(*pdfcore.PdfIndirectObject); ok {
		//log.Println("@@@@@@@@@@*******", indo.DefaultWriteString())
		//log.Printf("%T\n", indo.PdfObject)
		//log.Printf("%s\n", indo.PdfObject.String())
		//if fdict, ok := indo.PdfObject.(*pdfcore.PdfObjectDictionary); ok {
		//log.Println("~~~~~~~~~", fdict)
		//}
		//}
		//}

		//fontDict := page.Resources.Font

		err = locateSignatureLine(page)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return err
		}

	}

	return nil
}

func locateSignatureLine(page *pdf.PdfPage) error {

	pageContentStr, err := page.GetAllContentStreams()
	if err != nil {
		return err
	}

	cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)
	if err != nil {
		return err
	}

	// text, err := cstreamParser.ExtractText()
	// fmt.Println("------------", text)
	// if true {
	// 	return nil
	// }

	operations, err := cstreamParser.Parse()
	if err != nil {
		return err
	}

	var ok bool
	var fontDict *pdfcore.PdfObjectDictionary
	if fontDict, ok = page.Resources.Font.(*pdfcore.PdfObjectDictionary); ok {

	}

	lastFont := ""
	for _, op := range *operations {

		//fmt.Println("=========", op.Operand, op.Params)
		allOps[op.Operand] = "m"
		fmt.Println(op.Operand, op.Params)
		switch op.Operand {
		case Tm_setTextMatrix:
			if len(op.Params) == 6 {
				if _, ok := op.Params[4].(*pdfcore.PdfObjectFloat); ok {
					//x = float64(*val)
				}

				if _, ok := op.Params[5].(*pdfcore.PdfObjectFloat); ok {
					//y = float64(*val)
				}
			}

		case Tj_showText:
			if len(op.Params) == 1 {
				val, ok := op.Params[0].(*pdfcore.PdfObjectString)
				if ok {
					str := string(*val)
					//if strings.Contains(str, "a") {
					log.Println(str)
					log.Printf("Tj: %s\n", str)
					fmt.Print(str)
					//break
					//}
				}
			}

		case Tf_setFont:
			if len(op.Params) == 2 {
				log.Println("font", reflect.TypeOf(op.Params[0]))
				log.Println("font", reflect.TypeOf(op.Params[1]))
				if name, ok := op.Params[0].(*pdfcore.PdfObjectName); ok {
					font2 := fontDict.Get(*name)
					//log.Printf("%T\n", font2)
					if indo, ok := font2.(*pdfcore.PdfIndirectObject); ok {
						//log.Println("AAAAA1", indo.DefaultWriteString())
						//log.Printf("AAAAA2 %T\n", indo.PdfObject)
						//log.Printf("AAAAA3 %s\n", indo.PdfObject.String())
						if fdict, ok := indo.PdfObject.(*pdfcore.PdfObjectDictionary); ok {
							//log.Println("AAAAA4 ~~~~~~~~~", fdict)
							//log.Println("AAAAA4 BASEFONT", fdict.Get("BaseFont"))
							newFont := fdict.Get("BaseFont").String()
							if newFont != lastFont {
								fmt.Print("<FONT name=", newFont, ">")
								lastFont = newFont
							}
						}
					}
				}
			}
		case Tstar_nextLine:
			fmt.Println("")

		case TD_setLeadingMoveText:
			if len(op.Params) != 2 {
				break
			}
			x, err := makeFloat(op.Params[0])
			if err != nil {
				log.Println(err)
				break
			}
			if math.Abs(x) > 150.0 {
				fmt.Print(" ")
			}
			_, err = makeFloat(op.Params[0])
			if err != nil {
				log.Println(err)
				break
			}
			fmt.Println("")

		case Backslash_nextLineShowText: // fix
			fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%")

		case TJ_showSpacedText:

			for _, p := range op.Params {
				//fmt.Println("\n", i, "TJ++++", p)

				if textOuts, ok := p.(*pdfcore.PdfObjectArray); ok {
					var line string
					for _, v := range *textOuts {
						//fmt.Println("values", i, j, v, reflect.TypeOf(v))
						if str, ok := v.(*pdfcore.PdfObjectString); ok {
							//fmt.Println("X", str)
							line = line + str.String()
						} else {
							x, err := makeFloat(v)
							if err != nil {
								log.Println(err)
								break
							}
							if math.Abs(x) > 1000 {
								line = line + "\t"
							} else if math.Abs(x) > 150 {
								line = line + " "
							}
						}
					}
					if len(line) > 0 {
						fmt.Print(line)
					}
				}

				if str, ok := p.(*pdfcore.PdfObjectString); ok {
					fmt.Println("X", str)
				} else {
					if val, ok := p.(*pdfcore.PdfObjectFloat); ok {
						if float64(*val) < -30.0 || float64(*val) > 30.0 {
							fmt.Println("XS", " ")
						} else {
							fmt.Println("XS", "ZZ")
						}
					}
				}
			}
		}
		//lastOp = op
	}
	return nil
}

const Tm_setTextMatrix = "Tm"
const Tj_showText = "Tj"
const Tf_setFont = "Tf"
const Tstar_nextLine = "T*"
const TJ_showSpacedText = "TJ"
const TD_setLeadingMoveText = "TD"

const Backslash_nextLineShowText = "'"

func makeFloat(p pdfcore.PdfObject) (float64, error) {
	if val, ok := p.(*pdfcore.PdfObjectFloat); ok {
		return float64(*val), nil
	}

	if val, ok := p.(*pdfcore.PdfObjectInteger); ok {
		return float64(*val), nil
	}
	return 0.0, errors.New("Neither float nor int:" + reflect.TypeOf(p).String())
}
