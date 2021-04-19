package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Freepaq/MediaManagement/pkg/MediaUtils"
)

var actions []string

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Please check arguments")
		os.Exit(-1)
	}
	fmt.Println("Excution Time :" + MediaUtils.CurrentTime)
	mediaType := os.Args[1]
	action := os.Args[2]
	origin := os.Args[3]
	dest := os.Args[4]
	actions = strings.Split(action, ".")
	rows := MediaUtils.GetListOfFile(origin, mediaType)
	fmt.Println("Actions requested : " + strings.Join(actions, " - "))
	fmt.Println("Media Type : " + mediaType)
	fmt.Println("Eligible Photo Type : " + strings.Join(MediaUtils.EligiblePhotoFile, " "))
	fmt.Println("Eligible Video Type : " + strings.Join(MediaUtils.EligibleVideoFile, " "))
	total := len(rows)
	proc := 1
	fmt.Print(MediaUtils.SEPARATOR)
	nbPhoto := 0
	nbVideo := 0
	nbVidePro := 0
	nbPhotoPro := 0
	metaTypeVideo := make(map[string]int)
	metaTypePhoto := make(map[string]int)

	for _, file := range rows {
		fmt.Println(" " + strconv.Itoa(proc) + "/" + strconv.Itoa(total))
		proc++
		meta, err := MediaUtils.GetMeta(file)
		fmt.Println("Source file : [" + file + "]")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Creation : [" + meta.CreationDate.String() + "] taken from " + meta.MetaOrigin + " [" + meta.FromMeta + "]")

			if MediaUtils.Contains(actions, "RENAME") {
				MediaUtils.Rename(&meta)
			}
			if MediaUtils.Contains(actions, "COPY") {
				MediaUtils.Copy(&meta, dest, false)
			}
			if MediaUtils.Contains(actions, "REPLACE") {
				MediaUtils.Copy(&meta, dest, true)
			}
			if MediaUtils.Contains(actions, "MOVE") {
				if result := MediaUtils.Copy(&meta, dest, true); result {
					MediaUtils.Delete(meta)
				}
			}

			if meta.TypeOfMedia == MediaUtils.VIDEO {
				nbVideo++
				if meta.Proccessed {
					nbVidePro++
				}
				if count, ok := metaTypeVideo[meta.MetaOrigin]; ok {
					count++
					metaTypeVideo[meta.MetaOrigin] = count
				} else {
					metaTypeVideo[meta.MetaOrigin] = 1
				}
			}
			if meta.TypeOfMedia == MediaUtils.PHOTO {
				nbPhoto++
				if meta.Proccessed {
					nbPhotoPro++
				}
				if count, ok := metaTypePhoto[meta.MetaOrigin]; ok {
					count++
					metaTypePhoto[meta.MetaOrigin] = count
				} else {
					metaTypePhoto[meta.MetaOrigin] = 1
				}
			}
		}
		fmt.Print(MediaUtils.SEPARATOR)
	}
	fmt.Println("")
	fmt.Println("Number of Video : " + strconv.Itoa(nbVideo))
	fmt.Println("Number of Video copied : " + strconv.Itoa(nbVidePro))
	for key, m := range metaTypeVideo {
		fmt.Println("Number of Date taken from " + key + " : " + strconv.Itoa(m))
	}

	fmt.Println(MediaUtils.SEPARATOR)
	fmt.Println("Number of Photo : " + strconv.Itoa(nbPhoto))
	fmt.Println("Number of Photo copied : " + strconv.Itoa(nbPhotoPro))
	for key, m := range metaTypePhoto {
		fmt.Println("Number of Date taken from " + key + " : " + strconv.Itoa(m))
	}
	fmt.Println(MediaUtils.SEPARATOR)
	fmt.Println("End Time :" + MediaUtils.CurrentTime)

}
