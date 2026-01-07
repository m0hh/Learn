package cralwer



type linkFetcher struct {
	urlGetter URLGetter
	netDetector PrivateNetwrokDetector
}


func NewLinkFetcher(urlGetter URLGetter, netDetector PrivateNetwrokDetector) *linkFetcher {
	return &linkFetcher{
		urlGetter: urlGetter,
		netDetector: netDetector,
	}
}

