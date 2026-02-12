package main

import (
	"errors"
	"html/template"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/caltechlibrary/crossrefapi"
	"github.com/gofiber/contrib/v3/i18n"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v3"
	ni18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type AppConfig struct {
	// Do you want to access the CrossRef API (under the polite policy)?
	ApiCrossrefEmail *string `env:"API_CROSSREF_EMAIL"`

	// Do you want to give a different app name than "SaferLinks"?
	AppName string `env:"APP_NAME" envDefault:"SaferLinks"`

	// Do you want to set a custom hostname?
	HostAndPort string `env:"HOST_AND_PORT" envDefault:":8080"`

	// Do you need to continue redirect if the target site does not have any known restrictions?
	SkipIfNoKnownPolicies bool `env:"SKIP_IF_NO_KNOWN_POLICIES" envDefault:"false"`
}

var translator = i18n.New(&i18n.Config{
	RootPath: "./i18n",
	AcceptLanguages: []language.Tag{
		language.English,
		language.Indonesian,
	},
	DefaultLanguage: language.English,
})

func getAppConfig() (AppConfig, error) {
	var cfg AppConfig
	err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: false,
	})
	return cfg, err
}

func formatList(items []string, delimiter string, lastDelimiter string) string {
	if len(items) == 0 {
		return ""
	}

	if len(items) == 1 {
		return items[0]
	}

	// Join all but the last item with commas
	allButLast := items[:len(items)-1]
	last := items[len(items)-1]

	// Use a strings.Builder for efficiency
	var builder strings.Builder
	builder.WriteString(strings.Join(allButLast, delimiter))
	builder.WriteString(lastDelimiter)
	builder.WriteString(last)

	return builder.String()
}

func getRegionsString(c fiber.Ctx, policy PolicySummary) string {
	regionItems := []string{}
	for _, region := range policy.Regions {
		regionItems = append(regionItems, translator.MustLocalize(c, "country_"+region))
	}

	return formatList(regionItems, translator.MustLocalize(c, "delimiters.comma_space"), translator.MustLocalize(c, "delimiters.comma_or_space"))
}

func formatPolicyList(c fiber.Ctx, rules []string, policy PolicySummary, delimiter string, lastDelimiter string) string {
	if len(rules) == 0 {
		return ""
	}

	renderedItems := []string{}
	for i, rule := range rules {
		var i18nName string
		if i == 0 {
			i18nName = rule + ".phrase_start"
		} else {
			i18nName = rule + ".phrase"
		}
		renderedItems = append(renderedItems, translator.MustLocalize(c, &ni18n.LocalizeConfig{
			MessageID: i18nName,
			TemplateData: map[string]string{
				"Name":    policy.Name,
				"Regions": getRegionsString(c, policy),
			},
		}))
	}

	if len(renderedItems) == 1 {
		return renderedItems[0]
	}

	// Return as template.HTML to prevent double escaping in the template
	return formatList(renderedItems, delimiter, lastDelimiter)
}

func getMandatoryEssentialsFromPolicy(policy PolicySummary) []string {
	list := []string{}

	if policy.ContentAccess != nil {
		if policy.ContentAccess.PlatformSubscription != nil && *policy.ContentAccess.PlatformSubscription == AlwaysEnforced {
			list = append(list, "content_access.platform_subscription")
		}
		if policy.ContentAccess.AuthorSubscription != nil && *policy.ContentAccess.AuthorSubscription == AlwaysEnforced {
			list = append(list, "content_access.author_subscription")
		}
		if policy.ContentAccess.ContentPayment != nil && *policy.ContentAccess.ContentPayment == AlwaysEnforced {
			list = append(list, "content_access.content_payment")
		}
		if policy.ContentAccess.ContentPassword != nil && *policy.ContentAccess.ContentPassword == AlwaysEnforced {
			list = append(list, "content_access.content_password")
		}
	}

	if len(list) < 2 && policy.ClientApp != nil {
		if policy.ClientApp.DedicatedHardware != nil && *policy.ClientApp.DedicatedHardware == AlwaysEnforced {
			list = append(list, "client_app.dedicated_hardware")
		}
		if policy.ClientApp.DesktopApp != nil && *policy.ClientApp.DesktopApp == AlwaysEnforced {
			list = append(list, "client_app.desktop_app")
		}
		if policy.ClientApp.MobileApp != nil && *policy.ClientApp.MobileApp == AlwaysEnforced {
			list = append(list, "client_app.mobile_app")
		}
		if policy.ClientApp.MatchRegion != nil && *policy.ClientApp.MatchRegion == AlwaysEnforced {
			list = append(list, "client_app.match_region")
		}
		if policy.ClientApp.MustLatestClientSoftware != nil && *policy.ClientApp.MustLatestClientSoftware == AlwaysEnforced {
			list = append(list, "client_app.must_latest_client_software")
		}
		if policy.ClientApp.MustLatestHostSoftware != nil && *policy.ClientApp.MustLatestHostSoftware == AlwaysEnforced {
			list = append(list, "client_app.must_latest_host_software")
		}
	}

	if len(list) < 2 && policy.Account != nil {
		if policy.Account.LoggedIn != nil && *policy.Account.LoggedIn == AlwaysEnforced {
			list = append(list, "account.logged_in")
		}
		if policy.Account.MatchRegion != nil && *policy.Account.MatchRegion == AlwaysEnforced {
			list = append(list, "account.match_region")
		}
	}

	return list
}

