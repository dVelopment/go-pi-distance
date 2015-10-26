package distance

import (
  "fmt"
  "time"
  "github.com/stianeikeland/go-rpio"
)

const SPEED_OF_SOUND = 34320 // in cm per second

var (
  echoPin rpio.Pin
  triggerPin rpio.Pin
)

func ReadAverageDistance(numberOfReads int) (float64) {
  if (numberOfReads < 1) {
    numberOfReads = 1
  }
  limit := float64(numberOfReads)

  sum := 0.0
  reads := 0.0

  for reads < limit {
    sum = sum + ReadDistance()
    reads = reads + 1.0
    time.Sleep(10 * time.Millisecond)
  }

  return sum / reads
}

func ReadDistance() (float64) {
  var res rpio.State
  var start, stop time.Time
  var duration int
  fmt.Println("Read distance in cm")

  // trigger a 10µs burst
  fmt.Printf("Trigger burst for %dns\n", 10 * time.Microsecond)
  triggerPin.High()
  time.Sleep(10 * time.Microsecond)
  triggerPin.Low()

  // wait for echo
  fmt.Println("Wait for echo ...")
  res = echoPin.Read()
  for (res == rpio.Low) {
    start = time.Now()
    res = echoPin.Read()
  }

  fmt.Println("Wait for echo end ...")

  for (res == rpio.High) {
    stop = time.Now()
    res = echoPin.Read()
  }

  duration = stop.Nanosecond() - start.Nanosecond()

  fmt.Printf("duration: %d\n (start: %d - stop: %d)\n", duration, stop.Nanosecond(), start.Nanosecond())

  distance := float64(SPEED_OF_SOUND) / 2.0 * float64(duration) / 1000000000.0

  return distance
}

func Init(echo int, trigger int) (err error) {
  err = rpio.Open()

  if (err != nil) {
    return
  }

  fmt.Printf("echo pin: %d - trigger pin: %d\n", echo, trigger)

  echoPin = rpio.Pin(echo)
  echoPin.Input()

  triggerPin = rpio.Pin(trigger)
  triggerPin.Output()

  // ensure that triggerPin is set to low
  triggerPin.Low()

  fmt.Println("Waiting for sensor to settle")
  time.Sleep(2 * time.Second)

  return nil
}

func Close() {
  rpio.Close()
}
