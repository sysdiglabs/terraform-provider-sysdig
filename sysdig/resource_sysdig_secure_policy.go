package sysdig

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var defaultMatchActions = map[string]string{
	"accept": "DEFAULT_MATCH_EFFECT_ACCEPT",
	"deny":   "DEFAULT_MATCH_EFFECT_DENY",
	"none":   "DEFAULT_MATCH_EFFECT_NEXT",
}

var matchActions = map[string]string{
	"accept": "MATCH_EFFECT_ACCEPT",
	"deny":   "MATCH_EFFECT_DENY",
	"none":   "MATCH_EFFECT_NEXT",
}

func resourceSysdigSecurePolicy() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigPolicyCreate,
		Read:   resourceSysdigPolicyRead,
		Update: resourceSysdigPolicyUpdate,
		Delete: resourceSysdigPolicyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"severity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"container_scope": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"host_scope": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"falco_rule_name_regex": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notification_channels": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"stop", "pause"}, false),
						},
						"capture": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"seconds_after_event": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"seconds_before_event": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"processes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "deny", "none"}, false),
						},
						"whitelist": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"blacklist": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"containers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "deny", "none"}, false),
						},
						"whitelist": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"blacklist": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"syscalls": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "deny", "none"}, false),
						},
						"whitelist": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"blacklist": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"network": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inbound": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "deny", "none"}, false),
						},
						"outbound": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "deny", "none"}, false),
						},
						"listening_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"accept", "deny", "none"}, false),
									},
									"tcp": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"whitelist": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"blacklist": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
											},
										},
									},
									"udp": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"whitelist": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"blacklist": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"filesystem": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"whitelist": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"blacklist": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"readwrite": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"whitelist": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"blacklist": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"other_paths": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"accept", "deny", "none"}, true),
						},
					},
				},
			},
		},
	}
}

func resourceSysdigPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SysdigSecureClient)

	policy := policyFromResourceData(d)
	policy, err := client.CreatePolicy(policy)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(policy.ID))
	d.Set("version", policy.Version)

	return nil
}

