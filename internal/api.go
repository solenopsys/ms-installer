package internal

import (
	bl_kubernetes_tools "github.com/solenopsys/bl-kubernetes-tools"
	"google.golang.org/protobuf/proto"
	"hs-installer/pkg"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

type RequestProcessor struct {
	helmApi *pkg.HelmApi
}

func (r *RequestProcessor) getInstalled(message []byte) []byte {
	req := &pkg.GetChartsRequest{}
	proto.Unmarshal(message, req)

	chartsList, err := r.helmApi.ListHelmCharts("")
	if err != nil {
		klog.Infof("Error get list charts: %s", err)
	}

	var charts = make([]*pkg.Chart, 0)

	for _, chart := range chartsList.Items {
		charts = append(charts, &pkg.Chart{
			Name:       chart.ObjectMeta.Name,
			Digest:     string(chart.ObjectMeta.UID),
			Repository: chart.Spec.Repo,
			Version:    chart.Spec.Version,
		})

	}

	klog.Infof("RETURN CHARTS %s", charts)

	res := &pkg.ChartsResponse{Charts: charts}
	marshal, err := proto.Marshal(res)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal
}

func (r *RequestProcessor) installChart(message []byte) []byte {
	req := &pkg.InstallChartRequest{}
	proto.Unmarshal(message, req)

	chart := req.Chart
	r.helmApi.CreateHelmChartSimple(chart.Name, chart.Repository, chart.Version)

	klog.Infof("INSTALL CHART %s", chart)

	res := &pkg.OperationStatus{Status: "CHART_CREATED"}
	marshal, err := proto.Marshal(res)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal
}

func (r *RequestProcessor) uninstallChart(message []byte) []byte {
	req := &pkg.UninstallChartRequest{}
	proto.Unmarshal(message, req)

	r.helmApi.DeleteHelmChart(req.Digest)

	klog.Infof("UNINSTALL CHART %s", req.Digest)

	res := &pkg.OperationStatus{Status: "CHART_DELETED"}
	marshal, err := proto.Marshal(res)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal
}

func ProcessingFunction() func(message []byte, functionId uint8) []byte {

	var rc *rest.Config
	var err error
	rc, err = bl_kubernetes_tools.GetCubeConfig(false)
	if err != nil {
		klog.Error("error getting Kubernetes config:", err)
	}

	helmApi := pkg.NewAPI(rc)
	requestProcessor := RequestProcessor{helmApi: &helmApi}
	return func(message []byte, functionId uint8) []byte {
		klog.Infof("Pocessing function %s", functionId)

		if functionId == 1 {
			return requestProcessor.getInstalled(message)
		}
		if functionId == 2 {
			return requestProcessor.installChart(message)
		}
		if functionId == 3 {
			return requestProcessor.uninstallChart(message)
		}

		return []byte("FUNCTION_NOT_FOUND")
	}
}
