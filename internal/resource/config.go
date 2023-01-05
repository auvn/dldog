package resource

type DownloadItem struct {
	Name        string
	Destination string
	Fetch       struct {
		Git *struct {
			Url    string
			Tag    string
			Sha    string
			Branch string
			Cwd    string
			Files  []string
		}
		ArchiveUrl *struct {
			Url    string
			Format string
			Cwd    string
			Files  []string
		}
		FileUrl string
	}
}

type SkipConfig struct {
	Files []string
}