func policyFromResourceData(d *schema.ResourceData) Policy {
	policy := Policy{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Severity:       d.Get("severity").(int),
		ContainerScope: d.Get("container_scope").(bool),
		HostScope:      d.Get("host_scope").(bool),
		Enabled:        d.Get("enabled").(bool),
		Scope:          d.Get("filter").(string),
	}

	addActionsToPolicy(d, &policy)
	addProcessesToPolicy(d, &policy)
	addContainersToPolicy(d, &policy)
	addSyscallsToPolicy(d, &policy)
	addNetworkToPolicy(d, &policy)
	addFilesystemToPolicy(d, &policy)

	if falcoRuleNameRegex := d.Get("falco_rule_name_regex"); falcoRuleNameRegex != nil {
		policy.FalcoConfiguration = FalcoConfiguration{
			RuleNameRegEx: falcoRuleNameRegex.(string),
		}
	}

	if notificationChannelIdSet, ok := d.Get("notification_channels").(*schema.Set); ok {
		list := notificationChannelIdSet.List()
		for _, id := range list {
			if idStr, ok := id.(string); ok {
				idStr = strings.TrimSpace(idStr)
				idInt, _ := strconv.Atoi(idStr)
				policy.NotificationChannelIds = append(policy.NotificationChannelIds, idInt)
			}
		}
	}

	return policy
}
func addFilesystemToPolicy(d *schema.ResourceData, policy *Policy) {
	filesystem := d.Get("filesystem").([]interface{})
	if len(filesystem) == 0 {
		return
	}

	policy.FileSystemConfiguration = &FileSystemConfiguration{
		Fields: []string{
			"fd.name",
			"proc.cmdline",
			"proc.name",
		},
	}

	if otherFilesFilesystemAction := d.Get("filesystem.0.other_paths").(string); otherFilesFilesystemAction != "" {
		policy.FileSystemConfiguration.OnDefault = defaultMatchActions[otherFilesFilesystemAction]
	} else {
		policy.FileSystemConfiguration.OnDefault = defaultMatchActions["none"]
	}

	setToSlice := func(set *schema.Set) (result []string) {
		for _, element := range set.List() {
			result = append(result, element.(string))
		}
		return
	}

	if readWhitelist := d.Get("filesystem.0.read.0.whitelist").(*schema.Set); readWhitelist.Len() > 0 {
		policy.FileSystemConfiguration.List = append(policy.FileSystemConfiguration.List, FileSystemConfigurationListElement{
			Values:     setToSlice(readWhitelist),
			AccessType: "ACCESS_READ",
			OnMatch:    matchActions["accept"],
		})
	}

	if readBlacklist := d.Get("filesystem.0.read.0.blacklist").(*schema.Set); readBlacklist.Len() > 0 {
		policy.FileSystemConfiguration.List = append(policy.FileSystemConfiguration.List, FileSystemConfigurationListElement{
			Values:     setToSlice(readBlacklist),
			AccessType: "ACCESS_READ",
			OnMatch:    matchActions["deny"],
		})
	}

	if readwriteWhitelist := d.Get("filesystem.0.readwrite.0.whitelist").(*schema.Set); readwriteWhitelist.Len() > 0 {
		policy.FileSystemConfiguration.List = append(policy.FileSystemConfiguration.List, FileSystemConfigurationListElement{
			Values:     setToSlice(readwriteWhitelist),
			AccessType: "ACCESS_READWRITE",
			OnMatch:    matchActions["accept"],
		})
	}

	if readwriteBlacklist := d.Get("filesystem.0.readwrite.0.blacklist").(*schema.Set); readwriteBlacklist.Len() > 0 {
		policy.FileSystemConfiguration.List = append(policy.FileSystemConfiguration.List, FileSystemConfigurationListElement{
			Values:     setToSlice(readwriteBlacklist),
			AccessType: "ACCESS_READWRITE",
			OnMatch:    matchActions["deny"],
		})
	}
}

func addProcessesToPolicy(d *schema.ResourceData, policy *Policy) {
	processes := d.Get("processes").([]interface{})
	if len(processes) == 0 {
		return
	}

	policy.ProcessesConfiguration = &ProcessesConfiguration{
		Fields: []string{
			"proc.name",
			"proc.cmdline",
		},
	}

	if defaultProcessesAction := d.Get("processes.0.default").(string); defaultProcessesAction != "" {
		policy.ProcessesConfiguration.OnDefault = defaultMatchActions[defaultProcessesAction]
	} else {
		policy.ProcessesConfiguration.OnDefault = defaultMatchActions["none"]
	}

	if whitelist := d.Get("processes.0.whitelist").(*schema.Set); whitelist.Len() > 0 {
		element := ProcessesConfigurationElement{
			OnMatch: matchActions["accept"],
		}
		for _, value := range whitelist.List() {
			element.Values = append(element.Values, value.(string))
		}
		policy.ProcessesConfiguration.List = append(policy.ProcessesConfiguration.List, element)
	}
	if blacklist := d.Get("processes.0.blacklist").(*schema.Set); blacklist.Len() > 0 {
		element := ProcessesConfigurationElement{
			OnMatch: matchActions["deny"],
		}
		for _, value := range blacklist.List() {
			element.Values = append(element.Values, value.(string))
		}
		policy.ProcessesConfiguration.List = append(policy.ProcessesConfiguration.List, element)
	}
}

