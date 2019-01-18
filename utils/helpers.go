package utils

type Helper struct {
}

var UHelper Helper

//RemoveIntDuplicates de-duplicate int list
func (h Helper) RemoveIntDuplicates(elements []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

//RemoveStrDuplicates de-duplicate string list
func (h Helper) RemoveStrDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

//FormatSliceToIntMap convert slice of int to map
func (h Helper) FormatSliceToIntMap(all []int) map[int]int {
	bmap := make(map[int]int)
	for _, bv := range all {
		bmap[bv] = bv
	}
	return bmap
}
