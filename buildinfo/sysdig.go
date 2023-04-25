//go:build tf_acc_sysdig

package buildinfo

func init() {
	SysdigMonitor = true
	SysdigSecure = true
}
