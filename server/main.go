package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"server/http"
	"server/logger"
	"time"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	printCurrentTime()
	go http.StartServer(ctx)
	// go http.Crawler()

	<-http.StopChannel
	fmt.Println("call cancelFunc(), stop all goroutine")
	cancelFunc()
	fmt.Println("Process will shutdown after 5 secones")
	time.Sleep(5 * time.Second)
	fmt.Println("Process finished")
}

func printCurrentTime() {
	now := time.Now() //获取当前时间
	log.Printf("now:%v\n", now)
	log.Printf("%d-%d-%d %d:%d:%d %d\n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.UnixMilli()%1000)
	var (
		full  = "02-01- 03:04:05.000 PM Mon Jan"  // 12H AM PM
		full2 = "2006-01-02 15:04:05.000 Mon Jan" // 24H
		day   = "2006-01-02"
		day2  = "2006/01/02"
	)
	logger.GetLogger().LogAttrs(
		context.Background(),
		slog.LevelDebug,
		fmt.Sprintf("format as full %s\n", now.Format(full)),
	)
	log.Printf("format as full2 %s\n", now.Format(full2))
	log.Printf("format as day %s\n", now.Format(day))
	log.Printf("format as day2 %s\n", now.Format(day2))
	timeZone, timeZoneOffSet := now.Zone()
	log.Printf("timezone: %s, offset: %d\n", timeZone, timeZoneOffSet/60/60)
	log.Printf("current timestamp1:%v\n", now.Unix())      // ex: 1714115113
	log.Printf("current timestamp2:%v\n", now.UnixMilli()) // ex: 1714115113640
	log.Printf("current timestamp3:%v\n", now.UnixNano())  // ex: 1714115113640244400

	nextHour := now.Add(time.Hour)
	diff := nextHour.Sub(now)
	log.Printf("diff %v\n", diff)

	// There is no sub for time.Time, we can add negative Duration.
	previousHour := now.Add(-diff)
	log.Println("Previous Hour", previousHour)
}
