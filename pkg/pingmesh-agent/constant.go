package pingmesh_agent

const (
	DefaultPingmeshServerAddr = "0.0.0.0"
	DefaultPingmeshServerPort = "10009"
	DefaultPingmeshDownloadURL = "/getPinglist"
	DefaultPingmeshUploadURL   = "/uploadMetrics"
	DefaultPingmeshServiceAddr = "pms-service"

	// ping
	MetricsNamePingLatency       = `ping_latency_millonseconds`
	MetricsNamePingPackageDrop   = `ping_packageDrop_rate`
	MetricsNamePingTargetSuccess = `ping_target_success`

	// http
	MetricsNameHttpResolvedurationMillonseconds    = `http_resolveDuration_millonseconds`
	MetricsNameHttpTlsDurationMillonseconds        = `http_tlsDuration_millonseconds`
	MetricsNameHttpConnectDurationMillonseconds    = `http_connectDuration_millonseconds`
	MetricsNameHttpProcessingDurationMillonseconds = `http_processingDuration_millonseconds`
	MetricsNameHttpTransferDurationMillonseconds   = `http_transferDuration_millonseconds`
	MetricsNameHttpInterfaceSuccess                = `http_interface_success`

)
