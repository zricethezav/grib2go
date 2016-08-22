package grib2go

import (
    "testing"
    "os"
    "log"
)

func TestProcessGrib(t *testing.T) {
    log.Println("TestProcessGrib")
    pwd, _ := os.Getwd()
    f, err := os.Open(pwd+"/data/nldas.t12z.force-a.grb2f00")
    CheckError(err)
    ProcessGrib(f)
}
