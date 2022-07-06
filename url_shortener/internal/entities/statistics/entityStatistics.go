package statistics

type Statistics struct {
	IP     string `json:"ip"`
	Viewed string `json:"viewed"`
	Count  int64  `json:"count"`
	Long   string `json:"long"`
	Short  string `json:"short"`
	Admin  string `json:"admin"`
}
