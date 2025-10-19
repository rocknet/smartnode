package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"github.com/rocket-pool/smartnode/addons/obol"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/types/addons"
	cfgtypes "github.com/rocket-pool/smartnode/shared/types/config"
)

// The page wrapper for the Obol addon config
type AddonObolPage struct {
	addonsPage     *AddonsPage
	page           *page
	layout         *standardLayout
	masterConfig   *config.RocketPoolConfig
	addon          addons.SmartnodeAddon
	enabledBox     *parameterizedFormItem
	clusterRoleBox *parameterizedFormItem
	otherParams    []*parameterizedFormItem
}

// Creates a new page for the Obol addon settings
func NewAddonObolPage(addonsPage *AddonsPage, addon addons.SmartnodeAddon) *AddonObolPage {
	configPage := &AddonObolPage{
		addonsPage:   addonsPage,
		masterConfig: addonsPage.home.md.Config,
		addon:        addon,
	}

	configPage.createContent()

	configPage.page = newPage(
		addonsPage.page,
		"settings-addon-obol",
		addon.GetName(),
		addon.GetDescription(),
		configPage.layout.grid,
	)

	return configPage
}

// Get the underlying page
func (configPage *AddonObolPage) getPage() *page {
	return configPage.page
}

// Creates the content for the Obol settings page
func (configPage *AddonObolPage) createContent() {
	// Create the layout
	configPage.layout = newStandardLayout()
	configPage.layout.createForm(&configPage.masterConfig.Smartnode.Network, fmt.Sprintf("%s Settings", configPage.addon.GetName()))
	configPage.layout.setupEscapeReturnHomeHandler(configPage.addonsPage.home.md, configPage.addonsPage.page)

	// Get the config
	obolConfig := configPage.addon.GetConfig().(*obol.ObolConfig)

	// Get the parameters - separate enabled and cluster role from the rest
	enabledParam := configPage.addon.GetEnabledParameter()
	clusterRoleParam := &obolConfig.ClusterRole
	otherParams := []*cfgtypes.Parameter{}

	for _, param := range configPage.addon.GetConfig().GetParameters() {
		if param.ID != enabledParam.ID && param.ID != clusterRoleParam.ID {
			otherParams = append(otherParams, param)
		}
	}

	// Set up the form items
	configPage.enabledBox = createParameterizedCheckbox(enabledParam)
	configPage.clusterRoleBox = createParameterizedDropDown(clusterRoleParam, configPage.layout.descriptionBox)
	configPage.otherParams = createParameterizedFormItems(otherParams, configPage.layout.descriptionBox)

	// Map the parameters to the form items in the layout
	configPage.layout.mapParameterizedFormItems(configPage.enabledBox, configPage.clusterRoleBox)
	configPage.layout.mapParameterizedFormItems(configPage.otherParams...)

	// Set up the setting callbacks
	configPage.enabledBox.item.(*tview.Checkbox).SetChangedFunc(func(checked bool) {
		if enabledParam.Value == checked {
			return
		}
		enabledParam.Value = checked
		configPage.handleEnableChanged()
	})

	configPage.clusterRoleBox.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
		if clusterRoleParam.Value == clusterRoleParam.Options[index].Value {
			return
		}
		clusterRoleParam.Value = clusterRoleParam.Options[index].Value
		configPage.handleEnableChanged()
	})

	// Set up callback for NumOperators to trigger refresh when changed
	numOperatorsParam := &obolConfig.NumOperators
	for _, paramItem := range configPage.otherParams {
		if paramItem.parameter.ID == "numOperators" {
			paramItem.item.(*DropDown).SetSelectedFunc(func(text string, index int) {
				if numOperatorsParam.Value == numOperatorsParam.Options[index].Value {
					return
				}
				numOperatorsParam.Value = numOperatorsParam.Options[index].Value
				configPage.handleEnableChanged()
			})
			break
		}
	}

	// Do the initial draw
	configPage.handleEnableChanged()
}

// Handle enable/disable toggle
func (configPage *AddonObolPage) handleEnableChanged() {
	configPage.layout.form.Clear(true)
	configPage.layout.form.AddFormItem(configPage.enabledBox.item)

	// Only add the supporting stuff if the addon is enabled
	if configPage.addon.GetEnabledParameter().Value == false {
		return
	}

	// Add cluster role selector
	configPage.layout.form.AddFormItem(configPage.clusterRoleBox.item)

	// Get the config to check cluster role and operator count
	obolConfig := configPage.addon.GetConfig().(*obol.ObolConfig)
	role := obolConfig.ClusterRole.Value.(obol.ClusterRole)
	numOperatorsStr := obolConfig.NumOperators.Value.(string)

	// Convert numOperators string to int
	numOperators := 3 // default
	if n, err := strconv.Atoi(numOperatorsStr); err == nil {
		numOperators = n
	}

	// Filter parameters based on role
	paramsToShow := []*parameterizedFormItem{}

	// Parameters shown for both roles
	commonParams := map[string]bool{
		"containerTag":         true,
		"p2pPort":              true,
		"clusterDefinitionURL": true,
	}

	for _, param := range configPage.otherParams {
		// Common parameters are shown for both roles
		if commonParams[param.parameter.ID] {
			paramsToShow = append(paramsToShow, param)
			continue
		}

		// All other parameters are only shown for Creator role
		if role == obol.ClusterRole_Creator {
			// For ENR fields, only show the number needed (numOperators - 1)
			if strings.HasPrefix(param.parameter.ID, "operatorENR") {
				// Extract the ENR number from the ID (operatorENR1 -> 1, operatorENR2 -> 2, etc.)
				enrNum := int(param.parameter.ID[11] - '0') // Get the digit from "operatorENRX"

				// Only show ENRs up to (numOperators - 1)
				// e.g., if 4 operators, show ENR 1, 2, 3 (the 4th is the creator)
				if enrNum <= numOperators-1 {
					paramsToShow = append(paramsToShow, param)
				}
			} else {
				// Show all non-ENR creator parameters
				paramsToShow = append(paramsToShow, param)
			}
		}
	}

	// Add the filtered parameters
	configPage.layout.addFormItems(paramsToShow)
	configPage.layout.refresh()
}

// Handle a bulk redraw request
func (configPage *AddonObolPage) handleLayoutChanged() {
	configPage.handleEnableChanged()
}
