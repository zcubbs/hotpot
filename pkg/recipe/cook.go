package recipe

const (
	// CertResolver is the name of the cert-manager resolver
	CertResolver = "certResolver"
)

type Hooks struct {
	Pre  PreHook
	Post PostHook
}
type PreHook func(r *Recipe) error
type PostHook func(r *Recipe) error

// Cook runs recipe
func Cook(recipePath string, deps Dependencies, hooks ...Hooks) error {
	// load config
	recipe, err := Load(recipePath)
	if err != nil {
		return err
	}

	// debug recipe
	if recipe.Debug {
		printRecipe(recipe)
	}

	// validate config
	if err := validate(recipe); err != nil {
		return err
	}

	// preheat hooks
	for _, hook := range hooks {
		if err := hook.Pre(recipe); err != nil {
			return err
		}
	}

	// add steps
	if err := add(recipe,
		step{f: func(r *Recipe) error { return checkPrerequisites(r, deps.SystemInfo) }, c: recipe.Node.Check},
		step{f: func(r *Recipe) error { return installK3s(r, deps.K3s, deps.Helm, deps.FileSystem) }, c: recipe.K3s.Enabled},
		step{f: func(r *Recipe) error { return installK9s(r, deps.K9s) }, c: recipe.K9s.Enabled},
		step{f: createSecrets, c: recipe.Secrets.Enabled},
		step{f: func(r *Recipe) error { return installCertManager(r, deps.CertManager) }, c: recipe.CertManager.Enabled},
		step{f: func(r *Recipe) error { return installTraefik(r, deps.Traefik) }, c: recipe.Traefik.Enabled},
		step{f: func(r *Recipe) error { return installRancher(r, deps.Rancher) }, c: recipe.Rancher.Enabled},
		step{f: func(r *Recipe) error { return installArgocd(r, deps.ArgoCD) }, c: recipe.ArgoCD.Enabled},
		step{f: configureGitopsProjects, c: recipe.Gitops.Enabled},
		step{f: printKubeconfig, c: recipe.Debug},
	); err != nil {
		return err
	}

	// post cook hooks
	for _, hook := range hooks {
		if err := hook.Post(recipe); err != nil {
			return err
		}
	}

	return nil
}

func add(r *Recipe, steps ...step) error {
	for _, step := range steps {
		if !step.c {
			continue
		}
		if err := step.f(r); err != nil {
			return err
		}
	}

	return nil
}
