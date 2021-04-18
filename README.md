# MediaManager
## Purpose

Little utility to sort all Phot and video from a origin folder to a destination folder

Destination folder will be organised automatically with years and inside all years folders with months

Origin folder will be scanned and also the sub-folders

## Usage

MediaManager.exe TypeMedia ListOfAction originFolder DestFolder

TypeMedia : 

ALL : for Videos and Photo
PHOTO : for Photo only
VIDEO : for Video only

ListOfAction (separated by .) :

COPY : to copy media to destination File
MOVE : to move media to destination File
RENAME : to rename media before any other actions

Exemples :


MediaManager.exe ALL COPY.RENAME originFolder DestFolder
-> rename media and copy to destination folder

MediaManager.exe ALL MOVE.RENAME originFolder DestFolder
-> rename media and move to destination folder

MediaManager.exe PHOTO COPY.RENAME originFolder DestFolder
-> rename photo and copy to destination folder