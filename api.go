package main

import (
	"google.golang.org/protobuf/proto"
	"k8s.io/klog/v2"
)

func getInstalled(message []byte) []byte {
	return []byte("NOT_IMPLEMENTED")
}

func installChart(message []byte) []byte {
	req := &InstallChartRequest{}
	proto.Unmarshal(message, req)

	chart := req.Chart
	createChart(chart.Name, chart.Repository, chart.Version)

	klog.Infof("INSTALL CHART %s")

	res := &OperationStatus{Status: "CHART_CREATED"}
	marshal, err := proto.Marshal(res)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal
}

func uninstallChart(message []byte) []byte {
	return []byte("NOT_IMPLEMENTED")
}

func processingFunction() func(message []byte, functionId uint8) []byte {
	return func(message []byte, functionId uint8) []byte {
		klog.Infof("Pocessing function %s", functionId)

		if functionId == 1 {
			return getInstalled(message)
		}
		if functionId == 2 {
			return installChart(message)
		}
		if functionId == 3 {
			return uninstallChart(message)
		}

		return []byte("FUNCTION_NOT_FOUND")
	}
}
