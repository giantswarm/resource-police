package project

var (
	description = "Command line tool for reporting test clusters that live too long."
	gitSHA      = "n/a"
	name        = "resource-police"
	source      = "https://github.com/giantswarm/resource-police"
	version     = "1.2.0"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
