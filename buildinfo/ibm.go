//go:build tf_acc_ibm || tf_acc_ibm_monitor

package buildinfo

func init() {
	IBMMonitor = true
}
