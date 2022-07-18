package mutantDNA

import(
  "encoding/json"
  "errors"
)

func ScanningDNA(dna []string, letter2Check string) bool {
  oCounter:=0
  vCounter:=0
  aux_i:=0
  aux_j:=0

  repString := strings.Repeat(letter2Check, 4)

  for i := 0; i < len(dna) ; i++{
    // Horizontal ScanningDNA
    if(strings.Contains(dna[i], repString)){
      return true
    }

    oCounter = 0
    vCounter = 0
    for j := 0; j < len(dna[i]) ; j++{
      if (string(dna[i][j]) == letter2Check){
        oCounter++
        vCounter++

        // Vertical ScanningDNA
        if(len(dna)-i>=4){
          aux_i=i
          for vCounter < 4 {
            aux_i++
            if (string(dna[aux_i][j]) == letter2Check){
              vCounter++
            } else {
              break;
            }
          }
          if(vCounter == 4){
            return true
          } else {
            vCounter = 0
          }
        }
        // Oblique ScanningDNA
        if(len(dna[i])-j>=4){
          // Go on

          // Oblique ScanningDNA
          if(len(dna)-i>=4){
            // Downside
            aux_i = i
            aux_j = j
            for oCounter < 4 {
          		aux_i++
              aux_j++
              if (string(dna[aux_i][aux_j]) == letter2Check){
                oCounter++
              } else {
                oCounter=0
                break;
              }
          	}

            if(oCounter == 4){
              return true
            } else {
              oCounter=0
            }
          }

          if(i>=3){
            // Upside
            aux_i=i
            aux_j=j
            for oCounter < 4 {
          		aux_i--
              aux_j++
              if (string(dna[aux_i][aux_j]) == letter2Check){
                oCounter++
              } else {
                oCounter=0
                break;
              }
          	}

            if(oCounter == 4){
              return true
            } else {
              oCounter=0
            }
          }
        } else {
          oCounter=0
          break;
        }
      }
    }
  }
  return false
}
