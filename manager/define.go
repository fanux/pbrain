package manager

var (
	Host string
	Port string

	DBHost   string
	DBPort   string
	DBUser   string
	DBName   string
	DBPasswd string

	DockerHost string

	AllowedDomain string
)

type Node struct {
	ID      string
	IP      string
	Addr    string
	Name    string
	Cpus    int64
	Memory  int64
	Labels  map[string]string
	Version string
}
