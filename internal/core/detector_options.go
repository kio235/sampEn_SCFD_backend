package core

type DetectorOption func(*Detector)

func WithM(m int) DetectorOption {
	return func(d *Detector) { d.m = m }
}

func WithRCoeff(rCoeff float64) DetectorOption {
	return func(d *Detector) { d.rCoeff = rCoeff }
}

func WithWindowSize(windowSize int) DetectorOption {
	return func(d *Detector) { d.windowSize = windowSize }
}

func WithStepSize(stepSize int) DetectorOption {
	return func(d *Detector) { d.stepSize = stepSize }
}

func WithRelativeThreshold(relativeThreshold float64) DetectorOption {
	return func(d *Detector) { d.relativeThreshold = relativeThreshold }
}

func WithAbsoluteThreshold(absoluteThreshold float64) DetectorOption {
	return func(d *Detector) { d.absoluteThreshold = absoluteThreshold }
}

func WithFaultThreshold(faultThreshold int) DetectorOption {
	return func(d *Detector) { d.faultThreshold = faultThreshold }
}