func addContainersToPolicy(d *schema.ResourceData, policy *Policy) {
	containers := d.Get("containers").([]interface{})
	if len(containers) == 0 {
		return
	}

	policy.ContainerImagesConfiguration = &ContainerImagesConfiguration{
		Fields: []string{
			"container.id",
			"container.name",
			"container.image.id",
			"container.image",
		},
	}

	if defaultContainersAction := d.Get("containers.0.default").(string); defaultContainersAction != "" {
		policy.ContainerImagesConfiguration.OnDefault = defaultMatchActions[defaultContainersAction]
	} else {
		policy.ContainerImagesConfiguration.OnDefault = defaultMatchActions["none"]
	}

	if whitelist := d.Get("containers.0.whitelist").(*schema.Set); whitelist.Len() > 0 {
		element := ContainerImagesConfigurationElement{
			OnMatch: matchActions["accept"],
		}
		for _, value := range whitelist.List() {
			element.Values = append(element.Values, value.(string))
		}
		policy.ContainerImagesConfiguration.List = append(policy.ContainerImagesConfiguration.List, element)
	}
	if blacklist := d.Get("containers.0.blacklist").(*schema.Set); blacklist.Len() > 0 {
		element := ContainerImagesConfigurationElement{
			OnMatch: matchActions["deny"],
		}
		for _, value := range blacklist.List() {
			element.Values = append(element.Values, value.(string))
		}
		policy.ContainerImagesConfiguration.List = append(policy.ContainerImagesConfiguration.List, element)
	}
}

func addSyscallsToPolicy(d *schema.ResourceData, policy *Policy) {
	syscalls := d.Get("syscalls").([]interface{})
	if len(syscalls) == 0 {
		return
	}

	policy.SyscallsConfiguration = &SyscallsConfiguration{
		Fields: []string{
			"evt.type",
			"proc.cmdline",
			"proc.name",
		},
	}

	if defaultSyscallsAction := d.Get("syscalls.0.default").(string); defaultSyscallsAction != "" {
		policy.SyscallsConfiguration.OnDefault = defaultMatchActions[defaultSyscallsAction]
	} else {
		policy.SyscallsConfiguration.OnDefault = defaultMatchActions["none"]
	}

	if whitelist := d.Get("syscalls.0.whitelist").(*schema.Set); whitelist.Len() > 0 {
		element := SyscallsConfigurationElement{
			OnMatch: matchActions["accept"],
		}
		for _, value := range whitelist.List() {
			element.Values = append(element.Values, value.(string))
		}
		policy.SyscallsConfiguration.List = append(policy.SyscallsConfiguration.List, element)
	}
	if blacklist := d.Get("syscalls.0.blacklist").(*schema.Set); blacklist.Len() > 0 {
		element := SyscallsConfigurationElement{
			OnMatch: matchActions["deny"],
		}
		for _, value := range blacklist.List() {
			element.Values = append(element.Values, value.(string))
		}
		policy.SyscallsConfiguration.List = append(policy.SyscallsConfiguration.List, element)
	}
}

func addActionsToPolicy(d *schema.ResourceData, policy *Policy) {
	actions := d.Get("actions").([]interface{})
	if len(actions) == 0 {
		return
	}

	policy.Actions = []Action{}

	containerAction := d.Get("actions.0.container").(string)
	if containerAction != "" {
		containerAction = strings.ToUpper("POLICY_ACTION_" + containerAction)

		policy.Actions = append(policy.Actions, Action{Type: containerAction})
	}

	captureAction := d.Get("actions.0.capture").(map[string]interface{})
	if len(captureAction) != 0 {
		afterEventNs, _ := strconv.Atoi(d.Get("actions.0.capture.seconds_after_event").(string) + "000000000")
		beforeEventNs, _ := strconv.Atoi(d.Get("actions.0.capture.seconds_before_event").(string) + "000000000")
		policy.Actions = append(policy.Actions, Action{
			Type:                 "POLICY_ACTION_CAPTURE",
			IsLimitedToContainer: false,
			AfterEventNs:         afterEventNs,
			BeforeEventNs:        beforeEventNs,
		})
	}
}

