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
func updateKeys(inputMap gabs.Container) error {
	// in this case, the object is probably an array, so update each child
	if len(inputMap.ChildrenMap()) == 0 {
		for _, child := range inputMap.Children() {
			err := updateKeys(*child)
			if err != nil {
				return err
			}
		}
	} else {
		for key, child := range inputMap.ChildrenMap() {
			_, err := inputMap.Set(child.Data(), capitalize(key))
			if err != nil {
				return fmt.Errorf("failed to update new key %s", capitalize(key))
			}

			err = inputMap.Delete(key)
			if err != nil {
				return fmt.Errorf("failed to delete old key %s" + key)
			}

			// recurse to update child's keys
			err = updateKeys(*child)
			if err != nil {
				return err
			}
		}
	}

	return nil
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

			if container.Exists("name") {
				passthrough, _ := GetValueFromTemplate(container.S("name"))
				_, err = container.Set(passthrough, "Name")
				if err != nil {
					return nil, fmt.Errorf("Could not update Name field: %v", err)
				}

				err = container.Delete("name")
				if err != nil {
					return nil, fmt.Errorf("could not delete name in the Container definition: %w", err)
				}
			}

			if container.Exists("environment") {
				for _, env := range container.S("environment").Children() {
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
				err = updateKeys(*container.S("volumesFrom"))
				if err != nil {
					return nil, fmt.Errorf("could not update volumesFrom items: %v", err)
				}

				passthrough, _ := GetValueFromTemplate(container.S("volumesFrom"))
				parsedPassthrough, _ := gabs.ParseJSON([]byte(passthrough))
				_, err = container.Set(parsedPassthrough, "VolumesFrom")
				if err != nil {
					return nil, fmt.Errorf("could not update VolumesFrom field: %v", err)
				}

				err = container.Delete("volumesFrom")
				if err != nil {
					return nil, fmt.Errorf("could not delete volumesFrom in the container definition: %w", err)
				}
			}

			if container.Exists("linuxParameters") {
				err = updateKeys(*container.S("linuxParameters"))
				if err != nil {
					return nil, fmt.Errorf("could not update linuxParameters items: %v", err)
				}

				passthrough, _ := GetValueFromTemplate(container.S("linuxParameters"))
				parsedPassthrough, _ := gabs.ParseJSON([]byte(passthrough))
				_, err = container.Set(parsedPassthrough, "LinuxParameters")
				if err != nil {
					return nil, fmt.Errorf("could not update LinuxParameters field: %v", err)
				}

				err = container.Delete("linuxParameters")
				if err != nil {
					return nil, fmt.Errorf("could not delete linuxParameters in the COntainer definition: %w", err)
				}
			}
		}
	}

	return template.Bytes(), err
}
