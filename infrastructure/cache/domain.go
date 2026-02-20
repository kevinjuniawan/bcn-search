package cache

import "kevinjuniawan/bookcabin/internal"

var MapSortTypeToKey = map[internal.SortType]string{
	internal.SortLowestPriceType:      "price",
	internal.SortHighestPriceType:     "price",
	internal.SortShortestDurationType: "duration",
	internal.SortLongestDurationType:  "duration",
	internal.SortDepartureType:        "departure",
	internal.SortArrivalType:          "arrival",
	internal.SortBestValueType:        "best",
}
