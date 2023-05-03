package controlplane

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/percona/everest-backend/models"
)

type kubectlGetPodsResponse struct {
	Status kubectlGetPodsStatus `json:"status"`
}

type kubectlGetPodsStatus struct {
	Phase string `json:"phase"`
}

func Status(cluster *models.DBCluster) (string, string, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-o", "json", fmt.Sprintf("cluster-%d-ps-db-mysql-0", cluster.ID))
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", out.String(), err
	}

	var res kubectlGetPodsResponse
	if err := json.Unmarshal([]byte(out.String()), &res); err != nil {
		return "", out.String(), fmt.Errorf("cannot unmarshal kubectl get pods response. %w", err)
	}

	return res.Status.Phase, out.String(), nil
}
