package updater

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	A_LESS_THAN_B    = -1
	A_EQUAL_TO_B     = 0
	A_GREATER_THAN_B = 1
)

func convertVerToNum(ver string) int {
	var num int
	fields := strings.Split(ver, ".")
	for i, field := range fields {
		x, _ := strconv.Atoi(field)
		if i == 0 {
			num = num + (x << 24)
		}
		if i == 1 {
			num = num + (x << 16)
		}
		if i == 2 {
			num = num + (x << 8)
		}
		if i == 3 {
			num = num + x
		}
	}
	return num
}

// CompareVersions compares two versions and returns an integer that indicates
// their relationship in the sort order.
// Return a negative number if versionA is less than versionB, 0 if they're
// equal, a positive number if versionA is greater than versionB.
func CompareVersions(a string, b string) int {
	aNum := convertVerToNum(a)
	bNum := convertVerToNum(b)

	// fmt.Printf("a = %s (%d)\n", a, aNum)
	// fmt.Printf("b = %s (%d)\n", b, bNum)

	if aNum < bNum {
		return A_LESS_THAN_B
	}
	if aNum > bNum {
		return A_GREATER_THAN_B
	}
	//if aNum == bNum {
	// return 0
	//}
	return A_EQUAL_TO_B
}

func GetUpdateDetails(extractedFiles []string) (udt ConfigUDT, updates []string, err error) {
	udtFound := false

	for _, f := range extractedFiles {
		if path.Base(f) == "updtdetails.udt" {
			udt, err = ParseUDT(f)
			if err != nil {
				return ConfigUDT{}, updates, err
			}
			udtFound = true
		} else {
			updates = append(updates, f)
		}
	}

	if !udtFound {
		err := fmt.Errorf("no udt file found")
		return ConfigUDT{}, updates, err
	}

	return udt, updates, nil
}

func InstallUpdate(udt ConfigUDT, srcFiles []string, installDir string) error {
	// stop services
	// for _, s := range udt.ServiceToStopBeforeUpdate {
	// 	fmt.Printf("Stopping %s\n", ValueToString(&s))
	// }

	for _, f := range srcFiles {
		// fmt.Printf("Moving %s\n", f)
		err := MoveFile(f, installDir)
		if err != nil {
			return err
		}
	}

	// start services
	// for _, s := range udt.ServiceToStartAfterUpdate {
	// 	fmt.Printf("Starting %s\n", ValueToString(&s))
	// }

	return nil
}

func MoveFile(file string, dstDir string) error {
	dst := filepath.Join(dstDir, filepath.Base(file))
	// fmt.Println(dst)
	// Rename() returns *LinkError
	return os.Rename(file, dst)
}
