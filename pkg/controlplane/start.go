package controlplane

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/percona/everest-backend/models"
)

func Start(cluster *models.DBCluster) (string, error) {
	cmd := exec.Command("helm", "install", fmt.Sprintf("cluster-%d", cluster.ID), "percona/ps-db")
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return out.String(), err
	}

	return out.String(), nil
}
