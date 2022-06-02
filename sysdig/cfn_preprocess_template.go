/*
This file contains definition for the following functions.

terraformPreModifications is used to modify the container definitions passed to the cfn patcher such that it modifies casing issues in any ECS json content so that it can be processed by the kilt patcher.

GetValueFromTemplate is used to obtain key, value information from JSON object
*/

package sysdig

import (
	"context"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/rs/zerolog/log"
)

// GetValueFromTemplate can be used to obtain string value from JSON object
func GetValueFromTemplate(what *gabs.Container) (string, *gabs.Container) {
	var result string
	var fallback *gabs.Container

	switch v := what.Data().(type) {
	case string:
		result = v
		fallback = nil
	default:
		fallback = what
		result = "placeholder: " + what.String()
	}
	return result, fallback
}

func terraformPreModifications(ctx context.Context, patchedStack []byte) ([]byte, error) {
	l := log.Ctx(ctx)
	template, err := gabs.ParseJSON(patchedStack)
	if err != nil {
		l.Error().Err(err).Msg("failed to parse input fragment")
		return nil, err
	}

	// This code block is used when parsing ECS JSON format
	for _, resource := range template.S("Resources").ChildrenMap() {
		for _, container := range resource.S("Properties", "ContainerDefinitions").Children() {
			if container.Exists("image") {
				passthrough, _ := GetValueFromTemplate(container.S("image"))
				_, err = container.Set(passthrough, "Image")
				if err != nil {
					return nil, fmt.Errorf("Could not update Image field: %v", err)
				}

				err = container.Delete("image")
				if err != nil {
					return nil, fmt.Errorf("could not delete image in the Container definition: %w", err)
				}
			}

			if container.Exists("Environment") {
				for _, env := range container.S("Environment").Children() {
					if env.Exists("name") {
						name, _ := env.S("name").Data().(string)
						err = env.Delete("name")
						if err != nil {
							return nil, fmt.Errorf("Could not delete \"name\" key in Environment: %v", err)
						}
						_, err = env.Set(name, "Name")
						if err != nil {
							return nil, fmt.Errorf("Could not assign \"Name\" key in Environment: %v", err)
						}
					}
					if env.Exists("value") {
						value, _ := env.S("value").Data().(string)
						err = env.Delete("value")
						if err != nil {
							return nil, fmt.Errorf("Could not delete \"value\" key in Environment: %v", err)
						}
						_, err = env.Set(value, "Value")
						if err != nil {
							return nil, fmt.Errorf("Could not assign \"Value\" key in Environment: %v", err)
						}
					}
				}
			}

			if container.Exists("entryPoint") {
				for _, arg := range container.S("entryPoint").Children() {
					passthrough, _ := GetValueFromTemplate(arg)
					err = container.ArrayAppend(passthrough, "EntryPoint")
					if err != nil {
						return nil, fmt.Errorf("Could not append entrypoint values to EntryPoint: %v", err)
					}
				}

				err = container.Delete("entryPoint")
				if err != nil {
					return nil, fmt.Errorf("could not delete entryPoint in the Container definition: %w", err)
				}
			}
		}
	}

	return template.Bytes(), err
}
