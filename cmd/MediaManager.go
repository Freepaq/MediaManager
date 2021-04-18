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
	total := len(rows)
	proc := 1
	fmt.Print(MediaUtils.SEPARATOR)
	nbPhoto := 0
	nbVideo := 0
	nbVidePro := 0
	nbPhotoPro := 0

	for _, file := range rows {
		fmt.Println(" " + strconv.Itoa(proc) + "/" + strconv.Itoa(total))
		proc++
		meta, err := MediaUtils.GetMeta(file)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Source file : [" + file + "]")
			fmt.Println("Creation : [" + meta.CreationDate.String() + "] taken from " + meta.MetaOrigin)

		}
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
		fmt.Print(MediaUtils.SEPARATOR)
		if meta.TypeOfMedia == MediaUtils.VIDEO {
			nbVideo++
			if meta.Proccessed {
				nbVidePro++
			}
		}
		if meta.TypeOfMedia == MediaUtils.PHOTO {
			nbPhoto++
			if meta.Proccessed {
				nbPhotoPro++
			}
		}

	}
	fmt.Println("")
	fmt.Println("Number of Video : " + strconv.Itoa(nbVideo))
	fmt.Println("Number of Video copied : " + strconv.Itoa(nbVidePro))
	fmt.Println(MediaUtils.SEPARATOR)
	fmt.Println("Number of Photo : " + strconv.Itoa(nbPhoto))
	fmt.Println("Number of Photo copied : " + strconv.Itoa(nbPhotoPro))
	fmt.Println(MediaUtils.SEPARATOR)
	fmt.Println("End Time :" + MediaUtils.CurrentTime)

}
