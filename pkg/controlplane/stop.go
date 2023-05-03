package controlplane

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/percona/everest-backend/models"
)

func Stop(cluster *models.DBCluster) (string, error) {
	cmd := exec.Command("helm", "uninstall", fmt.Sprintf("cluster-%d", cluster.ID))
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return out.String(), err
	}

	return out.String(), nil
}
