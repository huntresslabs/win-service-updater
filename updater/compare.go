package updater

import (
	"regexp"
	"strings"
	"unicode"
)

// CompareVersions compares two versions and returns an integer that indicates
// their relationship in the sort order.
// Return a negative number if versionA is less than versionB, 0 if they're
// equal, a positive number if versionA is greater than versionB.
func CompareVersions(a string, b string) int {
	if a == "" {
		if b == "" {
			return A_EQUAL_TO_B
		}

		return A_LESS_THAN_B
	}

	if b == "" {
		return A_GREATER_THAN_B
	}

	// Convert version to lowercase, and replace all instances of "release candidate" with "rc"
	re := regexp.MustCompile(`release[\s]+candidate`)
	a = re.ReplaceAllString(strings.ToLower(a), "rc")
	b = re.ReplaceAllString(strings.ToLower(b), "rc")

	// compare indices
	iVerA := 0
	iVerB := 0

	lastAWasLetter := true
	lastBWasLetter := true

	for {
		greekIndA := iVerA
		greekIndB := iVerB

		objA := GetNextObject(a, &iVerA, &lastAWasLetter)
		objB := GetNextObject(b, &iVerB, &lastBWasLetter)

		if (!lastBWasLetter && objB != "") && (objA == "" || lastAWasLetter) {
			objA = "0"
			iVerA = greekIndA
		} else if (!lastAWasLetter && objA != "") && (objB == "" || lastBWasLetter) {
			objB = "0"
			iVerB = greekIndB
		}

		// find greek *index for A and B
		greekIndA = -1
		greekIndB = -1
		if lastAWasLetter {
			greekIndA = GetGreekIndex(objA)
		}
		if lastAWasLetter {
			greekIndA = GetGreekIndex(objA)
		}

		if objA == "" && objB == "" {
			return 0 //versions are equal
		}

		// objB != null
		if objA == "" {
			//if versionB has a greek word, then A is greater
			if greekIndB != -1 {
				return 1
			}

			return -1
		}

		// objA != null
		if objB == "" {
			//if versionA has a greek word, then B is greater
			if greekIndA != -1 {
				return -1
			}

			return 1
		}

		if unicode.IsDigit([]rune(objA)[0]) == unicode.IsDigit([]rune(objB)[0]) {
			var strComp int
			if unicode.IsDigit([]rune(objA)[0]) {
				//compare integers
				strComp = IntCompare(objA, objB)

				if strComp != 0 {
					return strComp
				}
			} else {
				if greekIndA == -1 && greekIndB == -1 {
					//compare non-greek strings
					strComp = strings.Compare(objA, objB)

					if strComp != 0 {
						return strComp
					}
				} else if greekIndA == -1 {
					return 1 //versionB has a greek word, thus A is newer
				} else if greekIndB == -1 {
					return -1 //versionA has a greek word, thus B is newer
				} else {
					//compare greek words
					if greekIndA > greekIndB {
						return 1
					}

					if greekIndB > greekIndA {
						return -1
					}
				}
			}
		} else if unicode.IsDigit([]rune(objA)[0]) {
			return 1 //versionA is newer than versionB
		} else {
			return -1 //verisonB is newer than versionA
		}
	}

	aNum := convertVerToNum(a)
	bNum := convertVerToNum(b)

	if aNum < bNum {
		return A_LESS_THAN_B
	}
	if aNum > bNum {
		return A_GREATER_THAN_B
	}

	return A_EQUAL_TO_B
}

func GetNextObject(version string, index *int, lastWasLetter *bool) string {
	//1 == string, 2 == int, -1 == neither
	StringOrInt := -1

	startIndex := *index

	for len(version) != *index {
		if StringOrInt == -1 {
			if unicode.IsLetter([]rune(version)[*index]) {
				startIndex = *index
				StringOrInt = 1
			} else if unicode.IsDigit([]rune(version)[*index]) {
				startIndex = *index
				StringOrInt = 2
			} else if *lastWasLetter && !unicode.IsSpace([]rune(version)[*index]) {
				*index++
				*lastWasLetter = false
				return "0"
			}
		} else if StringOrInt == 1 && !unicode.IsLetter([]rune(version)[*index]) {
			break
		} else if StringOrInt == 2 && !unicode.IsDigit([]rune(version)[*index]) {
			break
		}

		*index++
	}

	// set the last "type" retrieved
	*lastWasLetter = (StringOrInt == 1)

	// return the retitrved sub-string
	if StringOrInt == 1 || StringOrInt == 2 {
		return version[startIndex:*index]
	}

	// was neither a string nor and int
	return ""
}

func IntCompare(a string, b string) int {
	lastZero := -1

	// Clear any preceding zeros

	for i := 0; i < len(a); i++ {
		if a[i] != '0' {
			break
		}

		lastZero = i
	}

	if lastZero != -1 {
		if len(a) > lastZero+1 {
			a = a[lastZero+1 : len(a)-(lastZero+1)]
		}
	}

	lastZero = -1

	for i := 0; i < len(b); i++ {
		if b[i] != '0' {
			break
		}

		lastZero = i
	}

	if lastZero != -1 {
		if len(b) > lastZero+1 {
			b = b[lastZero+1 : len(b)-(lastZero+1)]
		}
	}

	if len(a) > len(b) {
		return 1
	}

	if len(a) < len(b) {
		return -1
	}

	return strings.Compare(a, b)
}

func GetGreekIndex(obj string) int {
	greekLetters := map[string]int{
		"alpha": 0, "beta": 1, "gamma": 2,
		"delta": 3, "epsilon": 4, "zeta": 5,
		"eta": 6, "theta": 7, "iota": 8,
		"kappa": 9, "lambda": 10, "mu": 11,
		"nu": 12, "xi": 13, "omicron": 14,
		"pi": 15, "rho": 16, "sigma": 17,
		"tau": 18, "upsilon": 19, "phi": 20,
		"chi": 21, "psi": 22, "omega": 23,
	}

	return greekLetters[obj]
}
