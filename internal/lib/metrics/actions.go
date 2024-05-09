package metrics

func (m metricsImpl) IncrementRestartedContainers() {
	m.restartedContainers.Inc()
}
func (m metricsImpl) SetDisabledContainers(count float64) {
	m.disabledContainers.Set(count)
}
func (m metricsImpl) SetRunningContainers(count float64) {
	m.runningContainers.Set(count)
}
func (m metricsImpl) SetUnhealthyContainers(count float64) {
	m.unhealthyContainers.Set(count)
}
func (m metricsImpl) SetIgnoredContainers(count float64) {
	m.ignoredContainers.Set(count)
}
