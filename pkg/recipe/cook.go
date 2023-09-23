package recipe

const (
	// CertResolver is the name of the cert-manager resolver
	CertResolver = "certResolver"
)

type Hooks struct {
	Pre  PreHook
	Post PostHook
}
type PreHook func() error
type PostHook func() error

// Cook runs recipe
func Cook(recipePath string, hooks ...Hooks) error {
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
		if err := hook.Pre(); err != nil {
			return err
		}
	}

	// add steps
	if err := add(recipe,
		step{f: checkPrerequisites, c: recipe.Ingredients.Node.Check},
		step{f: installK3s, c: recipe.Ingredients.K3s.Enabled},
		step{f: createSecrets, c: recipe.Ingredients.Secrets.Enabled},
		step{f: installCertManager, c: recipe.Ingredients.CertManager.Enabled},
		step{f: installTraefik, c: recipe.Ingredients.Traefik.Enabled},
		step{f: installArgocd, c: recipe.Ingredients.ArgoCD.Enabled},
		step{f: configureGitopsProjects, c: recipe.Ingredients.ArgoCD.Enabled},
		step{f: printKubeconfig, c: recipe.Debug},
	); err != nil {
		return err
	}

	// post cook hooks
	for _, hook := range hooks {
		if err := hook.Post(); err != nil {
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
