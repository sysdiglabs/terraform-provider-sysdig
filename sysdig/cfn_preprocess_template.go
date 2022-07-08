/*
This file contains definition for the following functions.

terraformPreModifications is used to modify the container definitions passed to the cfn patcher such that it modifies casing issues in any ECS json content so that it can be processed by the kilt patcher.

GetValueFromTemplate is used to obtain key, value information from JSON object
*/

package sysdig

import (
	"context"
	"fmt"
	"unicode"

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
		result = what.String()
	}
	return result, fallback
}

func capitalize(key string) string {
	r := []rune(key)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// updateKey recursively capitalizes the first letter of each key in the input object
func updateKeys(inputMap gabs.Container) {
	// in this case, the object is probably an array, so update each child
	if len(inputMap.ChildrenMap()) == 0 {
		for _, child := range inputMap.Children() {
			updateKeys(*child)
		}
	} else {
		for key, child := range inputMap.ChildrenMap() {
			_, err := inputMap.Set(child.Data(), capitalize(key))
			if err != nil {
				log.Error().Err(err).Msg("Failed to update key " + key)
			}

			err = inputMap.Delete(key)
			if err != nil {
				log.Error().Err(err).Msg("Failed to update key " + key)
			}

			// recurse to update child's keys
			updateKeys(*child)
		}
	}
}

func terraformPreModifications(ctx context.Context, patchedStack []byte) ([]byte, error) {
	l := log.Ctx(ctx)
	template, err := gabs.ParseJSON(patchedStack)
	if err != nil {
		l.Error().Err(err).Msg("failed to parse input fragment")
		return nil, err
	}

	l.Debug().Msg("starting terraformPreModifications")

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

			if container.Exists("environment") {
				updateKeys(*container.S("environment"))
				passthrough, _ := GetValueFromTemplate(container.S("environment"))
				parsedPassthrough, _ := gabs.ParseJSON([]byte(passthrough))
				_, err = container.Set(parsedPassthrough, "Environment")
				if err != nil {
					return nil, fmt.Errorf("Could not update Environment field: %v", err)
				}

				err = container.Delete("environment")
				if err != nil {
					return nil, fmt.Errorf("could not delete environment in the Container definition: %w", err)
				}
			}

			if container.Exists("entryPoint") {
				passthrough, _ := GetValueFromTemplate(container.S("entryPoint"))
				parsedPassthrough, _ := gabs.ParseJSON([]byte(passthrough))
				_, err = container.Set(parsedPassthrough, "EntryPoint")
				if err != nil {
					return nil, fmt.Errorf("Could not update EntryPoint field: %v", err)
				}

				err = container.Delete("entryPoint")
				if err != nil {
					return nil, fmt.Errorf("could not delete entryPoint in the Container definition: %w", err)
				}
			}

			if container.Exists("command") {
				passthrough, _ := GetValueFromTemplate(container.S("command"))
				parsedPassthrough, _ := gabs.ParseJSON([]byte(passthrough))
				_, err = container.Set(parsedPassthrough, "Command")
				if err != nil {
					return nil, fmt.Errorf("Could not update Command field: %v", err)
				}

				err = container.Delete("command")
				if err != nil {
					return nil, fmt.Errorf("could not delete command in the Container definition: %w", err)
				}
			}

			if container.Exists("volumesFrom") {
				updateKeys(*container.S("volumesFrom"))
				passthrough, _ := GetValueFromTemplate(container.S("volumesFrom"))
				parsedPassthrough, _ := gabs.ParseJSON([]byte(passthrough))
				_, err = container.Set(parsedPassthrough, "VolumesFrom")
				if err != nil {
					return nil, fmt.Errorf("Could not update VolumesFrom field: %v", err)
				}

				err = container.Delete("volumesFrom")
				if err != nil {
					return nil, fmt.Errorf("could not delete volumesFrom in the Container definition: %w", err)
				}
			}

			if container.Exists("linuxParameters") {
				updateKeys(*container.S("linuxParameters"))
				passthrough, _ := GetValueFromTemplate(container.S("linuxParameters"))
				parsedPassthrough, _ := gabs.ParseJSON([]byte(passthrough))
				_, err = container.Set(parsedPassthrough, "LinuxParameters")
				if err != nil {
					return nil, fmt.Errorf("Could not update LinuxParameters field: %v", err)
				}

				err = container.Delete("linuxParameters")
				if err != nil {
					return nil, fmt.Errorf("could not delete linuxParameters in the Container definition: %w", err)
				}
			}
		}
	}

	return template.Bytes(), err
}
