package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type EnforcementUrgency int

const (
	// NO
	_not_enforced EnforcementUrgency = iota

	// MAYBE; This policy is optionally enforced; users may fully opt-out from behavior.
	OptionallyEnforced

	// YES; This policy is only enforced under mass traffic.
	EnforcedIfMassTraffic

	// YES; This policy is only enforced when the content is private; public content may have no higher restrictions.
	EnforcedIfPrivateIntended

	// YES; This policy is especially enforced when the content is automatically assumed to be NSFW, regardless of explicit public/private access settings.
	EnforcedIfAutoNSFW

	// YES; This policy is only enforced when attempting to perform certain actions, such as clicking a button, or editing a shared document.
	EnforcedIfInteracting

	// YES; This policy is only enforced when the user attempts to view more parts of the content; such as scrolling down the page.
	EnforcedIfConsuming

	// YES; This policy is non-negotiable by the platform owner.
	AlwaysEnforced
)

func (u EnforcementUrgency) MarshalJSON() ([]byte, error) {
	switch u {
	case OptionallyEnforced:
		return json.Marshal("optionally_enforced")
	case EnforcedIfMassTraffic:
		return json.Marshal("enforced_if_mass_traffic")
	case EnforcedIfPrivateIntended:
		return json.Marshal("enforced_if_private_intended")
	case EnforcedIfAutoNSFW:
		return json.Marshal("enforced_if_auto_nsfw")
	case EnforcedIfInteracting:
		return json.Marshal("enforced_if_interacting")
	case EnforcedIfConsuming:
		return json.Marshal("enforced_if_consuming")
	case AlwaysEnforced:
		return json.Marshal("always_enforced")
	}
	return []byte{}, errors.New(fmt.Sprintf("Invalid EnforcementUrgency value: %v", u))
}

func (u EnforcementUrgency) MarshalYAML() ([]byte, error) {
	switch u {
	case OptionallyEnforced:
		return yaml.Marshal("optionally_enforced")
	case EnforcedIfMassTraffic:
		return yaml.Marshal("enforced_if_mass_traffic")
	case EnforcedIfPrivateIntended:
		return yaml.Marshal("enforced_if_private_intended")
	case EnforcedIfAutoNSFW:
		return yaml.Marshal("enforced_if_auto_nsfw")
	case EnforcedIfInteracting:
		return yaml.Marshal("enforced_if_interacting")
	case EnforcedIfConsuming:
		return yaml.Marshal("enforced_if_consuming")
	case AlwaysEnforced:
		return yaml.Marshal("always_enforced")
	}
	return []byte{}, errors.New(fmt.Sprintf("Invalid EnforcementUrgency value: %v", u))
}

func (u *EnforcementUrgency) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch strings.ToLower(value) {
	case "optionally_enforced":
		*u = OptionallyEnforced
	case "enforced_if_mass_traffic":
		*u = EnforcedIfMassTraffic
	case "enforced_if_private_intended":
		*u = EnforcedIfPrivateIntended
	case "enforced_if_auto_nsfw":
		*u = EnforcedIfAutoNSFW
	case "enforced_if_interacting":
		*u = EnforcedIfInteracting
	case "enforced_if_consuming":
		*u = EnforcedIfConsuming
	case "always_enforced":
		*u = AlwaysEnforced
	}
	return nil
}

func (u *EnforcementUrgency) UnmarshalYAML(data *yaml.Node) error {
	var value string
	if err := data.Decode(&value); err != nil {
		return err
	}

	switch strings.ToLower(value) {
	case "optionally_enforced":
		*u = OptionallyEnforced
	case "enforced_if_mass_traffic":
		*u = EnforcedIfMassTraffic
	case "enforced_if_private_intended":
		*u = EnforcedIfPrivateIntended
	case "enforced_if_auto_nsfw":
		*u = EnforcedIfAutoNSFW
	case "enforced_if_interacting":
		*u = EnforcedIfInteracting
	case "enforced_if_consuming":
		*u = EnforcedIfConsuming
	case "always_enforced":
		*u = AlwaysEnforced
	}
	return nil
}

