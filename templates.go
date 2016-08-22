package grib2go

import (
	"encoding/binary"
	"fmt"
	"os"
)

type GridTemplateDescriptor struct {
	EarthShape   uint8  // Shape of the Earth (See Code Table 3.2)
	RadiusFactor uint8  // Scale Factor of radius of spherical Earth
	RadiusValue  uint32 // Scale value of radius of spherical Earth
	MajorFactor  uint8  // Scale factor of major axis of oblate spheroid Earth
	MajorValue   uint32 // Scaled value of major axis of oblate spheroid Earth
	MinorFactor  uint8  // Scale factor of minor axis of oblate spheroid Earth
	MinorValue   uint32 // Scaled value of minor axis of oblate spheroid Earth
}

type GridTemplate0 struct {
	GridTemplateDesc GridTemplateDescriptor
	Ni               uint32 // Ni—number of points along a parallel
	Nj               uint32 // Nj—number of points along a meridian
	BasicAngle       uint32 // Basic angle of the initial production domain (see Note 1)
	BasicAngleDiv    uint32 // Subdivisions of basic angle used to define extreme
	// longitudes and latitudes, and direction increments
	// (see Note 1)
	La_1                    uint32 // latitude of first grid point (see Note 1)
	Lo_1                    uint32 // longitude of first grid point (see Note 1)
	ResolutionComponentFlag uint8  // Resolution and component flags
	La_2                    uint32 // latitude of last grid point
	Lo_2                    uint32 // longitude of last grid point
	Di                      uint32 // i direction increment (see Notes 1 and 5)
	Dj                      uint32 // j direction increment (see Note 1 and 5)
	ScanningMode            uint8  // Scanning mode (flags — see Flag Table 3.4 and Note 6)
}

type ScaleSurface struct {
	TypeOfFirst       uint8
	FirstFixedFactor  uint8
	FirstFixedValue   uint32
	TypeOfSecond      uint8
	SecondFixedFactor uint8
	SecondFixedValue  uint32
}

type ProductTemplate0 struct {
	CategoryParam         uint8
	ParamNum              uint8
	ProcessGenType        uint8
	BackgrndProcessID     uint8
	ForecastProcessID     uint8
	ObservationHours      uint16
	ObservationMinsCutoff uint8
	TimeRangeUnit         uint8
	ForecastTimeUnit      uint32
	ScaleDesc             ScaleSurface
}

type DataTemplateDescriptor struct {
	ReferenceValue     uint32
	BinaryScaleFactor  uint16
	DecimalScaleFactor uint16
	NumBitsPack        uint8
	TypeFieldValue     uint8
}

type DataTemplate0 struct {
	DataTemplateDesc DataTemplateDescriptor
}

type DataTemplate3 struct {
	DataTemplateDesc       DataTemplateDescriptor
	GroupSplitMethod       uint8
	MissingValMgmt         uint8
	PrimaryMissingValSub   uint32
	SecondaryMissingValSub uint32
	NumberGroups           uint32
	RefGroupWidths         uint8
	NumBitsGroupWidths     uint8
	RefGroupLength         uint32
	LengthIncr             uint8
	TrueLengthLastGroup    uint32
}

// TemplateHandler returns a template based on templateNumber and sectionNumber as an
// interface. Each section maps to a specific template type and each template type has
// multiple templates. For example: sectionNumber == 3 would return a grid template.
func TemplateHandler(f *os.File, templateNumber uint16, sectionNumber uint8) interface{} {
	fmt.Println("Enter Template Handler")
	fmt.Println("sectionNumber: ", sectionNumber)
	fmt.Println("templateNumber: ", templateNumber)
	switch sectionNumber {
	case SECTION_3:
		switch templateNumber {
		case 0:
			// Latitude/Longitude (See Template 3.0)    Also called Equidistant Cylindrical or Plate Caree
			var g0 GridTemplate0
			binary.Read(f, binary.BigEndian, &g0)
			return g0
		case 1:
			// TODO
		}
	case SECTION_4:
		switch templateNumber {
		case 0:
			var p0 ProductTemplate0
			binary.Read(f, binary.BigEndian, &p0)
			return p0
		case 1:
			//TODO
		}
	case SECTION_5:
		switch templateNumber {
		case 0:
			var d0 DataTemplate0
			binary.Read(f, binary.BigEndian, &d0)
			return d0
		case 1:
			//TODO
		}
	}
	return 0
}
