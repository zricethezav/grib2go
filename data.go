package grib2go

// This file serves as an interpreter of data defined in section 5's data template
// Finish em out as we run into them

// http://www.nco.ncep.noaa.gov/pmb/docs/grib2/grib2_doc.shtml
// defining some constants
// based on http://www.wmo.int/pages/prog/www/WMOCodes/Guides/GRIB/GRIB2_062006.pdf

import (
	"os"
)

// NOTE: follow section 3.3.2.1 for examples on how to process complex packing/spatial differencing
// messages

func template_5_3(f *os.File, d3 *DataTemplate3, dataLen uint64) {
	// NOTE: at this point the file pointer f is already offset to byte 6 of section 7
	// aka, the beginning of the actual data content.
	// Remember this is the data section (section 7) and we fill data based on
	// what is specified in section 5
	numberGroups := d3.NumberGroups

	// spatial differencing
	// (octets 6-ww in data template 7.3)
	if d3.NumBytesReqExtra != 0 {
		spatialDiffSize := uint32(d3.NumBytesReqExtra * 8)
		// At order 1, an initial field of values f is replaced by a new field of values g,
		// where g1 = f1, g2 = f2, ..., gn = fn - fn-1.
		// At order 2, the field of values g is itself replaced by a new field of values h,
		// where h1 = f1, h2 = f2 , h3 = g3- g2, ..., hn = gn - gn - 1.
		ival := f.readNBytes(f, spatialDiffSize)
	}

	// before we grab all the reference values and group widths/lengths we need to figure out how
	// many bytes are necessary for each read
	if d3.DataTemplateDesc.NumBitsPack%8 == 0 {
		numBytesRef := uint32(d3.DataTemplateDesc.NumBitsPack / 8)
	} else {
		numBytesRef := uint32(d3.DataTemplateDesc.NumBitsPack/8) + 1
	}

	// declare reference, group length, group width vals (assume max, 8 bytes when declaring slice)
	referenceVals := make([]uint64, numberGroups)
	groupWidths := make([]uint64, numberGroups)
	groupLengths := make([]uint64, numberGroups)

	for i := 0; i <= numberGroups; i++ {
		referenceVals[i] = f.readNBytes(f, numBytesRef)
	}
	for i := 0; i <= numberGroups; i++ {
		groupWidths[i] = f.readNBytes(f, numBytesRef)
	}
	for i := 0; i <= numberGroups; i++ {
		groupLengths[i] = f.readNBytes(f, numBytesRef)
	}
}
