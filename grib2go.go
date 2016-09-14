package grib2go

// http://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc.shtml
// defining some constants
// based on http://www.wmo.int/pages/prog/www/WMOCodes/Guides/GRIB/GRIB2_062006.pdf

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Message struct {
	SectionOne   Section1
	SectionTwo   Section2
	SectionThree Section3
	SectionFour  Section4
	SectionFive  Section5
	SectionSix   Section6
	SectionSeven Section7
}

// constants
const (
	SECTION_8      = 8
	SECTION_7      = 7
	SECTION_6      = 6
	SECTION_5      = 5
	SECTION_4      = 4
	SECTION_3      = 3
	SECTION_2      = 2
	SECTION_1      = 1
	GRIB_INDICATOR = "GRIB"
)

func CheckSection(f *os.File) (sectionHead SectionHead, endGrib bool) {
	binary.Read(f, binary.BigEndian, &sectionHead)
	if sectionHead.Len == 926365495 {
		return sectionHead, true
	}
	return sectionHead, false
}

// CheckError is a general error checker that panics if error e is not nil
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

// readNBytes reads n bytes from file f
func readNBytes(f *os.File, n uint32) []byte {
	buf := make([]byte, n)
	_, err := f.Read(buf)
	CheckError(err)
	return buf
}

// verifyGrib checks the first 4 bytes of file f and verify it == "GRIB"
// and seek forward 16 bytes from the start to get to section 1. Section 1 is
// always 16 bytes.
func verifyGrib(f *os.File) bool {
	buf := make([]byte, 4)
	_, err := f.Read(buf)
	CheckError(err)
	if string(buf) != GRIB_INDICATOR {
		return false
	}

	f.Seek(16, 0)
	return true
}

// ProcessGrib
// NOTE: for the first file we process, refer to page 49 in
// https://www.wmo.int/pages/prog/www/WMOCodes/Guides/GRIB/GRIB2_062006.pdf
// This is a simple Grib2 message with one data set
func ProcessGrib(f *os.File) {
	if verifyGrib(f) == false {
		panic("Not Grib")
	}
	var sectionHead SectionHead
	var endGrib bool
	for {
		var message Message
		currOffset, _ := f.Seek(0, 1)
		sectionHead, endGrib = CheckSection(f)
		if endGrib {
			fmt.Println("DONE")
			break
		}
		fmt.Println("---------------------")
		fmt.Printf("Section: %d\nLength of Section: %d bytes\nCurrent Offset: byte %d\n",
			sectionHead.Num, sectionHead.Len, currOffset)
		switch sectionHead.Num {
		case SECTION_1:
			// Identification Section
			message.SectionOne = SectionOne(f, &sectionHead, uint64(currOffset))
			fmt.Printf("%+v\n\n", message.SectionOne)
		case SECTION_2:
			// Local Use Section
			message.SectionTwo = SectionTwo(f, &sectionHead, uint64(currOffset))
			fmt.Printf("%+v\n\n", message.SectionTwo)
		case SECTION_3:
			// Grid Definition Section
			message.SectionThree = SectionThree(f, &sectionHead, uint64(currOffset))
			fmt.Printf("%+v\n\n", message.SectionThree)
		case SECTION_4:
			// Product Definition Section
			message.SectionFour = SectionFour(f, &sectionHead, uint64(currOffset))
			fmt.Printf("%+v\n\n", message.SectionFour)
		case SECTION_5:
			// Data Representation Section
			message.SectionFive = SectionFive(f, &sectionHead, uint64(currOffset))
			fmt.Printf("%+v\n\n", message.SectionFive)
		case SECTION_6:
			// Bit Map Section
			message.SectionSix = SectionSix(f, &sectionHead, uint64(currOffset))
		case SECTION_7:
			// Data Section
			message.SectionSeven = SectionSeven(f, &sectionHead, uint64(currOffset))
		case SECTION_8:
			// End Section
			SectionEight(f, &sectionHead, uint64(currOffset))
		}
		if sectionHead.Num == SECTION_8 {
			currOffset, _ := f.Seek(0, 1)
			fmt.Println("CURR OFFSET: ", currOffset)
			break
		}
	}
}