func addNetworkToPolicy(d *schema.ResourceData, policy *Policy) {
	network := d.Get("network").([]interface{})
	if len(network) == 0 {
		return
	}

	policy.NetworkConfiguration = &NetworkConfiguration{}
	policy.NetworkConfiguration.Inbound.Fields = []string{
		"fd.l4proto",
		"fd.sip",
		"fd.sport",
		"fd.cip",
		"fd.cport",
		"proc.cmdline",
		"proc.name",
	}
	policy.NetworkConfiguration.Outbound.Fields = policy.NetworkConfiguration.Inbound.Fields
	policy.NetworkConfiguration.ListeningPorts.Fields = []string{
		"fd.l4proto",
		"fd.sip",
		"fd.sport",
		"proc.cmdline",
		"proc.name",
	}

	if inboundAction := d.Get("network.0.inbound").(string); inboundAction != "" {
		policy.NetworkConfiguration.Inbound.OnDefault = defaultMatchActions[inboundAction]
	} else {
		policy.NetworkConfiguration.Inbound.OnDefault = defaultMatchActions["none"]
	}

	if outboundAction := d.Get("network.0.outbound").(string); outboundAction != "" {
		policy.NetworkConfiguration.Outbound.OnDefault = defaultMatchActions[outboundAction]
	} else {
		policy.NetworkConfiguration.Outbound.OnDefault = defaultMatchActions["none"]
	}

	if defaultAction := d.Get("network.0.listening_ports.0.default").(string); defaultAction != "" {
		policy.NetworkConfiguration.ListeningPorts.OnDefault = defaultMatchActions[defaultAction]
	} else {
		policy.NetworkConfiguration.ListeningPorts.OnDefault = defaultMatchActions["none"]
	}

	if tcpActions := d.Get("network.0.listening_ports.0.tcp.0").(map[string]interface{}); len(tcpActions) != 0 {
		if tcpWhitelist := d.Get("network.0.listening_ports.0.tcp.0.whitelist").(*schema.Set); tcpWhitelist.Len() != 0 {
			element := NetworkConfigurationListeningPortsListElement{
				NetProto: "PROTO_TCP",
				OnMatch:  matchActions["accept"],
			}
			for _, port := range tcpWhitelist.List() {
				element.Values = append(element.Values, strconv.Itoa(port.(int)))
			}
			policy.NetworkConfiguration.ListeningPorts.List = append(policy.NetworkConfiguration.ListeningPorts.List, element)
		}
		if tcpBlackList := d.Get("network.0.listening_ports.0.tcp.0.blacklist").(*schema.Set); tcpBlackList.Len() != 0 {
			element := NetworkConfigurationListeningPortsListElement{
				NetProto: "PROTO_TCP",
				OnMatch:  matchActions["deny"],
			}
			for _, port := range tcpBlackList.List() {
				element.Values = append(element.Values, strconv.Itoa(port.(int)))
			}
			policy.NetworkConfiguration.ListeningPorts.List = append(policy.NetworkConfiguration.ListeningPorts.List, element)
		}
	}

	if udpActions := d.Get("network.0.listening_ports.0.udp.0").(map[string]interface{}); len(udpActions) != 0 {
		if udpWhiteList := d.Get("network.0.listening_ports.0.udp.0.whitelist").(*schema.Set); udpWhiteList.Len() != 0 {
			element := NetworkConfigurationListeningPortsListElement{
				NetProto: "PROTO_UDP",
				OnMatch:  matchActions["accept"],
			}
			for _, port := range udpWhiteList.List() {
				element.Values = append(element.Values, strconv.Itoa(port.(int)))
			}
			policy.NetworkConfiguration.ListeningPorts.List = append(policy.NetworkConfiguration.ListeningPorts.List, element)
		}
		if udpBlackList := d.Get("network.0.listening_ports.0.udp.0.blacklist").(*schema.Set); udpBlackList.Len() != 0 {
			element := NetworkConfigurationListeningPortsListElement{
				NetProto: "PROTO_UDP",
				OnMatch:  matchActions["deny"],
			}
			for _, port := range udpBlackList.List() {
				element.Values = append(element.Values, strconv.Itoa(port.(int)))
			}
			policy.NetworkConfiguration.ListeningPorts.List = append(policy.NetworkConfiguration.ListeningPorts.List, element)
		}

	}

}

func resourceSysdigPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())
	policy, err := client.GetPolicyById(id)

	if err != nil {
		d.SetId("")
	}

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("container_scope", policy.ContainerScope)
	d.Set("host_scope", policy.HostScope)
	d.Set("enabled", policy.Enabled)
	d.Set("filter", policy.Scope)
	d.Set("version", policy.Version)
	for _, action := range policy.Actions {
		if action.Type != "POLICY_ACTION_CAPTURE" {
			action := strings.Replace(action.Type, "POLICY_ACTION_", "", 1)
			d.Set("actions.0.container", strings.ToLower(action))
		} else {
			d.Set("actions.0.capture.seconds_after_event", action.AfterEventNs/1000000000)
			d.Set("actions.0.capture.seconds_before_event", action.BeforeEventNs/100000000)
		}
	}
	resourceSysdigPolicyReadProcesses(policy, d)
	resourceSysdigPolicyReadContainers(policy, d)
	resourceSysdigPolicyReadSyscalls(policy, d)
	resourceSysdigPolicyReadNetwork(policy, d)
	resourceSysdigPolicyReadFilesystem(policy, d)

	var ncIds []string
	for _, idStr := range policy.NotificationChannelIds {
		ncIds = append(ncIds, strconv.Itoa(idStr))
	}
	d.Set("notification_channels", ncIds)

	if policy.FalcoConfiguration.RuleNameRegEx != "" {
		d.Set("falco_rule_name_regex", policy.FalcoConfiguration.RuleNameRegEx)
	}

	return nil
}
func resourceSysdigPolicyReadFilesystem(policy Policy, d *schema.ResourceData) {
	if policy.FileSystemConfiguration == nil {
		return
	}

	for key, value := range defaultMatchActions {
		if policy.FileSystemConfiguration.OnDefault == value {
			d.Set("filesystem.0.other_paths", key)
			break
		}
	}

	var readWhitelist []string
	var readBlacklist []string
	var readwriteWhitelist []string
	var readwriteBlacklist []string

	for _, element := range policy.FileSystemConfiguration.List {
		if element.AccessType == "ACCESS_READ" {
			if element.OnMatch == defaultMatchActions["accept"] {
				readWhitelist = append(readWhitelist, element.Values...)
			} else if element.OnMatch == defaultMatchActions["deny"] {
				readBlacklist = append(readBlacklist, element.Values...)
			}
		} else if element.AccessType == "ACCESS_READWRITE" {
			if element.OnMatch == defaultMatchActions["accept"] {
				readwriteWhitelist = append(readwriteWhitelist, element.Values...)
			} else if element.OnMatch == defaultMatchActions["deny"] {
				readwriteBlacklist = append(readwriteBlacklist, element.Values...)
			}
		}
	}

	d.Set("filesystem.0.read.whitelist", readWhitelist)
	d.Set("filesystem.0.read.blacklist", readBlacklist)
	d.Set("filesystem.0.readwrite.whitelist", readwriteWhitelist)
	d.Set("filesystem.0.readwrite.blacklist", readwriteBlacklist)
}

