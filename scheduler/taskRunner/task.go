package taskRunner

import (
	"errors"
	"log"
	"os"
	"sync"
	"videoServer/scheduler/dbops"
	"videoServer/scheduler/ossops"
)

func deleteVideo(vid string) error{
	err:= os.Remove(VIDEO_PATH+vid)
	if err!=nil && !os.IsNotExist(err){
		log.Printf("Deleting video error: %v",err)
		return err
	}

	//delete file in oss
	ossfn:="videos/"+vid
	bn:="bucketName XXX"
	ok:=ossops.DeleteObject(ossfn,bn)
	if !ok{
		log.Printf("Delet video error oss op failed")
		return errors.New("delete video error")
	}

	return nil
}

func VideoClearDispatcher(dc dataChan) error{
	res,err:=dbops.ReadVideoDeletionRecord(3)
	if err!=nil{
		log.Printf("Video clear dispatcher error: %v",err)
		return err
	}

	if len(res)==0{
		return errors.New("All tasks finished")
	}

	for _,id:=range res{
		dc <- id
	}
	return nil
}

func VideoClearExecutor(dc dataChan) error{

	errMap:=&sync.Map{}
	var err error
	forloop:
		for{
			select {
				case vid:=<-dc:
					go func(id interface{}) { //may duplicate read
						if err:=deleteVideo(id.(string));err!=nil{ //del in os
							errMap.Store(id,err)
							return
						}
						if err:=dbops.DelVideoDeletionRecord(id.(string));err!=nil{//del in db
							errMap.Store(id,err)
							return
						}
					}(vid)
				default:
					break forloop
			}
		}
	errMap.Range(func(k,v interface{}) bool{
		err = v.(error)
		if err!=nil{
			return false
		}
		return false
	})
	return err
}

