package nginx

var (
	sitesAvailable = "/etc/nginx/sites-available/"
	sitesEnabled   = "/etc/nginx/sites-enabled/"
)

type tmplConfig struct {
	Domain string
	Cert   string
	Target string
}
