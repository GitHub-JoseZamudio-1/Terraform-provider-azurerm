package cdn

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func ConvertFrontdoorTags(tagMap *map[string]string) map[string]*string {
	t := make(map[string]*string)

	if tagMap != nil {
		for k, v := range *tagMap {
			tagKey := k
			tagValue := v
			t[tagKey] = &tagValue
		}
	}

	return t
}

func ConvertBoolToEnabledState(isEnabled bool) track1.EnabledState {
	out := track1.EnabledState(track1.EnabledStateDisabled)

	if isEnabled {
		out = track1.EnabledState(track1.EnabledStateEnabled)
	}

	return out
}

func ConvertEnabledStateToBool(enabledState *track1.EnabledState) bool {
	if enabledState == nil {
		return false
	}

	return (*enabledState == track1.EnabledState(track1.EnabledStateEnabled))
}

func expandResourceReference(input string) *track1.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &track1.ResourceReference{
		ID: utils.String(input),
	}
}

func flattenResourceReference(input *track1.ResourceReference) string {
	result := ""
	if input == nil {
		return result
	}

	if input.ID != nil {
		result = *input.ID
	}

	return result
}

// func ConvertBoolToOriginsEnabledState(isEnabled bool) *afdorigins.EnabledState {
// 	out := afdorigins.EnabledState(afdorigins.EnabledStateDisabled)

// 	if isEnabled {
// 		out = afdorigins.EnabledState(afdorigins.EnabledStateEnabled)
// 	}

// 	return &out
// }

// func ConvertOriginsEnabledStateToBool(enabledState *afdorigins.EnabledState) bool {
// 	if enabledState == nil {
// 		return false
// 	}

// 	return (*enabledState == afdorigins.EnabledState(afdorigins.EnabledStateEnabled))
// }

func ConvertBoolToRouteHttpsRedirect(isEnabled bool) track1.HTTPSRedirect {
	out := track1.HTTPSRedirect(track1.HTTPSRedirectDisabled)

	if isEnabled {
		out = track1.HTTPSRedirect(track1.HTTPSRedirectEnabled)
	}

	return out
}

func ConvertRouteHttpsRedirectToBool(httpsRedirect *track1.HTTPSRedirect) bool {
	if httpsRedirect == nil {
		return false
	}

	return (*httpsRedirect == track1.HTTPSRedirect(track1.HTTPSRedirectEnabled))
}

func ConvertBoolToRouteLinkToDefaultDomain(isLinked bool) track1.LinkToDefaultDomain {
	out := track1.LinkToDefaultDomain(track1.LinkToDefaultDomainDisabled)

	if isLinked {
		out = track1.LinkToDefaultDomain(track1.LinkToDefaultDomainEnabled)
	}

	return out
}

func ConvertRouteLinkToDefaultDomainToBool(linkToDefaultDomain *track1.LinkToDefaultDomain) bool {
	if linkToDefaultDomain == nil {
		return false
	}

	return (*linkToDefaultDomain == track1.LinkToDefaultDomain(track1.LinkToDefaultDomainEnabled))
}

func IsValidDomain(i interface{}, k string) (warnings []string, errors []error) {
	if warn, err := validation.IsIPv6Address(i, k); len(err) == 0 {
		return warn, err
	}

	if warn, err := validation.IsIPv4Address(i, k); len(err) == 0 {
		return warn, err
	}

	// TODO: Figure out a better way to validate Doman Name if not and IP Address
	if warn, err := validation.StringIsNotEmpty(i, k); len(err) == 0 {
		return warn, err
	}

	return warnings, errors
}

func ValidateFrontdoorRuleSetName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, regexErrs := validate.RegExHelper(i, k, `(^[a-zA-Z])([\da-zA-Z]{1,88})([a-zA-Z]$)`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 1 and 90 characters in length and begin with a letter, end with a letter and may contain only letters and numbers, got %q`, v, k))
	}

	return nil, nil
}

func ValidateFrontdoorCacheDuration(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, regexErrs := validate.RegExHelper(i, k, `^([0-3]|([1-9][0-9])|([1-3][0-6][0-5])).((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between in the d.HH:MM:SS format and must be equal to or lower than %q, got %q`, v, "365.23:59:59", k))
	}

	return nil, nil
}

func ValidateContentTypes(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	// Per the IANA no whitespace is allowed in a Content Type
	if strings.Contains(v, " ") {
		return nil, append(errors, fmt.Errorf(`%q must not contain any whitespace, got %q`, k, v))
	}

	if strings.Contains(v, ";") {
		// Content Type has a parameter, error out
		return nil, append(errors, fmt.Errorf(`%q is not valid, Content Types with parameters are not allowed, got %q`, k, v))
	}

	if m, regexErrs := validate.RegExHelper(i, k, `^(application|audio|font|image|message|model|multipart|text|video)\/[-\w]+(\.[-\w]+)*([+][-\w]+)?$`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be a valid Content Type and a subtype concatenated with a slash(e.g. text/html), got %q`, k, v))
	}

	return nil, nil
}

func SchemaFrontdoorOperator() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			"Any",
			"Equal",
			"Contains",
			"BeginsWith",
			"EndsWith",
			"LessThan",
			"LessThanOrEqual",
			"GreaterThan",
			"GreaterThanOrEqual",
			"RegEx",
		}, false),
	}
}

func SchemaFrontdoorNegateCondition() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
}

func SchemaFrontdoorMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 10,

		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func SchemaFrontdoorRuleTransforms() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 6,

		Elem: &pluginsdk.Schema{
			Type:    pluginsdk.TypeString,
			Default: string(track1.TransformLowercase),
			ValidateFunc: validation.StringInSlice([]string{
				string(track1.TransformLowercase),
				string(track1.TransformRemoveNulls),
				string(track1.TransformTrim),
				string(track1.TransformUppercase),
				string(track1.TransformURLDecode),
				string(track1.TransformURLEncode),
			}, false),
		},
	}
}
