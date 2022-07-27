package main

import (
	"google.golang.org/protobuf/proto"
	"k8s.io/klog/v2"
	"strings"
)

func reposResponse(filter string) []byte {
	reposFromKube := loadReposFromKube()

	repos := &ReposResponse{}

	for nameRepo, urlRepo := range reposFromKube {
		repo := &Repo{}
		if "" == filter || strings.Contains(urlRepo, filter) || strings.Contains(nameRepo, filter) {
			repo.Name = nameRepo
			repo.Url = urlRepo
			repos.Repos = append(repos.Repos, repo)
		}
	}

	marshal, err := proto.Marshal(repos)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal
}

func getRepos(message []byte) []byte {
	req := &GetReposRequest{}
	proto.Unmarshal(message, req)

	filter := req.GetFilter()

	return reposResponse(filter)
}

func addRepo(message []byte) []byte {
	req := &AddRepoRequest{}
	proto.Unmarshal(message, req)

	repos := loadReposFromKube()
	repos[req.Url] = ""

	saveReposToKube(repos)

	res := &OperationStatus{Status: true}

	marshal, err := proto.Marshal(res)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal
}

func cacheConversion(filter string, cached map[string]*map[string][]HelmChart) []*Chart {
	var result = []*Chart{}
	for repoName, charts := range cached {
		for _, chartArray := range *charts {
			chartRest := chartArray[0]

			layout := "2006-01-02T15:04:05-0700"
			ch := Chart{
				Name:        chartRest.Name,
				Version:     chartRest.Version,
				Repo:        repoName,
				Description: chartRest.Description,
				Icon:        chartRest.Icon,
				Created:     chartRest.Created.Format(layout),
				Digest:      chartRest.Digest,
			}
			result = append(result, &ch)
		}
	}

	return result
}

func getCharts(message []byte) []byte {
	req := &GetChartsRequest{}
	proto.Unmarshal(message, req)
	repos := loadReposFromKube()

	cached := cachedReq(req.Reload, repos)

	response := &ChartsResponse{}

	conversion := cacheConversion("", cached)

	response.Charts = conversion

	marshal, err := proto.Marshal(response)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal

}

func getInstalled(message []byte) []byte {
	return []byte("NOT_IMPLEMENTED")
}

func installChart(message []byte) []byte {
	req := &InstallChartRequest{}
	proto.Unmarshal(message, req)

	repos := loadReposFromKube()

	cached := cachedReq(false, repos)
	for repoName, charts := range cached {
		for _, chartArray := range *charts {
			chartRest := chartArray[0]

			if chartRest.Digest == req.Digest {
				createChart(chartRest.Name, repos[repoName], chartRest.Version)
			}
		}
	}
	klog.Infof("INSTALL CHART %s")

	res := &OperationStatus{Status: true}
	marshal, err := proto.Marshal(res)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal
}

func uninstallChart(message []byte) []byte {
	return []byte("NOT_IMPLEMENTED")
}

func removeRepo(message []byte) []byte {
	req := &RemoveRepoRequest{}
	proto.Unmarshal(message, req)

	repos := loadReposFromKube()
	delete(repos, req.Name)
	saveReposToKube(repos)

	res := &OperationStatus{Status: true}
	marshal, err := proto.Marshal(res)

	if err != nil {
		klog.Infof("Marshal error: %s", err)
	}
	return marshal
}

func processingFunction() func(message []byte, streamId uint32, serviceId uint16, functionId uint16) []byte {
	return func(message []byte, streamId uint32, serviceId uint16, functionId uint16) []byte {
		if functionId == 1 {
			return getRepos(message)
		}
		if functionId == 2 {
			return addRepo(message)
		}
		if functionId == 3 {
			return removeRepo(message)
		}
		if functionId == 4 {
			return getCharts(message)
		}
		if functionId == 5 {
			return getInstalled(message)
		}
		if functionId == 6 {
			return installChart(message)
		}
		if functionId == 7 {
			return uninstallChart(message)
		}

		return []byte("FUNCTION_NOT_FOUND")
	}
}
