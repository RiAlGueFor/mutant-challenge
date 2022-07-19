package validators

import (
  "testing"
  "strings"
  "strconv"
)

type TestDNA struct {
    in  string
    out bool
}

var testDNA = []TestDNA{
  {"ATGCGA,CAGTAC,TCAGGT,AGCAGG,TTACGG,TCACTC",true},
  {"ACTGAC,TACGAC,GTCATG,AGCAGG,TTACGG,TCACTC",true},
  {"ATGCGA,CAGTAC,TCAGGT,AGCAGG,TTACGG,TCACTC",true},
  {"ATGCGA,CAGTAC,TCAGGT,AGCAGG,TTAGG,TCACTC",false},
  {"ATGCGA,CAGTAC,TCAGgT,AGCAGG,TTACGG,TCACTC",false},
  {"ATGCGA,CAGTAC,TCAGGT,AGCGG,TTACGG,TCACTC",false},
}

func TestIsValidDNA(t *testing.T) {
    var strSlice []string
    for _, test := range testDNA {
        strSlice = strings.Split(test.in, ",")
        isValid := IsDNAValid(strSlice)
        if isValid != test.out {
            t.Errorf("got %q, wanted %q",strconv.FormatBool(isValid), strconv.FormatBool(test.out))
        }
    }
}
