package main

type ConfigFile struct {
	Listen         string   `json:"listen"`
	StorageServers []string `json:"storage_servers"`
	BaseDir        string   `json:"base_dir"`
	CountSrv       int
}

type fileInfo struct {
	CrcSum      string `json:"crc_sum"`
	FirstServer int    `json:"first_server"`
}
