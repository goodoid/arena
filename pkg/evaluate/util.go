package evaluate

import (
	"fmt"
	"github.com/kubeflow/arena/pkg/apis/types"
	"io"
	batchv1 "k8s.io/api/batch/v1"
	"strings"
	"time"
)

const (
	formatLayout = "2006-01-02T15:04:05Z"
)

func printLine(w io.Writer, fields ...string) {
	buffer := strings.Join(fields, "\t")
	fmt.Fprintln(w, buffer)
}

func parseTime(strTime string) (time.Time, error) {
	return time.Parse(formatLayout, strTime)
}

func formatTime(t time.Time) string {
	return t.Format(formatLayout)
}

func buildEvaluateJob(job *batchv1.Job) *types.EvaluateJobInfo {
	modelName := ""
	modelVersion := ""
	modelPath := ""
	datasetPath := ""
	metricsPath := ""
	envs := job.Spec.Template.Spec.Containers[0].Env
	for _, env := range envs {
		if env.Name == "MODEL_NAME" {
			modelName = env.Value
		}
		if env.Name == "MODEL_VERSION" {
			modelVersion = env.Value
		}
		if env.Name == "MODEL_PATH" {
			modelPath = env.Value
		}
		if env.Name == "DATASET_PATH" {
			datasetPath = env.Value
		}
		if env.Name == "METRICS_PATH" {
			metricsPath = env.Value
		}
	}

	jobId := ""
	for k, v := range job.Labels {
		if k == "jobId" {
			jobId = v
		}
	}

	return &types.EvaluateJobInfo{
		UUID: string(job.UID),
		JobID: jobId,
		Name: job.Name,
		Namespace: job.Namespace,
		ModelName: modelName,
		ModelPath: modelPath,
		ModelVersion: modelVersion,
		DatasetPath: datasetPath,
		MetricsPath: metricsPath,
		CreationTimestamp: formatTime(job.CreationTimestamp.Time),
	}
}