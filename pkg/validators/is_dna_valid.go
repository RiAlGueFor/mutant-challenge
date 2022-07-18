package validators

import (
  "regex"
  "string"
)

func IsDNAValid(dna []string) bool{
  match := regexp.MustCompile(`[^ATCG]`)

  var initialLength int = len(dna[0])

	for k := 0 ; k < len(dna); k++ {
    if(len(dna[k])!=initialLength) {
      return false
    }
    if(match.MatchString(dna[k])) {
      return false
    }
	}
  return true
}
