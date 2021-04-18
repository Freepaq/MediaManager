package MediaUtils

import (
	"strings"
	"time"
)

var TimestampFormat = "2006-01-02_150405"
var timestampFormatLog = "2006_01_02"
var CurrentTime = time.Now().Format(TimestampFormat)

var EligiblePhotoFile = []string{".jpg", ".gif", ".jpeg", ".png"}
var EligibleVideoFile = []string{".mp4", ".mov", ".3gp"}
var Months = map[int]string{1: "01 - Janvier", 2: "02 - Février", 3: "03 - Mars",
	4: "04 - Avril", 5: "05 - Mai", 6: "06 - Juin", 7: "07 - Juillet",
	8: "08 - Août", 9: "09 - Septembre", 10: "10 - Octobre", 11: "11 - Novembre", 12: "12 - Décembre"}

const PHOTO = "PHOTO"
const VIDEO = "VIDEO"
const ALL = "ALL"
const METAORIGINFILE = "FILE"
const METAORIGINMETA = "META"
const SEPARATOR = "---------------------------------------------------------------------------------------------------"

type FileStruct struct {
	CreationDate   time.Time
	FullName       string
	MetaOrigin     string
	Name           string
	Extension      string
	OriginDir      string
	NewName        string
	NewFullName    string
	DestinationDir string
	TypeOfMedia    string
	Proccessed     bool
}

func IsPhotoEligible(ext string) bool {
	for _, e := range EligiblePhotoFile {
		if strings.ToLower(ext) == e {
			return true
		}
	}
	return false
}

func IsVideoEligible(ext string) bool {
	for _, e := range EligibleVideoFile {
		if strings.ToLower(ext) == e {
			return true
		}
	}
	return false
}

func IsMediEligible(ext string) bool {
	return IsPhotoEligible(ext) || IsVideoEligible(ext)
}

func Contains(actions []string, key string) bool {
	for _, e := range actions {
		if strings.ToLower(key) == strings.ToLower(e) {
			return true
		}
	}
	return false
}
