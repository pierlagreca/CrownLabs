package instance_controller

import (
	v1 "k8s.io/api/core/v1"
	virtv1 "kubevirt.io/client-go/api/v1"

	crownlabsv1alpha2 "github.com/netgroup-polito/CrownLabs/operators/api/v1alpha2"
)

// RetrievePhaseFromVM converts the VM phase to the corresponding one of the instance.
func (r *InstanceReconciler) RetrievePhaseFromVM(vm *virtv1.VirtualMachine) crownlabsv1alpha2.EnvironmentPhase {
	switch vm.Status.PrintableStatus {
	case virtv1.VirtualMachineStatusStarting:
		return crownlabsv1alpha2.EnvironmentPhaseStarting
	case virtv1.VirtualMachineStatusProvisioning:
		return crownlabsv1alpha2.EnvironmentPhaseImporting

	case virtv1.VirtualMachineStatusStopping:
		return crownlabsv1alpha2.EnvironmentPhaseStopping
	case virtv1.VirtualMachineStatusTerminating:
		return crownlabsv1alpha2.EnvironmentPhaseStopping
	case virtv1.VirtualMachineStatusStopped:
		return crownlabsv1alpha2.EnvironmentPhaseOff

	case virtv1.VirtualMachineStatusRunning:
		if vm.Status.Ready {
			return crownlabsv1alpha2.EnvironmentPhaseReady
		}
		return crownlabsv1alpha2.EnvironmentPhaseRunning

	default:
		return crownlabsv1alpha2.EnvironmentPhaseUnset
	}
}

// RetrievePhaseFromVMI converts the VMI phase to the corresponding one of the instance.
func (r *InstanceReconciler) RetrievePhaseFromVMI(vmi *virtv1.VirtualMachineInstance) crownlabsv1alpha2.EnvironmentPhase {
	if !vmi.DeletionTimestamp.IsZero() {
		return crownlabsv1alpha2.EnvironmentPhaseStopping
	}

	switch vmi.Status.Phase {
	case virtv1.VmPhaseUnset:
		return crownlabsv1alpha2.EnvironmentPhaseStarting
	case virtv1.Pending:
		return crownlabsv1alpha2.EnvironmentPhaseStarting
	case virtv1.Scheduling:
		return crownlabsv1alpha2.EnvironmentPhaseStarting
	case virtv1.Scheduled:
		return crownlabsv1alpha2.EnvironmentPhaseStarting

	case virtv1.Unknown:
		return crownlabsv1alpha2.EnvironmentPhaseFailed
	case virtv1.Failed:
		return crownlabsv1alpha2.EnvironmentPhaseFailed
	case virtv1.Succeeded:
		return crownlabsv1alpha2.EnvironmentPhaseFailed

	case virtv1.Running:
		if isVMIReady(vmi) {
			return crownlabsv1alpha2.EnvironmentPhaseReady
		}
		return crownlabsv1alpha2.EnvironmentPhaseRunning

	default:
		return crownlabsv1alpha2.EnvironmentPhaseUnset
	}
}

// isVMIReady checks whether a VMI is ready, depending on its conditions.
func isVMIReady(vmi *virtv1.VirtualMachineInstance) bool {
	for _, condition := range vmi.Status.Conditions {
		if condition.Type == virtv1.VirtualMachineInstanceReady {
			return condition.Status == v1.ConditionTrue
		}
	}

	return false
}