func getCrossRefApiClient() (*crossrefapi.CrossRefClient, error) {
	cfg, err := getAppConfig()
	if err == nil && cfg.ApiCrossrefEmail != nil && len(*cfg.ApiCrossrefEmail) > 0 {
		log.Info("Attempting to create a CrossRef client")
		crapi, err := crossrefapi.NewCrossRefClient(cfg.AppName, *cfg.ApiCrossrefEmail)

		if err != nil {
			log.Error("Error when creating CrossRef client")
		}
		return crapi, err
	}
	return nil, errors.New("CrossRef email not supplied")
}

func main() {
	engine := html.New("./views", ".html")

	funcMap := template.FuncMap{
		"TP": func(c interface{}, msg string, policy PolicySummary) template.HTML {
			return template.HTML(translator.MustLocalize(c.(fiber.Ctx), &ni18n.LocalizeConfig{
				MessageID: msg,
				TemplateData: map[string]string{
					"Name":    policy.Name,
					"Regions": getRegionsString(c.(fiber.Ctx), policy),
				},
			}))
		},
		"T": func(c interface{}, msg interface{}) string {
			return translator.MustLocalize(c.(fiber.Ctx), msg)
		},
		"derefInt": ___[int],
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"decodeEnforcementUrgency": func(status *EnforcementUrgency) string {
			if status == nil || *status <= 0 {
				return ""
			} else {
				str, _ := status.MarshalYAML()
				str2, _ := strings.CutSuffix(string(str), "\n")
				return str2
			}
		},
		"getColorSchemeFromEnforcementUrgency": func(status *EnforcementUrgency) string {
			if status == nil || *status <= 0 {
				return "light"
			} else if *status == OptionallyEnforced {
				return "secondary"
			} else if *status < AlwaysEnforced {
				return "warning"
			} else {
				return "danger"
			}
		},
		"getMandatoryEssentialsFromPolicy": getMandatoryEssentialsFromPolicy,
		"formatPolicyList":                 formatPolicyList,
		"renderRequirementsDetail": func(c fiber.Ctx, policy PolicySummary, url string) template.HTML {
			return template.HTML(translator.MustLocalize(c, &ni18n.LocalizeConfig{
				MessageID: "redirect_page.policies_intro",
				TemplateData: map[string]string{
					"Name": policy.Name,
					"Url":  url,
				},
			}))
		},
		"renderRequirementsWarning": func(c fiber.Ctx, policy PolicySummary) template.HTML {
			return template.HTML(translator.MustLocalize(c, &ni18n.LocalizeConfig{
				MessageID: "redirect_page.title_warning",
				TemplateData: map[string]string{
					"Requirements": formatPolicyList(
						c,
						getMandatoryEssentialsFromPolicy(policy),
						policy,
						translator.MustLocalize(c, "delimiters.comma_space"),
						translator.MustLocalize(c, "delimiters.comma_and_space"),
					),
				},
			}))
		},
	}
	engine.AddFuncMap(funcMap)

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	cfg, err := getAppConfig()
	if err != nil {
		log.Error(err)
	}

	crapi, err := getCrossRefApiClient()

	// Always empty crapi if failed
	if err != nil {
		crapi = nil
	}

	app.Get("/static/*", static.New("./static"))
	app.Get("/redirect/*", func(c fiber.Ctx) {
		handleRedirect(c, crapi)
	})

	app.Listen(cfg.HostAndPort)
}
