package checksum

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
)

func CalculateLocalChecksum(asks map[string]float64, bids map[string]float64) uint32 {
	asksSlc := mapToSortedSlice(asks, false)
	bidsSlc := mapToSortedSlice(bids, true)
	localChecksumStr := checksumString(asksSlc, bidsSlc)

	return crc32.ChecksumIEEE([]byte(localChecksumStr))
}

func checksumString(asks, bids []string) string {
	maxLevels := 25
	var checksumItems []string
	for i := 0; i < maxLevels; i++ {
		if i < len(bids) {
			checksumItems = append(checksumItems, bids[i])
		}
		if i < len(asks) {
			checksumItems = append(checksumItems, asks[i])
		}
	}

	return strings.Join(checksumItems, ":")
}

func mapToSortedSlice(m map[string]float64, desc bool) []string {
	sortedSlc := make([][2]float64, 0, len(m))
	for priceStr, quantity := range m {
		if quantity == 0 {
			continue
		}
		price, _ := strconv.ParseFloat(priceStr, 64)
		sortedSlc = append(sortedSlc, [2]float64{price, quantity})
	}

	sort.Slice(sortedSlc, func(i, j int) bool {
		if desc {
			return sortedSlc[i][0] > sortedSlc[j][0]
		}

		return sortedSlc[i][0] < sortedSlc[j][0]
	})

	maxLevels := 25
	res := make([]string, 0, min(len(sortedSlc), maxLevels))
	for i, val := range sortedSlc {
		if i >= maxLevels {
			break
		}
		priceStr := strconv.FormatFloat(val[0], 'f', -1, 64)
		quantityStr := strconv.FormatFloat(val[1], 'f', -1, 64)
		res = append(res, fmt.Sprintf("%s:%s", priceStr, quantityStr))
	}

	return res
}
