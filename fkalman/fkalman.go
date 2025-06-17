package fkalman

type KalmanFilter struct {
	X float64 // текущее состояние (оценка азимута)
	P float64 // ошибка оценки
	Q float64 // шум процесса
	R float64 // шум измерения
}

func NewKalmanFilter(initialEstimate, processNoise, measurementNoise float64) *KalmanFilter {
	return &KalmanFilter{
		X: initialEstimate,
		P: 1.0,
		Q: processNoise,
		R: measurementNoise,
	}
}

func (kf *KalmanFilter) Update(measurement float64) float64 {
	// Предсказание
	kf.P += kf.Q

	// Коэффициент Калмана
	K := kf.P / (kf.P + kf.R)

	// Обновление оценки
	kf.X += K * angleDiff(measurement, kf.X)

	// Обновление ошибки
	kf.P *= (1 - K)

	return kf.X
}

// Корректная разница между углами с учётом круга
func angleDiff(a, b float64) float64 {
	d := a - b
	for d > 180 {
		d -= 360
	}
	for d < -180 {
		d += 360
	}
	return d
}
