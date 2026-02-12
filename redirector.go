package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"gopkg.in/yaml.v3"
)

func __[T any](item T) *T {
	return &item
}

func ___[T any](item *T) T {
	return *item
}

// Returns a [PolicySummary].
func handleDOIData(doi string, crapi *crossrefapi.CrossRefClient) PolicySummary {

	// Define the default site policy
	policy := PolicySummary{
		Captcha: __(EnforcedIfMassTraffic),
		ContentAccess: __(ContentAccessPolicy{
			PlatformSubscription: __(EnforcedIfPrivateIntended),
			ContentPayment:       __(EnforcedIfPrivateIntended),
		}),
	}

	if crapi == nil {
		return policy
	}

	log.Info("Attempting to search CrossRef for " + doi)
	works, err := crapi.Works(doi)
	if err != nil {
		log.Error(err)
		return policy
	}

	if works.Message != nil && works.Message.License != nil {
		for _, license := range works.Message.License {
			if strings.Contains(license.URL, "creativecommons.org") {
				policy.ContentAccess.AuthorSubscription = nil
				policy.ContentAccess.ContentPayment = nil
			}
		}
	}

	return policy
}

func handleGenericData(url *url.URL) *PolicySummary {
	if url == nil {
		return nil
	}

	hostname := resolveHostnameAlias(url.Host)

	rawRulesFile, err := os.ReadFile(fmt.Sprintf("./data/rules/%s.yaml", strings.ReplaceAll(hostname, "*", "$")))
	if err != nil {
		log.Error(err)
		return nil
	}

	var policy PolicySummary
	err = yaml.Unmarshal(rawRulesFile, &policy)
	if err != nil {
		log.Error(err)
		return nil
	}

	a, _ := json.Marshal(policy)
	log.Info(string(a))

	return __(policy)
}

func handleRedirect(c fiber.Ctx, crapi *crossrefapi.CrossRefClient) {
	targetUrl := c.Params("*")
	cfg, cfgErr := getAppConfig()

	// Step 0: Validate the URL
	parsedTargetUrl, err := url.Parse(targetUrl)
	if err != nil || parsedTargetUrl == nil {
		c.Status(400)
	}

	var policy *PolicySummary = nil

	// Step 1: Check if this is a DOI
	if parsedTargetUrl.Host == "doi.org" {
		policy = __(handleDOIData(parsedTargetUrl.Path, crapi))
	} else {
		policy = handleGenericData(parsedTargetUrl)
	}

	if policy == nil && cfgErr != nil && cfg.SkipIfNoKnownPolicies == true {
		c.Redirect().Status(302).To(targetUrl)
	} else {
		err := c.Render("redirect", fiber.Map{
			"Ctx": c,

			"AppName":        cfg.AppName,
			"OriginAppName":  c.Params("Origin", "this website"),
			"TargetHostname": parsedTargetUrl.Host,
			"TargetURL":      targetUrl,
			"Policy":         policy,

			// EnforcementUrgency data
			"OptionallyEnforced":        __(OptionallyEnforced),
			"EnforcedIfMassTraffic":     __(EnforcedIfMassTraffic),
			"EnforcedIfPrivateIntended": __(EnforcedIfPrivateIntended),
			"EnforcedIfAutoNSFW":        __(EnforcedIfAutoNSFW),
			"EnforcedIfInteracting":     __(EnforcedIfInteracting),
			"EnforcedIfConsuming":       __(EnforcedIfConsuming),
			"AlwaysEnforced":            __(AlwaysEnforced),
		})
		if err != nil {
			log.Error(err)
			c.Status(500)
		}
	}
}
