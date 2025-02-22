package task_metadata

import (
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/flant/shell-operator/pkg/hook/binding_context"
	. "github.com/flant/shell-operator/pkg/hook/types"
	"github.com/flant/shell-operator/pkg/task"
	"github.com/flant/shell-operator/pkg/task/dump"
	"github.com/flant/shell-operator/pkg/task/queue"
)

func Test_HookMetadata_Access(t *testing.T) {
	g := NewWithT(t)

	Task := task.NewTask(HookRun).
		WithMetadata(HookMetadata{
			HookName:    "test-hook",
			BindingType: Schedule,
			BindingContext: []BindingContext{
				{Binding: "each_1_min"},
				{Binding: "each_5_min"},
			},
			AllowFailure: true,
		})

	hm := HookMetadataAccessor(Task)

	g.Expect(hm.HookName).Should(Equal("test-hook"))
	g.Expect(hm.BindingType).Should(Equal(Schedule))
	g.Expect(hm.AllowFailure).Should(BeTrue())
	g.Expect(hm.BindingContext).Should(HaveLen(2))
	g.Expect(hm.BindingContext[0].Binding).Should(Equal("each_1_min"))
	g.Expect(hm.BindingContext[1].Binding).Should(Equal("each_5_min"))
}

func Test_HookMetadata_QueueDump_Task_Description(t *testing.T) {
	g := NewWithT(t)

	logLabels := map[string]string{
		"hook": "hook1.sh",
	}
	q := queue.NewTasksQueue()

	q.AddLast(task.NewTask(EnableKubernetesBindings).
		WithMetadata(HookMetadata{
			HookName: "hook1.sh",
			Binding:  string(EnableKubernetesBindings),
		}))

	q.AddLast(task.NewTask(HookRun).
		WithMetadata(HookMetadata{
			HookName:    "hook1.sh",
			BindingType: OnKubernetesEvent,
			Binding:     "monitor_pods",
		}).
		WithLogLabels(logLabels).
		WithQueueName("main"))

	q.AddLast(task.NewTask(HookRun).
		WithMetadata(HookMetadata{
			HookName:     "hook1.sh",
			BindingType:  Schedule,
			AllowFailure: true,
			Binding:      "every 1 sec",
			Group:        "monitor_pods",
		}).
		WithLogLabels(logLabels).
		WithQueueName("main"))

	queueDump := dump.TaskQueueToText(q)

	g.Expect(queueDump).Should(ContainSubstring("hook1.sh"), "Queue dump should reveal a hook name.")
	g.Expect(queueDump).Should(ContainSubstring("EnableKubernetesBindings"), "Queue dump should reveal EnableKubernetesBindings.")
	g.Expect(queueDump).Should(ContainSubstring(":kubernetes:"), "Queue dump should show kubernetes binding.")
	g.Expect(queueDump).Should(ContainSubstring(":schedule:"), "Queue dump should show schedule binding.")
	g.Expect(queueDump).Should(ContainSubstring("group=monitor_pods"), "Queue dump should show group name.")
}
