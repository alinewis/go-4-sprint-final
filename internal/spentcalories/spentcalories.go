package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("incorrect number of arguments")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	action := parts[1]

	dur, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	return steps, action, dur, nil
}

func distance(steps int, height float64) float64 {
	lenOfStep := height * stepLengthCoefficient
	distanceM := lenOfStep * float64(steps)
	return distanceM / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distanceKm := distance(steps, height)
	return distanceKm / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, action, dur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var result string

	switch action {
	case "Бег":
		cal, err := RunningSpentCalories(steps, weight, height, dur)
		if err != nil {
			return "", fmt.Errorf("ошибка при подсчёте калорий: %w", err)
		}

		hours := dur.Hours()
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, dur)

		result = fmt.Sprintf("Тип тренировки: %s\n", action) +
			fmt.Sprintf("Длительность: %.2f ч.\n", hours) +
			fmt.Sprintf("Дистанция: %.2f км.\n", dist) +
			fmt.Sprintf("Скорость: %.2f км/ч\n", speed) +
			fmt.Sprintf("Сожгли калорий: %.2f\n", cal)

	case "Ходьба":
		cal, err := WalkingSpentCalories(steps, weight, height, dur)
		if err != nil {
			return "", fmt.Errorf("ошибка при подсчёте калорий: %w", err)
		}

		hours := dur.Hours()
		dist := distance(steps, height)
		speed := meanSpeed(steps, height, dur)

		result = fmt.Sprintf("Тип тренировки: %s\n", action) +
			fmt.Sprintf("Длительность: %.2f ч.\n", hours) +
			fmt.Sprintf("Дистанция: %.2f км.\n", dist) +
			fmt.Sprintf("Скорость: %.2f км/ч\n", speed) +
			fmt.Sprintf("Сожгли калорий: %.2f\n", cal)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input parameters")
	}

	speed := meanSpeed(steps, height, duration)
	return weight * speed * duration.Minutes() / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input parameters")
	}

	speed := meanSpeed(steps, height, duration)
	calories := weight * speed * duration.Minutes() / minInH
	return calories * walkingCaloriesCoefficient, nil
}
