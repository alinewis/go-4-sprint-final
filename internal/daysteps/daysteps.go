package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, spentcalories.ErrNumbOfArgs
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, spentcalories.ErrStepsNeg
	}

	dur, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	if dur <= 0 {
		return 0, 0, spentcalories.ErrDurNeg
	}

	return steps, dur, nil

}

func DayActionInfo(data string, weight, height float64) string {
	steps, dur, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distM := stepLength * float64(steps)
	distKm := distM / float64(mInKm)

	cal, err := spentcalories.WalkingSpentCalories(steps, weight, height, dur)
	if err != nil {
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\n", steps) +
		fmt.Sprintf("Дистанция составила %.2f км.\n", distKm) +
		fmt.Sprintf("Вы сожгли %.2f ккал.\n", cal)

	return result
}