type AccountPolicy struct {
	// Is the content only accessible to users logged in to an account?
	LoggedIn *EnforcementUrgency `json:"logged_in" yaml:"logged_in"`

	// Is the user account required to be registered from the correct regional variant of the digital platform; accounts registered from another region of the same service brand are not allowed?
	MatchRegion *EnforcementUrgency `json:"match_region" yaml:"match_region"`

	// Is the user account must not be banned by the author who publishes the specified content?
	NotBannedByAuthor *EnforcementUrgency `json:"not_banned_by_author" yaml:"not_banned_by_author"`

	// Is the user account must not be banned by the whole platform?
	NotBannedByPlatform *EnforcementUrgency `json:"not_banned_by_platform" yaml:"not_banned_by_platform"`
}

type ClientAppPolicy struct {
	// Does the content need to be viewed on a desktop computer device with a specified application software?
	DesktopApp *EnforcementUrgency `json:"desktop_app" yaml:"desktop_app"`

	// Does the content need to be viewed on a mobile device with a specified mobile application software?
	MobileApp *EnforcementUrgency `json:"mobile_app" yaml:"mobile_app"`

	// Does the content need to be viewed on a dedicated hardware, such as a Smart TV?
	DedicatedHardware *EnforcementUrgency `json:"dedicated_hardware" yaml:"dedicated_hardware"`

	// Does the content need to be viewed from the latest version of the client software either supplied or endorsed by the platform owner?
	MustLatestClientSoftware *EnforcementUrgency `json:"must_latest_client_software" yaml:"must_latest_client_software"`

	// Does the content need to be viewed from the latest version of the host software (e.g., desktop operating system or the IoT hardware)?
	MustLatestHostSoftware *EnforcementUrgency `json:"must_latest_host_software" yaml:"must_latest_host_software"`

	// Does the content need to be viewed from the correct regional variant of the client software/hardware. If this policy is enforced for the whole account, use [AccountPolicy.MatchRegion] flag instead?
	MatchRegion *EnforcementUrgency `json:"match_region" yaml:"match_region"`
}

type ContentAccessPolicy struct {
	// Does any user need to subscribe to the content author, free or paid, in order to view the content?
	AuthorSubscription *EnforcementUrgency `json:"author_subscription" yaml:"author_subscription"`

	// Is the user account must subscribe to a specified pricing tier to view the content?
	PlatformSubscription *EnforcementUrgency `json:"platform_subscription" yaml:"platform_subscription"`

	// Does any user need to pay to access the content?
	ContentPayment *EnforcementUrgency `json:"content_payment" yaml:"content_payment"`

	// Does any user need to have a content access code, password, passcode, PIN?
	ContentPassword *EnforcementUrgency `json:"content_password" yaml:"content_password"`
}

type PlatformLinks struct {
	// Does the site have a link to create a new account?
	CreateAccount *string `json:"create_account" yaml:"create_account"`

	// Does the site have a link to sign in to an existing account?
	SignIn *string `json:"sign_in" yaml:"sign_in"`

	// Does the site have a link for account recovery ("forgot password?")?
	ForgotPassword *string `json:"forgot_password" yaml:"forgot_password"`

	// Does the site have a link to the End User License Agreement, Terms of Use, or Terms of Service?
	TermsOfUse *string `json:"terms_of_use" yaml:"terms_of_use"`

	// Does the site have a link to the Cookies Policy?
	Cookies *string `json:"cookies" yaml:"cookies"`

	// Does the site have a link to the Privacy Policy?
	Privacy *string `json:"privacy" yaml:"privacy"`

	// Does the site have a link to the Community Guidelines or Moderation Guidelines?
	Moderation *string `json:"moderation" yaml:"moderation"`

	// Does the site have a link to the License terms for content?
	License *string `json:"license" yaml:"license"`
}

type PolicySummary struct {
	// The name of the website or platform
	Name string `json:"name" yaml:"name"`

	// The supported regions of the website, platform, or content
	Regions []string `json:"regions" yaml:"regions"`

	// Does the site have any known aliases?
	DomainAliases []string `json:"domain_aliases" yaml:"domain_aliases"`

	// Does the site require users to complete
	Captcha *EnforcementUrgency `json:"captcha" yaml:"captcha"`

	// Does the site require to be accessed on a specified client/viewer application software of hardware?
	ClientApp *ClientAppPolicy `json:"client_app" yaml:"client_app"`

	// Does the site require to be accessed from a specified account?
	Account *AccountPolicy `json:"account" yaml:"account"`

	// Does the content have specific access policy?
	ContentAccess *ContentAccessPolicy `json:"content_access" yaml:"content_access"`

	// Does the site, platform, or content contain official regulations?
	Links *PlatformLinks `json:"links" yaml:"links"`
}