func resourceSysdigPolicyReadNetwork(policy Policy, d *schema.ResourceData) {
	if policy.NetworkConfiguration == nil {
		return // No network configuration set
	}

	d.Set("network.0.inbound", defaultMatchActions["none"])
	d.Set("network.0.outbound", defaultMatchActions["none"])
	d.Set("network.0.listening_ports.0.default", defaultMatchActions["none"])
	for key, value := range defaultMatchActions {
		if policy.NetworkConfiguration.Inbound.OnDefault == value {
			d.Set("network.0.inbound", key)
		}
		if policy.NetworkConfiguration.Outbound.OnDefault == value {
			d.Set("network.0.outbound", key)
		}
		if policy.NetworkConfiguration.ListeningPorts.OnDefault == value {
			d.Set("network.0.listening_ports.0.default", key)
		}
	}

	var tcpAccept []int
	var tcpDeny []int
	var udpAccept []int
	var udpDeny []int

	stringSliceToIntSlice := func(data []string) (result []int) {
		for _, value := range data {
			integer, _ := strconv.Atoi(value)
			result = append(result, integer)
		}
		return
	}

	for _, element := range policy.NetworkConfiguration.ListeningPorts.List {
		if element.NetProto == "PROTO_TCP" {
			if element.OnMatch == matchActions["accept"] {
				tcpAccept = append(tcpAccept, stringSliceToIntSlice(element.Values)...)
			} else if element.OnMatch == matchActions["deny"] {
				tcpDeny = append(tcpDeny, stringSliceToIntSlice(element.Values)...)
			}
		} else if element.NetProto == "PROTO_UDP" {
			if element.OnMatch == matchActions["accept"] {
				udpAccept = append(udpAccept, stringSliceToIntSlice(element.Values)...)
			} else if element.OnMatch == matchActions["deny"] {
				udpDeny = append(udpDeny, stringSliceToIntSlice(element.Values)...)
			}
		}
	}
	d.Set("network.0.listening_ports.0.tcp.0.whitelist", tcpAccept)
	d.Set("network.0.listening_ports.0.tcp.0.blacklist", tcpDeny)
	d.Set("network.0.listening_ports.0.udp.0.whitelist", udpAccept)
	d.Set("network.0.listening_ports.0.udp.0.blacklist", udpDeny)
}

func resourceSysdigPolicyReadProcesses(policy Policy, d *schema.ResourceData) {
	if policy.ProcessesConfiguration == nil {
		return
	}

	for key, value := range defaultMatchActions {
		if policy.ProcessesConfiguration.OnDefault == value {
			d.Set("processes.0.default", key)
			break
		}
	}
	var whiteList []string
	var blackList []string
	for _, element := range policy.ProcessesConfiguration.List {
		if defaultMatchActions["accept"] == element.OnMatch {
			whiteList = append(whiteList, element.Values...)
		} else if defaultMatchActions["deny"] == element.OnMatch {
			blackList = append(blackList, element.Values...)
		}
	}
	d.Set("processes.0.whitelist", whiteList)
	d.Set("processes.0.blacklist", blackList)
}

func resourceSysdigPolicyReadContainers(policy Policy, d *schema.ResourceData) {
	if policy.ContainerImagesConfiguration == nil {
		return
	}

	for key, value := range defaultMatchActions {
		if policy.ContainerImagesConfiguration.OnDefault == value {
			d.Set("containers.0.default", key)
			break
		}
	}
	var whiteList []string
	var blackList []string
	for _, element := range policy.ContainerImagesConfiguration.List {
		if defaultMatchActions["accept"] == element.OnMatch {
			whiteList = append(whiteList, element.Values...)
		} else if defaultMatchActions["deny"] == element.OnMatch {
			blackList = append(blackList, element.Values...)
		}
	}
	d.Set("containers.0.whitelist", whiteList)
	d.Set("containers.0.blacklist", blackList)
}

func resourceSysdigPolicyReadSyscalls(policy Policy, d *schema.ResourceData) {
	if policy.SyscallsConfiguration == nil {
		return
	}

	for key, value := range defaultMatchActions {
		if policy.SyscallsConfiguration.OnDefault == value {
			d.Set("syscalls.0.default", key)
			break
		}
	}
	whiteList := []string{}
	blackList := []string{}
	for _, element := range policy.SyscallsConfiguration.List {
		if defaultMatchActions["accept"] == element.OnMatch {
			whiteList = append(whiteList, element.Values...)
		} else if defaultMatchActions["deny"] == element.OnMatch {
			blackList = append(blackList, element.Values...)
		}
	}
	d.Set("syscalls.0.whitelist", whiteList)
	d.Set("syscalls.0.blacklist", blackList)

}

func resourceSysdigPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())

	return client.DeletePolicy(id)
}

func resourceSysdigPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SysdigSecureClient)

	policy := policyFromResourceData(d)
	policy.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	policy.ID = id

	_, err := client.UpdatePolicy(policy)
	return err
}
