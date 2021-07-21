package instance_controller_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	virtv1 "kubevirt.io/client-go/api/v1"

	crownlabsv1alpha2 "github.com/netgroup-polito/CrownLabs/operators/api/v1alpha2"
	instance_controller "github.com/netgroup-polito/CrownLabs/operators/pkg/instance-controller"
)

var _ = Describe("Status Inspection", func() {

	Describe("The statusinspection.RetrievePhaseFromVM function", func() {
		var reconciler instance_controller.InstanceReconciler

		ForgeVM := func(status virtv1.VirtualMachinePrintableStatus) *virtv1.VirtualMachine {
			return &virtv1.VirtualMachine{Status: virtv1.VirtualMachineStatus{PrintableStatus: status, Ready: false}}
		}

		ForgeReadyVM := func() *virtv1.VirtualMachine {
			return &virtv1.VirtualMachine{Status: virtv1.VirtualMachineStatus{
				PrintableStatus: virtv1.VirtualMachineStatusRunning,
				Ready:           true,
			}}
		}

		BeforeEach(func() {
			reconciler = instance_controller.InstanceReconciler{}
		})

		DescribeTable("Correctly returns the expected instance phase",
			func(vm *virtv1.VirtualMachine, expected crownlabsv1alpha2.EnvironmentPhase) {
				Expect(reconciler.RetrievePhaseFromVM(vm)).To(Equal(expected))
			},
			Entry("When the VM is starting", ForgeVM(virtv1.VirtualMachineStatusStarting), crownlabsv1alpha2.EnvironmentPhaseStarting),
			Entry("When the VM is provisioning", ForgeVM(virtv1.VirtualMachineStatusProvisioning), crownlabsv1alpha2.EnvironmentPhaseImporting),
			Entry("When the VM is stopping", ForgeVM(virtv1.VirtualMachineStatusStopping), crownlabsv1alpha2.EnvironmentPhaseStopping),
			Entry("When the VM is terminating", ForgeVM(virtv1.VirtualMachineStatusTerminating), crownlabsv1alpha2.EnvironmentPhaseStopping),
			Entry("When the VM is off", ForgeVM(virtv1.VirtualMachineStatusStopped), crownlabsv1alpha2.EnvironmentPhaseOff),
			Entry("When the VM is running", ForgeVM(virtv1.VirtualMachineStatusRunning), crownlabsv1alpha2.EnvironmentPhaseRunning),
			Entry("When the VM is ready", ForgeReadyVM(), crownlabsv1alpha2.EnvironmentPhaseReady),
			Entry("When the VM status is unknown", ForgeVM(virtv1.VirtualMachineStatusUnknown), crownlabsv1alpha2.EnvironmentPhaseUnset),
		)
	})

	Describe("The statusinspection.RetrievePhaseFromVMI function", func() {
		var reconciler instance_controller.InstanceReconciler

		ForgeVMI := func(phase virtv1.VirtualMachineInstancePhase) *virtv1.VirtualMachineInstance {
			return &virtv1.VirtualMachineInstance{Status: virtv1.VirtualMachineInstanceStatus{
				Phase: phase,
				Conditions: []virtv1.VirtualMachineInstanceCondition{
					{Type: virtv1.VirtualMachineInstanceReady, Status: v1.ConditionFalse},
					{Type: virtv1.VirtualMachineInstanceIsMigratable, Status: v1.ConditionTrue},
					{Type: virtv1.VirtualMachineInstancePaused, Status: v1.ConditionFalse},
				},
			}}
		}

		ForgeReadyVMI := func() *virtv1.VirtualMachineInstance {
			return &virtv1.VirtualMachineInstance{Status: virtv1.VirtualMachineInstanceStatus{
				Phase: virtv1.Running,
				Conditions: []virtv1.VirtualMachineInstanceCondition{
					{Type: virtv1.VirtualMachineInstanceReady, Status: v1.ConditionTrue},
					{Type: virtv1.VirtualMachineInstanceIsMigratable, Status: v1.ConditionTrue},
					{Type: virtv1.VirtualMachineInstancePaused, Status: v1.ConditionFalse},
				},
			}}
		}

		ForgeStoppingVMI := func() *virtv1.VirtualMachineInstance {
			timestamp := metav1.NewTime(time.Now())
			return &virtv1.VirtualMachineInstance{ObjectMeta: metav1.ObjectMeta{DeletionTimestamp: &timestamp}}
		}

		BeforeEach(func() {
			reconciler = instance_controller.InstanceReconciler{}
		})

		DescribeTable("Correctly returns the expected instance phase",
			func(vmi *virtv1.VirtualMachineInstance, expected crownlabsv1alpha2.EnvironmentPhase) {
				Expect(reconciler.RetrievePhaseFromVMI(vmi)).To(Equal(expected))
			},
			Entry("When the VMI status is unset", ForgeVMI(virtv1.VmPhaseUnset), crownlabsv1alpha2.EnvironmentPhaseStarting),
			Entry("When the VMI is pending", ForgeVMI(virtv1.Pending), crownlabsv1alpha2.EnvironmentPhaseStarting),
			Entry("When the VMI is scheduling", ForgeVMI(virtv1.Scheduling), crownlabsv1alpha2.EnvironmentPhaseStarting),
			Entry("When the VMI is scheduled", ForgeVMI(virtv1.Scheduled), crownlabsv1alpha2.EnvironmentPhaseStarting),
			Entry("When the VMI is running", ForgeVMI(virtv1.Running), crownlabsv1alpha2.EnvironmentPhaseRunning),
			Entry("When the VMI is ready", ForgeReadyVMI(), crownlabsv1alpha2.EnvironmentPhaseReady),
			Entry("When the VMI status is unknown", ForgeVMI(virtv1.Unknown), crownlabsv1alpha2.EnvironmentPhaseFailed),
			Entry("When the VMI status is failed", ForgeVMI(virtv1.Failed), crownlabsv1alpha2.EnvironmentPhaseFailed),
			Entry("When the VMI status is succeeded", ForgeVMI(virtv1.Succeeded), crownlabsv1alpha2.EnvironmentPhaseFailed),
			Entry("When the VMI is being deleted", ForgeStoppingVMI(), crownlabsv1alpha2.EnvironmentPhaseStopping),
		)
	})
})
