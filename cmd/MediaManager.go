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
	fmt.Println(MediaUtils.SEPARATOR)
	nbPhoto := 0
	nbVideo := 0
	for _, file := range rows {

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
		fmt.Println(MediaUtils.SEPARATOR)
		if meta.TypeOfMedia == MediaUtils.VIDEO {
			nbVideo++
		}
		if meta.TypeOfMedia == MediaUtils.PHOTO {
			nbPhoto++
		}

	}
	fmt.Println("Number of Video : " + strconv.Itoa(nbVideo))
	fmt.Println("Number of Photo : " + strconv.Itoa(nbPhoto))
	fmt.Println("End Time :" + MediaUtils.CurrentTime)
}
