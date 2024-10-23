package domain

import (
	"fmt"
	"github.com/gofiber/fiber/v3/log"
	"runtime"
	"time"
)

func Collector() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	go func() {
		for range time.Tick(10 * time.Minute) {
			runtime.ReadMemStats(&m)

			text0 := "Starting Garbage Collector..."

			log.Debug(text0)
			Logger.Debug(text0)

			text1 := fmt.Sprintf("Before GC: HeapAlloc = %.2f MB, TotalAlloc = %.2f MB, NumGC = %v",
				float64(m.HeapAlloc)/1024/1024, float64(m.TotalAlloc)/1024/1024, m.NumGC)

			log.Debug(text1)
			Logger.Debug(text1)

			runtime.GC()

			runtime.ReadMemStats(&m)

			text2 := fmt.Sprintf("After GC: HeapAlloc = %.2f MB, TotalAlloc = %.2f MB, NumGC = %v",
				float64(m.HeapAlloc)/1024/1024, float64(m.TotalAlloc)/1024/1024, m.NumGC)

			log.Debug(text2)
			Logger.Debug(text2)

			text3 := "Garbage Collector executed!"

			log.Debug(text3)
			Logger.Debug(text3)
		}
	}()

}
