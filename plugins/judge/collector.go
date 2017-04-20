package judge

//Collector is
type Collector interface {
	Collect() int
}

//NewCollector is
func NewCollector(name string) Collector {
	return &ZkrCollector{}
}
