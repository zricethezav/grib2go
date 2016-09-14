package grib2go

import (
	"encoding/binary"
	"fmt"
	"os"
)

// type defintions
const GENID_TO_YEAR_S1 = 6

type Time struct {
	Year   uint16 // year
	Month  uint8  // month + 1
	Day    uint8  // day
	Hour   uint8  // hour
	Minute uint8  // minute
	Second uint8  // second
}

type SectionHead struct {
	Len uint32
	Num uint8
}

type Section1 struct {
	OriginatingCenter         uint16
	OriginatingSubCenter      uint16
	MasterTablesVersion       uint8
	LocalTablesVersion        uint8
	ReferenceTimeSignificance uint8
	TimeAtDate                Time
	ProductionStatus          uint8
	Type                      uint8
}

type Section2 struct {
	LocalUse []uint8
}

type Section3 struct {
	Source                   uint8
	DataPointCount           uint32
	PointCountOctets         uint8
	PointCountInterpretation uint8
	TemplateNumber           uint16
	GridTemplate             interface{}
}

type Section4 struct {
	CoordinatesCount uint16
	TemplateNumber   uint16
	ProductTemplate  interface{}
}

type Section5 struct {
	PointsNumber   uint32
	TemplateNumber uint16
	DataTemplate   interface{}
}

type Section6 struct {
	SectionHead
	BitmapIndicator uint8
	Bitmap          interface{}
}

type Section7 struct {
	SectionHead
	Data interface{}
}

// readBytes reads n bytes from file f
func readBytes(f *os.File, elements ...interface{}) {
	for _, element := range elements {
		err := binary.Read(f, binary.BigEndian, element)
		CheckError(err)
	}
}

//SectionOne
func SectionOne(f *os.File, sectionHead *SectionHead, currOffset uint64) (sectionOne Section1) {
	fmt.Println("Enter Section One")
	binary.Read(f, binary.BigEndian, &sectionOne)
	return sectionOne
}

// TODO SectionTwo
func SectionTwo(f *os.File, sectionHead *SectionHead, currOffset uint64) (sectionTwo Section2) {
	fmt.Println("Enter Section Two")
	return sectionTwo
}

// SectionThree
// Section 3 will vary depending upon the map projection the grid is defined.
// A table listing more detailed information for each NDFD geographic region
// (CONUS, CONUS subsectors, Puerto Rico, Hawaii, Guam, Alaska) can be found here.
// Latitude and longitude values are encoded in 10 -6 degrees.  All longitudes are east.
func SectionThree(f *os.File, sectionHead *SectionHead, currOffset uint64) (sectionThree Section3) {
	fmt.Println("Enter Section Three")
	readBytes(f, &sectionThree.Source, &sectionThree.DataPointCount,
		&sectionThree.PointCountOctets, &sectionThree.PointCountInterpretation,
		&sectionThree.TemplateNumber)
	sectionThree.GridTemplate = TemplateHandler(f, sectionThree.TemplateNumber, sectionHead.Num)
	return sectionThree
}

// SectionFour
func SectionFour(f *os.File, sectionHead *SectionHead, currOffset uint64) (sectionFour Section4) {
	fmt.Println("Enter Section Four")
	readBytes(f, &sectionFour.CoordinatesCount, &sectionFour.TemplateNumber)
	sectionFour.ProductTemplate = TemplateHandler(f, sectionFour.TemplateNumber, sectionHead.Num)
	return sectionFour
}

// SectionFive
func SectionFive(f *os.File, sectionHead *SectionHead, currOffset uint64) (sectionFive Section5) {
	fmt.Println("Enter Section Five")
	readBytes(f, &sectionFive.PointsNumber, &sectionFive.TemplateNumber)
	sectionFive.DataTemplate = TemplateHandler(f, sectionFive.TemplateNumber, sectionHead.Num)
	return sectionFive
}

// SectionSix
func SectionSix(f *os.File, sectionHead *SectionHead, currOffset uint64) (sectionSix Section6) {
	fmt.Println("Enter Section Six")
	readBytes(f, &sectionSix.BitmapIndicator)
	if sectionSix.BitmapIndicator == 0 {
		// bytes 7 - nn describe the bitmap
		fmt.Println("bitmapIndicator = 0")
		sectionSix.Bitmap = readNBytes(f, uint32(sectionHead.Len-6))
	}
	return sectionSix
}

// SectionSeven
func SectionSeven(f *os.File, sectionHead *SectionHead, currOffset uint64) (sectionSeven Section7) {
	fmt.Println("Enter Section Seven")
	sectionSeven.Data = readNBytes(f, uint32(sectionHead.Len-5))
	return sectionSeven
}

func SectionEight(f *os.File, sectionHead *SectionHead, currOffset uint64) {
	fmt.Println("Enter Section Seven")
}
