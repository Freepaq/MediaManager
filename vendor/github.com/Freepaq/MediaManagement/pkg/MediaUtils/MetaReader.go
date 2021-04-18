package MediaUtils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	log "github.com/dsoprea/go-logging"
)

type MyMapping map[string]interface{}

func ReadVideoMeta(fname string, fileStr *FileStruct) error {
	//fmt.Println("Current file :" + fname)
	cmd := exec.Command("mediainfo/MediaInfo.exe", "--Output=JSON", "--Logfile=text.txt", fname)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error reading METADATA :" + fmt.Sprint(err) + ": " + string(output))
		return err
	}
	var encodeDate string
	resultingMap := MyMapping{}
	//	file, _ := ioutil.ReadFile("temp.txt")
	if err := json.Unmarshal(output, &resultingMap); err != nil {
		fmt.Println("json.Compact:", err)
		if serr, ok := err.(*json.SyntaxError); ok {
			fmt.Println("Occurred at offset:", serr.Offset)
			return serr
		}
	} else {
		encodeDate = search(resultingMap, []string{"DateTime", "Encoded_Date", "File_Created_Date"})
	}
	if encodeDate == "" {
		readFromFile(fname, fileStr)
	} else {
		time, err := time.Parse("2006-01-02T15:04:05", encodeDate)
		if nil != err {
			fmt.Println(err)
			return err
		} else {
			fileStr.CreationDate = time
			fileStr.MetaOrigin = METAORIGINMETA
		}
	}
	return nil
}

func search(str MyMapping, value []string) string {
	var result = ""
	var media MyMapping
	media = str["media"].(map[string]interface{})
	//	fmt.Println(reflect.TypeOf(media))
	if nil != media {
		var track []interface{}
		track = media["track"].([]interface{})
		if nil != track {
			for _, m := range track {
				for key, v := range m.(map[string]interface{}) {
					if Contains(value, key) {

						var date = v.(string)
						if date != "" {
							dateArray := strings.Split(date, " ")
							return dateArray[1] + "T" + dateArray[2]
						}
					}
				}
			}
		}
	}
	return result
}

func ReadPhotoMeta(fname string, fileStr *FileStruct) {
	x, err := extractTag(fname, "DateTime")
	if err != nil {
		readFromFile(fname, fileStr)
	} else {
		fileStr.CreationDate, _ = time.Parse("2006:01:02 15:04:05", x)
		fileStr.MetaOrigin = METAORIGINMETA
	}
}

func extractTag(fname string, tagName string) (string, error) {
	rawExif, err := exif.SearchFileAndExtractExif(fname)
	if err != nil {
		return "", err
	}
	im, err := exifcommon.NewIfdMappingWithStandard()
	log.PanicIf(err)
	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return "", err
	}

	rootIfd := index.RootIfd

	// We know the tag we want is on IFD0 (the first/root IFD).
	results, err := rootIfd.FindTagWithName(tagName)
	if err != nil {
		return "", err
	}

	if len(results) != 1 {
		err := log.Wrap("Error reading tags")
		return "", err
	}

	ite := results[0]

	valueRaw, err := ite.Value()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	value := valueRaw.(string)

	return value, nil
}

func readFromFile(fname string, fileStr *FileStruct) {
	fileStat, err := os.Stat(fname)
	if err != nil {
		log.PanicIf(err)
	}
	fileStr.CreationDate = fileStat.ModTime()
	fileStr.MetaOrigin = METAORIGINFILE
}
