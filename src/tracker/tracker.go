package tracker

import (
	"log"
	"steam"
	"storage"
	"sync"
)

func worker(appIdChan chan int, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()

	for appId := range appIdChan {
		count, err := steam.GetUserCount(appId)
		if err != nil {
			log.Print(err)
		}
		storage.MakeUsageRecord(appId, count)
		log.Println(appId, "-", count)
	}
}

func MakeRecords() error {
	apps, err := storage.AllTrackableApps()
	if err != nil {
		return err
	}

	appIdChan := make(chan int)
	wg := new(sync.WaitGroup)

	// Adding routines to workgroup and running then
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go worker(appIdChan, wg)
	}

	// Processing all links by spreading them to `free` goroutines
	for _, app := range apps {
		appIdChan <- app.Id
	}

	// Closing channel (waiting in goroutines won't continue any more)
	close(appIdChan)

	// Waiting for all goroutines to finish (otherwise they die as main routine dies)
	wg.Wait()
	return nil
}

func UpdateApps() error {
	apps, err := steam.GetApps()
	if err != nil {
		return err
	}
	err = storage.UpdateApps(apps)
	return err
}
